package validator

import (
	"context"
	"sync"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func Validate(ctx context.Context, input interface{}) error {
	once.Do(func() {
		validate = validator.New()
	})

	err := validate.StructCtx(ctx, input)
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrValidationFailed, "", err, nil)
	}
	return nil
}
