package instance

import (
	"context"

	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
)

type instanceSvc struct{}

func (is *instanceSvc) Start(ctx context.Context) error {
	return nil
}

func (is *instanceSvc) Stop(ctx context.Context) error {
	return nil
}

func (is *instanceSvc) Extend(ctx context.Context) error {
	return nil
}

func New() instanceDom.Service {
	return &instanceSvc{}
}
