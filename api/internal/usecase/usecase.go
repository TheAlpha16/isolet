package usecase

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/api/infra"
	"github.com/TheAlpha16/isolet/api/infra/cache"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	challengeDom "github.com/TheAlpha16/isolet/api/internal/domain/challenge"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	emailDom "github.com/TheAlpha16/isolet/api/internal/domain/email"
	eventDom "github.com/TheAlpha16/isolet/api/internal/domain/event"
	profileDom "github.com/TheAlpha16/isolet/api/internal/domain/profile"
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/internal/repository"
	authUc "github.com/TheAlpha16/isolet/api/internal/usecase/auth"
	challengeUc "github.com/TheAlpha16/isolet/api/internal/usecase/challenge"
	cvUc "github.com/TheAlpha16/isolet/api/internal/usecase/configvars"
	emailUc "github.com/TheAlpha16/isolet/api/internal/usecase/email"
	eventUc "github.com/TheAlpha16/isolet/api/internal/usecase/event"
	profileUc "github.com/TheAlpha16/isolet/api/internal/usecase/profile"
	scoreUc "github.com/TheAlpha16/isolet/api/internal/usecase/score"
	teamUc "github.com/TheAlpha16/isolet/api/internal/usecase/team"
	tokenUc "github.com/TheAlpha16/isolet/api/internal/usecase/token"
	userUc "github.com/TheAlpha16/isolet/api/internal/usecase/user"
)

type Usecases struct {
	User       userDom.Usecase
	Team       teamDom.Usecase
	Auth       authDom.Usecase
	Token      tokenDom.Usecase
	ConfigVars cvDom.Usecase
	Email      emailDom.Usecase
	Event      eventDom.Usecase
	Profile    profileDom.Usecase
	Challenge  challengeDom.Usecase
	Score      scoreDom.Usecase
}

func New(ctx context.Context, wg *sync.WaitGroup, cache cache.Cache, repos *repository.Repositories, infra *infra.Infra) *Usecases {
	// helpers
	cv := cvUc.New(ctx, repos.ConfigVars, infra.CNC)
	token := tokenUc.New(cache)
	email := emailUc.New(ctx, wg, cv)

	// business
	user := userUc.New(cache, repos.User, cv)
	auth := authUc.New(user, token, cv, email, infra.JWT)
	team := teamUc.New(repos.Team, user, auth, token, cv, infra.JWT)
	event := eventUc.New(cv)
	score := scoreUc.New(repos.Score, team, cache)
	challenge := challengeUc.New(repos.Challenge, cv, score, wg)
	profile := profileUc.New(user, team, challenge)

	return &Usecases{
		User:       user,
		Team:       team,
		Auth:       auth,
		Token:      token,
		ConfigVars: cv,
		Email:      email,
		Event:      event,
		Profile:    profile,
		Challenge:  challenge,
		Score:      score,
	}
}
