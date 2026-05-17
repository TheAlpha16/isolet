package usecase

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/oracle/external"
	"github.com/TheAlpha16/isolet/oracle/infra"
	"github.com/TheAlpha16/isolet/oracle/infra/cache"
	authDom "github.com/TheAlpha16/isolet/oracle/internal/domain/auth"
	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	cvDom "github.com/TheAlpha16/isolet/oracle/internal/domain/configvars"
	emailDom "github.com/TheAlpha16/isolet/oracle/internal/domain/email"
	eventDom "github.com/TheAlpha16/isolet/oracle/internal/domain/event"
	factDom "github.com/TheAlpha16/isolet/oracle/internal/domain/fact"
	instanceDom "github.com/TheAlpha16/isolet/oracle/internal/domain/instance"
	manifestDom "github.com/TheAlpha16/isolet/oracle/internal/domain/manifest"
	profileDom "github.com/TheAlpha16/isolet/oracle/internal/domain/profile"
	scoreDom "github.com/TheAlpha16/isolet/oracle/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/oracle/internal/domain/team"
	tokenDom "github.com/TheAlpha16/isolet/oracle/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/oracle/internal/domain/user"
	"github.com/TheAlpha16/isolet/oracle/internal/repository"
	authUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/auth"
	challengeUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/challenge"
	cvUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/configvars"
	emailUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/email"
	eventUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/event"
	factUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/fact"
	instanceUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/instance"
	manifestUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/manifest"
	profileUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/profile"
	scoreUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/score"
	teamUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/team"
	tokenUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/token"
	userUc "github.com/TheAlpha16/isolet/oracle/internal/usecase/user"
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
	Instance   instanceDom.Usecase
	Fact       factDom.Usecase
	Manifest   manifestDom.Usecase
}

func New(ctx context.Context, wg *sync.WaitGroup, c cache.Cache, repos *repository.Repositories, svc *infra.Infra, ext *external.Services) *Usecases {
	// helpers
	cv := cvUc.New(ctx, repos.ConfigVars, svc.CNC)
	token := tokenUc.New(c)
	email := emailUc.New(ctx, wg, cv)

	// business
	user := userUc.New(c, repos.User, cv)
	auth := authUc.New(user, token, cv, email, svc.JWT)
	team := teamUc.New(repos.Team, user, auth, token, cv, svc.JWT)
	event := eventUc.New(cv)
	score := scoreUc.New(repos.Score, team, c)
	challenge := challengeUc.New(repos.Challenge, repos.Instance, cv, score, wg)
	profile := profileUc.New(c, user, team, challenge)
	manifest := manifestUc.New(repos.Manifest)
	instance := instanceUc.New(ctx, repos.Instance, ext.Instance, c, challenge, manifest, cv, wg)
	fact := factUc.New(instance)

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
		Instance:   instance,
		Fact:       fact,
		Manifest:   manifest,
	}
}
