package instance

import (
	"context"

	k8sInfra "github.com/TheAlpha16/isolet/api/infra/k8s"
	"github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"

	"k8s.io/client-go/kubernetes"
)

type instanceSvc struct {
	client *kubernetes.Clientset
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
	client, err := k8sInfra.NewK8sClient()
	if err != nil {
		return nil, errors.Raise(ctx, errors.ErrK8sConnectionFailed, "", err, nil)
	}

	return &instanceSvc{
		client: client,
	}, nil
}
