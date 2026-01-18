package consumer

import (
	"context"

	errorDom "github.com/TheAlpha16/isolet/api/internal/domain/errors"
	factDom "github.com/TheAlpha16/isolet/api/internal/domain/fact"

	"github.com/goccy/go-json"
)

func deserializeInstanceFact(ctx context.Context, data []byte) (factDom.Fact, error) {
	var f factDom.InstanceFact
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, errorDom.Raise(ctx, errorDom.ErrFactDeserialization, "failed to deserialize instance fact", err, nil)
	}
	return &f, nil
}
