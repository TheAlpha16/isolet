package user

import "context"

type Usecase interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	ExistsByEmailOrUsername(ctx context.Context, email, username string) (*IdentifierExistence, error)
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	CheckIdentifiers(ctx context.Context, email, username string) (*IdentifierExistence, error)
}
