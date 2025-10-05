package instance

import (
	"context"

	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
)

type instanceService struct{}

func (s *instanceService) Start(ctx context.Context) error {
	return nil
}

func (s *instanceService) Stop(ctx context.Context) error {
	return nil
}

func (s *instanceService) Extend(ctx context.Context) error {
	return nil
}

func New() instanceDom.Service {
	return &instanceService{}
}
