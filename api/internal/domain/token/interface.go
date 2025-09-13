package token

import "context"

type Usecase interface {
	// Create stores a new token in cache.
	// Optionally, additional key-value pairs can be provided via extras.
	// The token and extras are written atomically, and all keys share the same TTL.
	Create(ctx context.Context, token *Token, extras map[string]string) error
	Fetch(ctx context.Context, id *TokenIdentifier) (*Token, error)
	CountAuthTokens(ctx context.Context, userID int64) (int, error)
}
