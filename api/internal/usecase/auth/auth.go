package auth

import (
	"context"
	"time"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	"github.com/TheAlpha16/isolet/api/internal/domain/errors"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils"
)

type authImpl struct {
	repo    authDom.Repository
	userUc  userDom.Usecase
	tokenUc tokenDom.Usecase
}

func (a *authImpl) Login(ctx context.Context, input *authDom.LoginInput) (*authDom.Session, error) {
	// hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errors.RaiseInternal(ctx, "failed to hash password", err, errors.ExtraData{"password": input.Password})
	}
	input.Password = hashedPassword

	return &authDom.Session{
		UserID:    6969,
		Token:     "some-token-here",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (a *authImpl) Register(ctx context.Context, input *authDom.RegisterInput) (*authDom.Session, error) {
	config := utils.GetConfig()

	// hash the password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errors.RaiseInternal(ctx, "failed to hash password", err, errors.ExtraData{"password": input.Password})
	}
	input.Password = hashedPassword

	// check if username or email is already taken
	existence, err := a.userUc.ExistsByEmailOrUsername(ctx, input.Email, input.Username)
	if err != nil {
		return nil, err
	}
	if existence.EmailExists {
		return nil, errors.Raise(ctx, errors.ErrUserEmailTaken, "", nil, nil)
	}
	if existence.UsernameExists {
		return nil, errors.Raise(ctx, errors.ErrUserUsernameTaken, "", nil, nil)
	}

	// verify the email in case enabled
	// TODO implement email send
	token := tokenDom.Token{
		TokenIdentifier: tokenDom.TokenIdentifier{
			ID:       utils.RandomUUID(),
			EntityID: utils.RandomUUID(),
			Purpose:  tokenDom.TokenEmailVerification,
		},
		Metadata: map[string]string{
			"email":    input.Email,
			"username": input.Username,
			"password": input.Password,
		},
	}
	token.UpdateTime()
	token.ExpiresAt = token.CreatedAt.Add(config.Token.EmailVerificationValidity)

	// store email and username for existence lookup later
	extras := map[string]string{
		userDom.EmailCacheKey(input.Email):       "",
		userDom.UsernameCacheKey(input.Username): "",
	}

	err = a.tokenUc.Create(ctx, &token, extras)
	if err != nil {
		return nil, err
	}

	return &authDom.Session{
		UserID:    6969,
		Token:     "registration-token",
		ExpiresAt: token.ExpiresAt.Unix(),
	}, nil
}

func New(repo authDom.Repository, userUc userDom.Usecase, tokenUc tokenDom.Usecase) authDom.Usecase {
	return &authImpl{
		repo:    repo,
		userUc:  userUc,
		tokenUc: tokenUc,
	}
}
