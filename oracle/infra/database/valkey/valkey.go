package valkey

import (
	"context"
	"crypto/tls"

	"github.com/TheAlpha16/isolet/oracle/utils"

	"github.com/valkey-io/valkey-go"
)

func NewValkey(ctx context.Context) (valkey.Client, error) {
	config := utils.GetConfig()
	opts, err := valkey.ParseURL(config.Valkey.Address[0])
	if err != nil {
		return nil, err
	}
	opts.ClientName = config.Name
	opts.Username = config.Valkey.Username
	opts.Password = config.Valkey.Password

	if config.Valkey.UseTLS {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, err
	}

	err = client.Do(ctx, client.B().Ping().Build()).Error()
	if err != nil {
		return nil, err
	}

	return client, nil
}
