package usecase

import (
	"github.com/TheAlpha16/isolet/api/infra"
	"github.com/TheAlpha16/isolet/api/infra/cache"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/internal/repository"
	authUc "github.com/TheAlpha16/isolet/api/internal/usecase/auth"
	cvUc "github.com/TheAlpha16/isolet/api/internal/usecase/configvars"
	tokenUc "github.com/TheAlpha16/isolet/api/internal/usecase/token"
	userUc "github.com/TheAlpha16/isolet/api/internal/usecase/user"
)

type Usecases struct {
	User       userDom.Usecase
	Auth       authDom.Usecase
	Token      tokenDom.Usecase
	ConfigVars cvDom.Usecase
}

func New(cache cache.Cache, repos *repository.Repositories, infra *infra.Infra) *Usecases {
	token := tokenUc.New(cache)
	user := userUc.New(cache, repos.User)
	cv := cvUc.New(repos.ConfigVars)
	auth := authUc.New(repos.Auth, user, token, cv, infra.JWT)

	return &Usecases{
		User:       user,
		Auth:       auth,
		Token:      token,
		ConfigVars: cv,
	}
}
