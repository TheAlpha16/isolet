package team

import "context"

type Usecase interface {
	Create(ctx context.Context, input *CreateInput) error
}

type Repository interface{}
