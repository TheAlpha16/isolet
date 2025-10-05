package external

import (
	instanceExt "github.com/TheAlpha16/isolet/api/external/instance"
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
)

type Services struct {
	Instance instanceDom.Service
}

func New() *Services {
	instance := instanceExt.New()

	return &Services{
		Instance: instance,
	}
}
