package usecase

import (
	"context"
	"sync"

	"github.com/TheAlpha16/isolet/api/infra"
	"github.com/TheAlpha16/isolet/api/infra/cache"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	emailDom "github.com/TheAlpha16/isolet/api/internal/domain/email"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/internal/repository"
	authUc "github.com/TheAlpha16/isolet/api/internal/usecase/auth"
	cvUc "github.com/TheAlpha16/isolet/api/internal/usecase/configvars"
	emailUc "github.com/TheAlpha16/isolet/api/internal/usecase/email"
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
}

func New(ctx context.Context, wg *sync.WaitGroup, cache cache.Cache, repos *repository.Repositories, infra *infra.Infra) *Usecases {
	// helpers
	cv := cvUc.New(ctx, repos.ConfigVars, infra.CNC)
	token := tokenUc.New(cache)
	email := emailUc.New(ctx, wg, cv)

	// business
	user := userUc.New(cache, repos.User)
	team := teamUc.New(repos.Team)
	auth := authUc.New(user, token, cv, email, infra.JWT)

	return &Usecases{
		User:       user,
		Team:       team,
		Auth:       auth,
		Token:      token,
		ConfigVars: cv,
		Email:      email,
	}
}
