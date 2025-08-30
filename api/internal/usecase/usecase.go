package usecase

import (
	domain_user "github.com/TheAlpha16/isolet/api/internal/domain/user"

	"github.com/TheAlpha16/isolet/api/internal/repository"

	usecase_user "github.com/TheAlpha16/isolet/api/internal/usecase/user"
)

type Usecases struct {
	User domain_user.Usecase
}

func New(repos *repository.Repositories) *Usecases {
	user := usecase_user.New(repos.User)
	return &Usecases{
		User: user,
	}
}
