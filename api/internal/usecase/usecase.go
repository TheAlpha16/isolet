package usecase

import (
	"github.com/TheAlpha16/isolet/api/infra/cache"
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	tokenDom "github.com/TheAlpha16/isolet/api/internal/domain/token"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/internal/repository"
	authUc "github.com/TheAlpha16/isolet/api/internal/usecase/auth"
	tokenUc "github.com/TheAlpha16/isolet/api/internal/usecase/token"
	userUc "github.com/TheAlpha16/isolet/api/internal/usecase/user"
)

type Usecases struct {
	User  userDom.Usecase
	Auth  authDom.Usecase
	Token tokenDom.Usecase
}

func New(cache cache.Cache, repos *repository.Repositories) *Usecases {
	user := userUc.New(repos.User)
	auth := authUc.New(repos.Auth)
	token := tokenUc.New(cache)
	return &Usecases{
		User:  user,
		Auth:  auth,
		Token: token,
	}
}
