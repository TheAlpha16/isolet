package user

import (
	"github.com/TheAlpha16/isolet/api/utils/validator"

	govalidator "github.com/go-playground/validator/v10"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleAuthor  Role = "author"
	RoleCaptain Role = "captain"
	RolePlayer  Role = "player"
)

var validRoles = map[Role]struct{}{
	RoleAdmin:   {},
	RoleAuthor:  {},
	RoleCaptain: {},
	RolePlayer:  {},
}

func (r Role) IsValid() bool {
	_, ok := validRoles[r]
	return ok
}

func roleValidator(fl govalidator.FieldLevel) bool {
	return Role(fl.Field().String()).IsValid()
}

func init() {
	validator.RegisterValidation("role", roleValidator)
}
