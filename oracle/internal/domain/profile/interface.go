package profile

import "context"

type Usecase interface {
	Me(ctx context.Context) (*Me, error)
	Team(ctx context.Context) (*Team, error)
}
