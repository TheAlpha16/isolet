package user

import "time"

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	TeamID    *int64
	Role      Role
	IsBanned  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
