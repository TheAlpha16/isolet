package usecase

import (
	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
	"github.com/TheAlpha16/isolet/api/internal/repository"
	authUc "github.com/TheAlpha16/isolet/api/internal/usecase/auth"
	userUc "github.com/TheAlpha16/isolet/api/internal/usecase/user"
)

type Usecases struct {
	User userDom.Usecase
	Auth authDom.Usecase
}

func New(repos *repository.Repositories) *Usecases {
	user := userUc.New(repos.User)
	auth := authUc.New(repos.Auth)
	return &Usecases{
		User: user,
		Auth: auth,
	}
}
