package team

import (
	"context"
	"strconv"

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
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	token, err := t.getOrCreateInviteToken(ctx, teamID)
	if err != nil {
		return nil, err
	}

	jwtToken, err := t.jwtSvc.Sign(ctx, jwt.NewTeamInviteClaims(token.ID, teamID, token.ExpiresAt))
	if err != nil {
		return nil, err
	}

	inviteLink, err := t.buildInviteLink(ctx, jwtToken)
	if err != nil {
		return nil, err
	}

	return &teamDom.GenerateInviteOutput{
		InviteLink: inviteLink,
	}, nil
}

func (t *teamImpl) AcceptInvite(ctx context.Context, inviteToken string) (*authDom.Session, error) {
	userID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyUserID)

	claims, err := t.jwtSvc.Verify(ctx, inviteToken)
	if err != nil {
		return nil, err
	}

	if claims.Purpose != tokenDom.TokenTeamInvite {
		return nil, errorDom.Raise(ctx, errorDom.ErrTokenMalformed, "", nil, nil)
	}

	tokenIdentfier := &tokenDom.TokenIdentifier{
		ID:       claims.JWTID,
		EntityID: claims.Subject,
		Purpose:  claims.Purpose,
	}

	if _, err = t.tokenUc.Fetch(ctx, tokenIdentfier); err != nil {
		return nil, err
	}

	if err := t.userUc.JoinTeam(ctx, userID, *claims.TeamID); err != nil {
		return nil, err
	}

	user, err := t.userUc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return t.reissueAuthSession(ctx, user)
}

func (t *teamImpl) GetNameByIDs(ctx context.Context, teamIDs []int64) (map[int64]string, error) {
	return t.repo.GetNameByIDs(ctx, teamIDs)
}

func (t *teamImpl) reissueAuthSession(ctx context.Context, user *userDom.User) (*authDom.Session, error) {
	// revoke all the auth tokens
	if err := t.tokenUc.RevokeEntityTokens(ctx, tokenDom.TokenAuth, strconv.FormatInt(user.ID, 10)); err != nil {
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

func (t *teamImpl) getOrCreateInviteToken(ctx context.Context, teamID int64) (*tokenDom.Token, error) {
	config := utils.GetConfig()

	// check if there's a token already
	tokens, err := t.tokenUc.FetchEntityTokens(ctx, tokenDom.TokenTeamInvite, strconv.FormatInt(teamID, 10))
	if err != nil {
		return nil, err
	}
	if len(tokens) > 0 {
		return tokens[0], nil
	}

	// generate new token
	token := tokenDom.Token{
		TokenIdentifier: tokenDom.TokenIdentifier{
			ID:       utils.RandomUUID(),
			EntityID: strconv.FormatInt(teamID, 10),
			Purpose:  tokenDom.TokenTeamInvite,
		},
	}
	token.UpdateTime()
	token.ExpiresAt = token.CreatedAt.Add(config.Token.TeamInviteValidity)

	err = t.tokenUc.Create(ctx, &token, nil)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (t *teamImpl) buildInviteLink(ctx context.Context, token string) (string, error) {
	publicURL := t.cvUc.GetString(ctx, cvDom.PublicURL)
	inviteLink, err := utils.BuildLink(
		publicURL, token,
		[]string{
			utils.GetConfig().Rest.APIVersionPrefix,
			utils.RouteTeam,
			utils.RouteTeamInvite,
		})
	if err != nil {
		return "", err
	}
	return inviteLink, nil
}

func (t *teamImpl) GetByID(ctx context.Context, id int64) (*teamDom.Team, error) {
	return t.repo.GetByID(ctx, id)
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
