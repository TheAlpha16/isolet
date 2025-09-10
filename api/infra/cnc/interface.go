package cnc

import "context"

type CNC interface {
	Register(name string, handler func(ctx context.Context, params map[string]any) error) error
	Trigger(ctx context.Context, name string, params map[string]any) error
}
