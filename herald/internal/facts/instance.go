package facts

import (
	"github.com/TheAlpha16/isolet/herald/utils/errors"

	"github.com/goccy/go-json"
)

const (
	InstanceExpiredFactType FactType = "InstanceExpired"
)

type InstanceFact struct {
	BaseFact
	ID          string `json:"id"`
	ChallengeID int64  `json:"challenge_id"`
	TeamID      *int64 `json:"team_id"`
}

func (f *InstanceFact) Key() []byte {
	return []byte(f.ID)
}

func (f *InstanceFact) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *InstanceFact) Validate() error {
	if err := f.BaseFact.Validate(); err != nil {
		return err
	}

	if f.ID == "" {
		return errors.Raise(errors.ErrFactInvalidKey, "InstanceFact has invalid ID", nil)
	}

	return nil
}
