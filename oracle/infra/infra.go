package infra

import (
	"context"

	"github.com/TheAlpha16/isolet/oracle/infra/cnc"
	"github.com/TheAlpha16/isolet/oracle/infra/jwt"
	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/valkey-io/valkey-go"
)

type Infra struct {
	JWT jwt.JWT
	CNC cnc.CNC
}

func New(ctx context.Context, valkeyClient valkey.Client) (*Infra, error) {
	config := utils.GetConfig()
	cnc, err := cnc.NewCNC(ctx, valkeyClient)
	if err != nil {
		return nil, err
	}

	return &Infra{
		JWT: jwt.NewJWT(config.Token.SigningKey),
		CNC: cnc,
	}, nil
}
