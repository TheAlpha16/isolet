package cnc

import (
	"context"

	"github.com/TheAlpha16/isolet/api/utils"

	cncgo "github.com/TheAlpha16/cnc-go"
	"github.com/valkey-io/valkey-go"
)

type cncImpl struct {
	client cncgo.CNC
}

func (c *cncImpl) Trigger(ctx context.Context, name string, params map[string]any) error {
	cmd := cncgo.Command{
		Name:       cncgo.CommandName(name),
		Parameters: params,
	}
	return c.client.TriggerCommand(ctx, cmd)
}

func (c *cncImpl) Register(name string, handler func(ctx context.Context, params map[string]any) error) error {
	wrapped := func(ctx context.Context, cmd cncgo.Command) error {
		return handler(ctx, cmd.Parameters)
	}
	return c.client.RegisterHandler(cncgo.CommandName(name), wrapped)
}

func NewCNC(ctx context.Context, cache valkey.Client) (CNC, error) {
	config := utils.GetConfig()
	client := cncgo.NewCNCWithValkey(cache, config.CNC.Channel)
	cnc := &cncImpl{client: client}

	if err := client.Start(ctx); err != nil {
		return nil, err
	}
	return cnc, nil
}
