package usecase

import (
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"

	"github.com/TheAlpha16/isolet/api/internal/repository"

	userUc "github.com/TheAlpha16/isolet/api/internal/usecase/user"
)

type Usecases struct {
	User userDom.Usecase
}

func New(repos *repository.Repositories) *Usecases {
	user := userUc.New(repos.User)
	return &Usecases{
		User: user,
	}
}
