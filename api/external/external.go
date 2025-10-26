package external

import (
	"context"

	instanceExt "github.com/TheAlpha16/isolet/api/external/instance"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
)

type Services struct {
	Instance instanceDom.Service
}

func New(ctx context.Context) (*Services, error) {
	instance, err := instanceExt.New(ctx)
	if err != nil {
		return nil, err
	}

	return &Services{
		Instance: instance,
	}, nil
}
