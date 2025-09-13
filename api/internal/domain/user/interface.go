package user

import "context"

type Usecase interface {
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User, fields []string) error
	GetByEmailOrUsername(ctx context.Context, email, username string) (*User, error)

	// ExistsByEmailOrUsername checks whether the given email or username exists.
	// It queries the cache first (email, then username) and returns immediately
	// on the first positive match; otherwise, it falls back to the database.
	ExistsByEmailOrUsername(ctx context.Context, email, username string) error
}

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User, fields []string) error
	GetByEmailOrUsername(ctx context.Context, email, username string) (*User, error)
}
