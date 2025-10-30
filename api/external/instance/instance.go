package instance

import (
	"context"

	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"

	"github.com/TheAlpha16/isolet/tide"
)

type instanceSvc struct {
	tide tide.Handler
}

func (is *instanceSvc) Start(ctx context.Context) error {
	return nil
}

func (is *instanceSvc) Stop(ctx context.Context) error {
	return nil
}

func (is *instanceSvc) Extend(ctx context.Context) error {
	return nil
}

func New(ctx context.Context) (instanceDom.Service, error) {
	// client, err := k8sInfra.NewK8sClient()
	// if err != nil {
	// 	return nil, errors.Raise(ctx, errors.ErrK8sConnectionFailed, "", err, nil)
	// }

	return &instanceSvc{}, nil
}
