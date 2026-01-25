package profile

import (
	"context"

	"github.com/TheAlpha16/isolet/oracle/infra/cache"
	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	"github.com/TheAlpha16/isolet/oracle/internal/domain/common"
	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	profileDom "github.com/TheAlpha16/isolet/oracle/internal/domain/profile"
	scoreDom "github.com/TheAlpha16/isolet/oracle/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/oracle/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/oracle/internal/domain/user"
	"github.com/TheAlpha16/isolet/oracle/utils"
)

type profileImpl struct {
	userUc      userDom.Usecase
	teamUc      teamDom.Usecase
	challengeUc challengeDom.Usecase
	cache       cache.Cache
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

	submissions, err := p.challengeUc.GetTeamSubmissions(ctx, teamID)
	if err != nil {
		return nil, err
	}

	subDTOs := []*challengeDom.SubmissionDTO{}
	for _, sub := range submissions {
		subDTOs = append(subDTOs, sub.ToDTO())
	}

	score, rank, err := p.getTeamScoreAndRank(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &profileDom.Team{
		Team:        *team,
		Submissions: subDTOs,
		Score:       score,
		Rank:        utils.IntOrNil(rank),
		Members:     team.Members,
	}, nil
}

func (p *profileImpl) getTeamScoreAndRank(ctx context.Context, teamID int64) (rank int, score int, err error) {
	teamRank, err := p.cache.ZRevRank(ctx, scoreDom.ScoreboardCacheKey, scoreDom.ScoreboardMember(teamID))
	if err != nil {
		if !errorDom.Is(err, errorDom.ErrCacheZSetMissingMember) {
			return 0, 0, err
		}
		return 0, 0, nil
	}

	teamScore, err := p.cache.ZScore(ctx, scoreDom.ScoreboardCacheKey, scoreDom.ScoreboardMember(teamID))
	if err != nil {
		if !errorDom.Is(err, errorDom.ErrCacheZSetMissingMember) {
			return 0, 0, err
		}
		return 0, 0, nil
	}

	// need to increase rank by 1 due to zscore handling
	if teamRank > 0 {
		teamRank++
	}

	return int(teamScore), int(teamRank), nil
}

func New(cache cache.Cache, userUc userDom.Usecase, teamUc teamDom.Usecase, challengeUc challengeDom.Usecase) profileDom.Usecase {
	return &profileImpl{
		userUc:      userUc,
		teamUc:      teamUc,
		challengeUc: challengeUc,
		cache:       cache,
	}
}
