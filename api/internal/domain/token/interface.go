package token

import "context"

type Usecase interface {
	Create(ctx context.Context, token *Token) error
	Fetch(ctx context.Context, id *TokenIdentifier) (*Token, error)
}
