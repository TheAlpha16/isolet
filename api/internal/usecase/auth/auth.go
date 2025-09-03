package auth

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/infra/jwt"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils"
)

type authImpl struct {
	repo    authDom.Repository
	userUc  userDom.Usecase
	tokenUc tokenDom.Usecase
	jwtSvc  jwt.JWT
}

func (a *authImpl) Login(ctx context.Context, input *authDom.LoginInput) (*authDom.Session, error) {
	// hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errorDom.RaiseInternal(ctx, "failed to hash password", err, common.ExtraData{"password": input.Password})
	}
	input.Password = hashedPassword

	return &authDom.Session{
		UserID:    6969,
		Token:     "some-token-here",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
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
	if err := a.ensureEmailAndUsernameAvailable(ctx, input.Email, input.Username); err != nil {
		return nil, err
	}

	// verify the email in case enabled
	// TODO implement email send
	token, err := a.createEmailVerificationToken(ctx, input.Email, input.Username, input.Password)
	if err != nil {
		return nil, err
	}

	// TODO think from here
	// I have user_id as a required field in the jwt claims
	// this means that i am expecting the user_id to be present in the token for all tokens which might not be possible everytime
	jwtClaims := jwt.Claims{
		JWTID:     token.ID,
		UserID:    6969,
		Role:      userDom.RolePlayer,
		Purpose:   tokenDom.TokenEmailVerification,
		CreatedAt: token.CreatedAt,
		ExpiresAt: token.ExpiresAt,
	}

	jwtToken, err := a.jwtSvc.Sign(ctx, &jwtClaims)
	if err != nil {
		return nil, err
	}

	return &authDom.Session{
		UserID:    6969,
		Token:     jwtToken,
		ExpiresAt: token.ExpiresAt.Unix(),
	}, nil
}

func (a *authImpl) ensureEmailAndUsernameAvailable(ctx context.Context, email, username string) error {
	existence, err := a.userUc.ExistsByEmailOrUsername(ctx, email, username)
	if err != nil {
		return err
	}
	if existence.EmailExists {
		return errorDom.Raise(ctx, errorDom.ErrUserEmailTaken, "", nil, nil)
	}
	if existence.UsernameExists {
		return errorDom.Raise(ctx, errorDom.ErrUserUsernameTaken, "", nil, nil)
	}
	return nil
}

func (a *authImpl) createEmailVerificationToken(ctx context.Context, email, username, password string) (*tokenDom.Token, error) {
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
		return nil, err
	}
	return &token, nil
}

func New(repo authDom.Repository, userUc userDom.Usecase, tokenUc tokenDom.Usecase, jwtSvc jwt.JWT) authDom.Usecase {
	return &authImpl{
		repo:    repo,
		userUc:  userUc,
		tokenUc: tokenUc,
		jwtSvc:  jwtSvc,
	}
}
