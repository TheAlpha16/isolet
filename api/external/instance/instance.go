package instance

import (
	"context"
	"net/http"

	k8sInfra "github.com/TheAlpha16/isolet/api/infra/k8s"
	"github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"

	tide "github.com/TheAlpha16/isolet/tide/api/v1"
	"github.com/TheAlpha16/isolet/tide/sdk"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type instanceSvc struct {
	client sdk.Handler
}

func (is *instanceSvc) Start(ctx context.Context) error {
	instance, err := is.buildInstance(ctx)
	if err != nil {
		return err
	}

	err = is.client.CreateInstance(ctx, instance)
	if err != nil {
		return errors.Raise(ctx, errors.ErrInstanceCreationFailed, "", err, nil)
	}

	return nil
}

func (is *instanceSvc) Stop(ctx context.Context) error {
	return nil
}

func (is *instanceSvc) Extend(ctx context.Context) error {
	return nil
}

func (is *instanceSvc) buildInstance(ctx context.Context) (*tide.Instance, error) {
	return &tide.Instance{}, nil
}

func New(ctx context.Context) (instanceDom.Service, error) {
	k8sClient, err := k8sInfra.NewK8sClient(
		&http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	)
	if err != nil {
		return nil, errors.Raise(ctx, errors.ErrK8sConnectionFailed, "", err, nil)
	}

	return &instanceSvc{
		client: sdk.New(k8sClient),
	}, nil
}
