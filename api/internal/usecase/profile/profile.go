package profile

import (
	"context"

	"github.com/TheAlpha16/isolet/api/internal/domain/common"
	profileDom "github.com/TheAlpha16/isolet/api/internal/domain/profile"
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/utils"
)

type profileImpl struct {
	userUc  userDom.Usecase
	teamUc  teamDom.Usecase
	scoreUc scoreDom.Usecase
}

func (p *profileImpl) Me(ctx context.Context) (*profileDom.Me, error) {
	userID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyUserID)

	var user *userDom.User
	var team *teamDom.Team

	user, err := p.userUc.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.TeamID != nil {
		team, err = p.teamUc.GetByID(ctx, *user.TeamID)
		if err != nil {
			return nil, err
		}
	}

	return &profileDom.Me{
		User: *user,
		Team: team,
	}, nil
}

func (p *profileImpl) Team(ctx context.Context) (*profileDom.Team, error) {
	teamID := common.GetFieldFromExtraData[int64](ctx, utils.ContextKeyTeamID)

	team, err := p.teamUc.GetByID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	scoreRecords, err := p.scoreUc.GetTeamScoreRecords(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &profileDom.Team{
		Team:    *team,
		Records: scoreRecords,
	}, nil
}

func New(userUc userDom.Usecase, teamUc teamDom.Usecase, scoreUc scoreDom.Usecase) profileDom.Usecase {
	return &profileImpl{
		userUc:  userUc,
		teamUc:  teamUc,
		scoreUc: scoreUc,
	}
}
