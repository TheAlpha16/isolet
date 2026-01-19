package user

import "context"

type Usecase interface {
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User, fields []string) error
	GetByEmailOrUsername(ctx context.Context, email, username string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	ExistsByEmailOrUsername(ctx context.Context, email, username string) error
	JoinTeam(ctx context.Context, userID int64, teamID int64) error
}

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User, fields []string) error
	GetByEmailOrUsername(ctx context.Context, email, username string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	JoinTeam(ctx context.Context, userID int64, teamID int64, userLimit int) error
}
