package user

import "context"

type Usecase interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}
