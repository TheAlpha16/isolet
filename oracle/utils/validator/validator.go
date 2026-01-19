package validator

import (
	"context"
	"fmt"
	"sync"

	errorDom "github.com/TheAlpha16/isolet/oracle/internal/domain/errors"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func buildErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", fe.Field(), fe.Param())
	case "token_purpose":
		return fmt.Sprintf("%s (%v) is not a valid token purpose", fe.Field(), fe.Value())
	case "role":
		return fmt.Sprintf("%s (%v) is not a valid role", fe.Field(), fe.Value())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

func Validate(ctx context.Context, input interface{}) error {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
	})

	err := validate.StructCtx(ctx, input)
	if err != nil {
		if verrs, ok := err.(validator.ValidationErrors); ok && len(verrs) > 0 {
			return errorDom.Raise(ctx, errorDom.ErrValidationFailed, buildErrorMessage(verrs[0]), nil, nil)
		}
		return errorDom.Raise(ctx, errorDom.ErrValidationFailed, "", err, nil)
	}
	return nil
}

func RegisterValidation(tag string, fn validator.Func) error {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
	})
	return validate.RegisterValidation(tag, fn)
}
