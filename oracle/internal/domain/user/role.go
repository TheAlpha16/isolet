package user

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"
	"github.com/TheAlpha16/isolet/oracle/utils/logger"
	"github.com/TheAlpha16/isolet/oracle/utils/validator"

	govalidator "github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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
	if err := validator.RegisterValidation("role", roleValidator); err != nil {
		errorDom.RaiseToSentry(context.Background(), err)
		logger.GetAppLogger().Fatal("failed to register role validator", zap.Error(err))
	}
}
