package challenge

import "context"

type Usecase interface {
	List(ctx context.Context) ([]*Challenge, error)
}

type Repository interface {
	GetAll(ctx context.Context) ([]*Challenge, error)
}
