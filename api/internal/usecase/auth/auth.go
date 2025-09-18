package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/TheAlpha16/isolet/api/infra/jwt"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	emailDom "github.com/TheAlpha16/isolet/api/internal/domain/email"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils"
)

type authImpl struct {
	userUc  userDom.Usecase
	tokenUc tokenDom.Usecase
	cvUc    cvDom.Usecase
	emailUc emailDom.Usecase
	jwtSvc  jwt.JWT
}

func (a *authImpl) Login(ctx context.Context, input *authDom.LoginInput) (*authDom.Session, error) {
	user, err := a.userUc.GetByEmailOrUsername(ctx, strings.ToLower(input.Identifier), input.Identifier)
	if err != nil {
		if errorDom.IsSameError(err, errorDom.ErrUserNotFound) {
			return nil, errorDom.Raise(ctx, errorDom.ErrAuthInvalidCredentials, "", nil, nil)
		}
		return nil, err
	}

	// verify the password
	if !utils.ComparePassword(user.Password, input.Password) {
		return nil, errorDom.Raise(ctx, errorDom.ErrAuthInvalidCredentials, "", nil, nil)
	}

	// don't allow banned users to Login
	if user.IsBanned {
		return nil, errorDom.Raise(ctx, errorDom.ErrAuthUserBanned, "", nil, nil)
	}

	activeSessionCount, err := a.tokenUc.CountEntityTokens(ctx, tokenDom.TokenAuth, fmt.Sprintf("%d", user.ID))
	if err != nil {
		return nil, err
	}
	if activeSessionCount >= a.cvUc.GetInt(ctx, cvDom.AuthMaxSessions) {
		return nil, errorDom.Raise(ctx, errorDom.ErrAuthMaxSessionsReached, "max auth sessions reached", nil, nil)
	}

	token, jwtToken, err := a.GenerateAuthToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &authDom.Session{
		UserID:    user.ID,
		Token:     jwtToken,
		ExpiresAt: token.ExpiresAt.Unix(),
	}, nil
}

func (a *authImpl) Register(ctx context.Context, input *authDom.RegisterInput) (*authDom.Session, error) {
	// hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errorDom.RaiseInternal(ctx, "failed to hash password", err, common.ExtraData{"password": input.Password})
	}
	input.Password = hashedPassword

	// check if username or email is already taken
	if err := a.userUc.ExistsByEmailOrUsername(ctx, input.Email, input.Username); err != nil {
		return nil, err
	}

	// verify the email in case enabled
	if a.cvUc.GetBool(ctx, cvDom.EmailVerificationEnabled) {
		_, jwtToken, err := a.generateEmailVerificationToken(ctx, input.Email, input.Username, input.Password)
		if err != nil {
			return nil, err
		}

		if err = a.emailUc.SendEmailAsync(ctx, &emailDom.EmailInput{
			EmailIdentifier: emailDom.EmailIdentifier{
				Username: input.Username,
				To:       input.Email,
			},
			Type:  emailDom.TypeVerification,
			Token: jwtToken,
		}); err != nil {
			return nil, err
		}
		return nil, nil
	}

	user, err := a.createUser(ctx, input.Email, input.Username, input.Password)
	if err != nil {
		return nil, err
	}

	token, jwtToken, err := a.GenerateAuthToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &authDom.Session{
		UserID:    user.ID,
		Token:     jwtToken,
		ExpiresAt: token.ExpiresAt.Unix(),
	}, nil
}

func (a *authImpl) Verify(ctx context.Context, verifyToken string) error {
	claims, err := a.jwtSvc.Verify(ctx, verifyToken)
	if err != nil {
		return err
	}

	if claims.Purpose != tokenDom.TokenEmailVerification {
		return errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "", nil, nil)
	}

	tokenIdentifier := &tokenDom.TokenIdentifier{
		ID:       claims.JWTID,
		EntityID: claims.Subject,
		Purpose:  claims.Purpose,
	}

	token, err := a.tokenUc.Fetch(ctx, tokenIdentifier)
	if err != nil {
		return err
	}

	var email, username, password string
	var ok bool
	if email, ok = token.Metadata["email"]; !ok {
		return errorDom.Raise(ctx, errorDom.ErrTokenExpiredInvalid, "", nil, nil)
	}
	if username, ok = token.Metadata["username"]; !ok {
		return errorDom.Raise(ctx, errorDom.ErrTokenExpiredInvalid, "", nil, nil)
	}
	if password, ok = token.Metadata["password"]; !ok {
		return errorDom.Raise(ctx, errorDom.ErrTokenExpiredInvalid, "", nil, nil)
	}

	_, err = a.createUser(ctx, email, username, password)
	if err != nil {
		return err
	}

	return a.tokenUc.Delete(ctx, tokenIdentifier)
}

func (a *authImpl) ForgotPassword(ctx context.Context, input *authDom.ForgotPasswordInput) error {
	if !a.cvUc.GetBool(ctx, cvDom.PasswordResetEnabled) {
		return errorDom.Raise(ctx, errorDom.ErrAuthPasswordResetDisabled, "", nil, nil)
	}

	user, err := a.userUc.GetByEmailOrUsername(ctx, input.Email, "")
	if err != nil {
		if errorDom.IsSameError(err, errorDom.ErrUserNotFound) {
			// user not found -> do nothing
			return nil
		}
		return err
	}

	activeSessionCount, err := a.tokenUc.CountEntityTokens(ctx, tokenDom.TokenPasswordReset, fmt.Sprintf("%d", user.ID))
	if err != nil {
		return err
	}
	if activeSessionCount >= a.cvUc.GetInt(ctx, cvDom.PasswordResetMaxSessions) {
		return errorDom.Raise(ctx, errorDom.ErrAuthMaxSessionsReached, "max password reset sessions reached", nil, nil)
	}

	_, jwtToken, err := a.generatePasswordResetToken(ctx, user)
	if err != nil {
		return err
	}

	return a.emailUc.SendEmailAsync(ctx, &emailDom.EmailInput{
		EmailIdentifier: emailDom.EmailIdentifier{
			Username: user.Username,
			To:       user.Email,
		},
		Type:  emailDom.TypePasswordReset,
		Token: jwtToken,
	})
}

func (a *authImpl) ResetPassword(ctx context.Context, input *authDom.ResetPasswordInput) error {
	claims, err := a.jwtSvc.Verify(ctx, input.Token)
	if err != nil {
		return err
	}

	if claims.Purpose != tokenDom.TokenPasswordReset {
		return errorDom.Raise(ctx, errorDom.ErrTokenExpiredInvalid, "", nil, nil)
	}

	tokenIdentifier := &tokenDom.TokenIdentifier{
		ID:       claims.JWTID,
		EntityID: claims.Subject,
		Purpose:  claims.Purpose,
	}

	_, err = a.tokenUc.Fetch(ctx, tokenIdentifier)
	if err != nil {
		return err
	}

	// hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return errorDom.RaiseInternal(ctx, "failed to hash password", err, common.ExtraData{"password": input.Password})
	}
	input.Password = hashedPassword

	// update the user's password
	if err := a.userUc.Update(ctx, &userDom.User{
		ID:       *claims.UserID,
		Password: input.Password,
	}, []string{"password"}); err != nil {
		return err
	}

	// delete the token
	return a.tokenUc.Delete(ctx, tokenIdentifier)
}

func (a *authImpl) GenerateAuthToken(ctx context.Context, user *userDom.User) (*tokenDom.Token, string, error) {
	config := utils.GetConfig()
	token := tokenDom.Token{
		TokenIdentifier: tokenDom.TokenIdentifier{
			ID:       utils.RandomUUID(),
			EntityID: strconv.FormatInt(user.ID, 10),
			Purpose:  tokenDom.TokenAuth,
		},
	}
	token.UpdateTime()
	token.ExpiresAt = token.CreatedAt.Add(config.Token.AuthValidity)

	err := a.tokenUc.Create(ctx, &token, nil)
	if err != nil {
		return nil, "", err
	}

	jwtToken, err := a.jwtSvc.Sign(ctx, jwt.NewAuthClaims(token.ID, user.ID, user.TeamID, user.Role, token.ExpiresAt))
	if err != nil {
		return nil, "", err
	}
	return &token, jwtToken, nil
}

func (a *authImpl) Logout(ctx context.Context) error {
	userID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyUserID)
	sessionID := common.GetFieldFromExtraData[string](ctx, utils.ContextKeySessionID)

	tokenIdentifier := &tokenDom.TokenIdentifier{
		ID:       sessionID,
		EntityID: strconv.FormatInt(userID, 10),
		Purpose:  tokenDom.TokenAuth,
	}

	return a.tokenUc.Delete(ctx, tokenIdentifier)
}

func (a *authImpl) generateEmailVerificationToken(ctx context.Context, email, username, password string) (*tokenDom.Token, string, error) {
	config := utils.GetConfig()
	token := tokenDom.Token{
		TokenIdentifier: tokenDom.TokenIdentifier{
			ID:       utils.RandomUUID(),
			EntityID: utils.RandomUUID(),
			Purpose:  tokenDom.TokenEmailVerification,
		},
		Metadata: map[string]string{
			"email":    email,
			"username": username,
			"password": password,
		},
	}
	token.UpdateTime()
	token.ExpiresAt = token.CreatedAt.Add(config.Token.EmailVerificationValidity)

	// store email and username for existence lookup later
	extras := map[string]string{
		userDom.EmailCacheKey(email):       "",
		userDom.UsernameCacheKey(username): "",
	}

	err := a.tokenUc.Create(ctx, &token, extras)
	if err != nil {
		return nil, "", err
	}

	jwtToken, err := a.jwtSvc.Sign(ctx, jwt.NewEmailVerificationClaims(token.ID, token.EntityID, token.ExpiresAt))
	if err != nil {
		return nil, "", err
	}
	return &token, jwtToken, nil
}

func (a *authImpl) generatePasswordResetToken(ctx context.Context, user *userDom.User) (*tokenDom.Token, string, error) {
	config := utils.GetConfig()
	token := tokenDom.Token{
		TokenIdentifier: tokenDom.TokenIdentifier{
			ID:       utils.RandomUUID(),
			EntityID: fmt.Sprintf("%d", user.ID),
			Purpose:  tokenDom.TokenPasswordReset,
		},
	}
	token.UpdateTime()
	token.ExpiresAt = token.CreatedAt.Add(config.Token.PasswordResetValidity)

	err := a.tokenUc.Create(ctx, &token, nil)
	if err != nil {
		return nil, "", err
	}

	jwtToken, err := a.jwtSvc.Sign(ctx, jwt.NewPasswordResetClaims(token.ID, user.ID, token.ExpiresAt))
	if err != nil {
		return nil, "", err
	}
	return &token, jwtToken, nil
}

func (a *authImpl) createUser(ctx context.Context, email, username, password string) (*userDom.User, error) {
	user := &userDom.User{
		Email:    email,
		Username: username,
		Password: password,
	}

	user, err := a.userUc.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func New(userUc userDom.Usecase, tokenUc tokenDom.Usecase, cvUc cvDom.Usecase, emailUc emailDom.Usecase, jwtSvc jwt.JWT) authDom.Usecase {
	return &authImpl{
		userUc:  userUc,
		tokenUc: tokenUc,
		cvUc:    cvUc,
		emailUc: emailUc,
		jwtSvc:  jwtSvc,
	}
}
