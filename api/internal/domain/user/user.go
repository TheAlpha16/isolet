package user

import "time"

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleAuthor  Role = "author"
	RoleCaptain Role = "captain"
	RolePlayer  Role = "player"
)

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

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleAuthor, RoleCaptain, RolePlayer:
		return true
	}
	return false
}
