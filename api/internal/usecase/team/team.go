package team

import (
	"context"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils"
)

type teamImpl struct {
	repo    teamDom.Repository
	authUc  authDom.Usecase
	userUc  userDom.Usecase
	tokenUc tokenDom.Usecase
}

func (t *teamImpl) Create(ctx context.Context, input *teamDom.CreateInput) (*authDom.Session, error) {
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	input.Password = hashedPassword

	team, err := t.repo.Create(ctx, &teamDom.Team{})
	if err != nil {
		return nil, err
	}

	captain, err := t.userUc.GetByID(ctx, team.CaptainID)
	if err != nil {
		return nil, err
	}

	// revoke all the auth tokens
	if err := t.tokenUc.RevokeUserTokens(ctx, tokenDom.TokenAuth, captain.ID); err != nil {
		return nil, err
	}

	token, jwtToken, err := t.authUc.GenerateAuthToken(ctx, captain)
	if err != nil {
		return nil, err
	}

	return &authDom.Session{
		UserID:    team.CaptainID,
		Token:     jwtToken,
		ExpiresAt: token.ExpiresAt.Unix(),
	}, nil
}

func New(repo teamDom.Repository, userUc userDom.Usecase, authUc authDom.Usecase, tokenUc tokenDom.Usecase) teamDom.Usecase {
	return &teamImpl{
		repo:    repo,
		authUc:  authUc,
		userUc:  userUc,
		tokenUc: tokenUc,
	}
}
