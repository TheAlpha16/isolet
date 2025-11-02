package instance

import (
	"context"
	"net/http"

	k8sInfra "github.com/TheAlpha16/isolet/api/infra/k8s"
	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"

	"github.com/TheAlpha16/isolet/tide/sdk"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type instanceSvc struct {
	client sdk.Handler
}

func (is *instanceSvc) Start(ctx context.Context, inst *instanceDom.Instance) error {
	instance := toTideInstance(ctx, inst)

	err := is.client.CreateInstance(ctx, instance)
	if err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceCreationFailed, "", err, nil)
	}

	// TODO wait for the instance and populate the necessary fields

	return nil
}

func (is *instanceSvc) Stop(ctx context.Context, inst *instanceDom.Instance) error {
	return nil
}

func (is *instanceSvc) Extend(ctx context.Context, inst *instanceDom.Instance) error {
	instance := toTideInstance(ctx, inst)

	if err := is.client.UpdateInstance(ctx, instance); err != nil {
		return errorDom.Raise(ctx, errorDom.ErrInstanceUpdateFailed, "failed to extend instance expiry", err, nil)
	}

	return nil
}

func New(ctx context.Context) (instanceDom.Service, error) {
	k8sClient, err := k8sInfra.NewK8sClient(
		&http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	)
	if err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrK8sConnectionFailed, "", err, nil)
	}

	return &instanceSvc{
		client: sdk.New(k8sClient),
	}, nil
}
