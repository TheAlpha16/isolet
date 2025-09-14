package team

import (
	"context"
	"fmt"

	"github.com/TheAlpha16/isolet/api/infra/jwt"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
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
	cvUc    cvDom.Usecase
	jwtSvc  jwt.JWT
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

	return t.reissueAuthSession(ctx, captain)
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

	// update the user's team
	if err := t.userUc.JoinTeam(ctx, userID, team.ID); err != nil {
		return nil, err
	}

	user, err := t.userUc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return t.reissueAuthSession(ctx, user)
}

func (t *teamImpl) GenerateInvite(ctx context.Context) (*teamDom.GenerateInviteOutput, error) {
	config := utils.GetConfig()
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	token := tokenDom.Token{
		TokenIdentifier: tokenDom.TokenIdentifier{
			ID:       utils.RandomUUID(),
			EntityID: fmt.Sprintf("%d", teamID),
			Purpose:  tokenDom.TokenTeamInvite,
		},
	}
	token.UpdateTime()
	token.ExpiresAt = token.CreatedAt.Add(config.Token.TeamInviteValidity)

	err := t.tokenUc.Create(ctx, &token, nil)
	if err != nil {
		return nil, err
	}

	jwtToken, err := t.jwtSvc.Sign(ctx, jwt.NewTeamInviteClaims(token.ID, teamID, token.ExpiresAt))
	if err != nil {
		return nil, err
	}

	publicURL := t.cvUc.GetString(ctx, cvDom.PublicURL)
	inviteLink, err := utils.BuildLink(
		publicURL, jwtToken,
		[]string{
			utils.GetConfig().Rest.APIVersionPrefix,
			utils.RouteTeamInvite,
		})
	if err != nil {
		return nil, err
	}

	return &teamDom.GenerateInviteOutput{
		InviteLink: inviteLink,
	}, nil
}

func (t *teamImpl) reissueAuthSession(ctx context.Context, user *userDom.User) (*authDom.Session, error) {
	// revoke all the auth tokens
	if err := t.tokenUc.RevokeEntityTokens(ctx, tokenDom.TokenAuth, fmt.Sprintf("%d", user.ID)); err != nil {
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

func New(repo teamDom.Repository, userUc userDom.Usecase, authUc authDom.Usecase, tokenUc tokenDom.Usecase, cvUc cvDom.Usecase, jwtSvc jwt.JWT) teamDom.Usecase {
	return &teamImpl{
		repo:    repo,
		authUc:  authUc,
		userUc:  userUc,
		tokenUc: tokenUc,
		cvUc:    cvUc,
		jwtSvc:  jwtSvc,
	}
}
