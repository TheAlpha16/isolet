package instance

import "context"

type Usecase interface {
	Start(ctx context.Context, input *StartInput) (*Instance, error)
}

type Repository interface{}
