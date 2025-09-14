package team

import (
	"context"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
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
	userID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyUserID)

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}
	input.Password = hashedPassword

	team, err := t.repo.Create(ctx, &teamDom.Team{
		Name:      input.TeamName,
		CaptainID: userID,
		Password:  input.Password,
	})
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

func (t *teamImpl) Join(ctx context.Context, input *teamDom.JoinInput) (*authDom.Session, error) {
	userID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyUserID)

	team, err := t.repo.GetByName(ctx, input.TeamName)
	if err != nil {
		if errorDom.IsSameError(err, errorDom.ErrTeamNotFound) {
			return nil, errorDom.Raise(ctx, errorDom.ErrAuthInvalidCredentials, "", nil, nil)
		}
		return nil, err
	}

	// verify the password
	if !utils.ComparePassword(team.Password, input.Password) {
		return nil, errorDom.Raise(ctx, errorDom.ErrAuthInvalidCredentials, "", nil, nil)
	}

	user, err := t.userUc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// update the user's team
	user.TeamID = &team.ID
	if err := t.userUc.Update(ctx, user, []string{"team_id"}); err != nil {
		return nil, err
	}

	// revoke all the auth tokens
	if err := t.tokenUc.RevokeUserTokens(ctx, tokenDom.TokenAuth, user.ID); err != nil {
		return nil, err
	}

	token, jwtToken, err := t.authUc.GenerateAuthToken(ctx, user)
	if err != nil {
		return nil, err
	}

	return &authDom.Session{
		UserID:    user.ID,
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
