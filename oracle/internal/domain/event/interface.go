package event

import "context"

type Usecase interface {
	Info(ctx context.Context) (*InfoOutput, error)
}
