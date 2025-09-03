package infra

import (
	"github.com/TheAlpha16/isolet/api/infra/jwt"
	"github.com/TheAlpha16/isolet/api/utils"
)

type Infra struct {
	JWT jwt.JWT
}

func New() *Infra {
	config := utils.GetConfig()

	return &Infra{
		JWT: jwt.NewJWT(config.Token.SigningKey),
	}
}
