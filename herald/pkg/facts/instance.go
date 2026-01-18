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
	val, err := json.Marshal(f)
	if err != nil {
		return nil, errors.Raise(errors.ErrFactMarshalFailed, "failed to marshal InstanceFact", err)
	}
	return val, nil
}

func (f *InstanceFact) Validate() error {
	if err := f.BaseFact.Validate(); err != nil {
		return err
	}

	if f.ID == "" {
		return errors.Raise(errors.ErrFactInvalidKey, "InstanceFact has invalid ID", nil)
	}

	if f.ChallengeID <= 0 {
		return errors.Raise(errors.ErrFactInvalidField, "InstanceFact has invalid ChallengeID", nil)
	}

	if f.TeamID != nil && *f.TeamID <= 0 {
		return errors.Raise(errors.ErrFactInvalidField, "InstanceFact has invalid TeamID", nil)
	}

	return nil
}
