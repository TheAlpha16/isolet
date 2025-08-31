package auth

import "context"

type Repository interface{}

type Usecase interface {
	Login(ctx context.Context, input *LoginInput) (*Session, error)
}
