package jwt

import "context"

type JWT interface {
	Sign(ctx context.Context, claims *Claims) (string, error)
	Verify(ctx context.Context, token string) (*Claims, error)
}
