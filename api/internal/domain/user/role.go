package user

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleAuthor  Role = "author"
	RoleCaptain Role = "captain"
	RolePlayer  Role = "player"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleAuthor, RoleCaptain, RolePlayer:
		return true
	}
	return false
}
