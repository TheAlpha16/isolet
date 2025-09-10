package infra

import (
	"github.com/TheAlpha16/isolet/api/infra/cnc"
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	"github.com/TheAlpha16/isolet/api/utils"

	"github.com/valkey-io/valkey-go"
)

type Infra struct {
	JWT jwt.JWT
	CNC cnc.CNC
}

func New(valkeyClient valkey.Client) *Infra {
	config := utils.GetConfig()

	return &Infra{
		JWT: jwt.NewJWT(config.Token.SigningKey),
		CNC: cnc.NewCNC(valkeyClient),
	}
}
