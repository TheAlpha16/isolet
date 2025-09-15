package profile

import "context"

type Usecase interface {
	Me(ctx context.Context) (*Me, error)
}
