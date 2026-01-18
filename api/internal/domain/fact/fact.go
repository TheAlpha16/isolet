package fact

import (
	"fmt"
	"time"
)

const (
	InstanceExpiredFactType FactType = "InstanceExpired"
)

type FactType string

type Fact interface {
	FactType() FactType
	Key() []byte
	OccuredAt() time.Time
	Validate() error
}

type BaseFact struct {
	Type FactType  `json:"type"`
	At   time.Time `json:"at"`
}

func (f *BaseFact) FactType() FactType {
	return f.Type
}

func (f *BaseFact) OccuredAt() time.Time {
	return f.At
}

func (f *BaseFact) Validate() error {
	if f.Type == "" {
		return fmt.Errorf("fact type is required")
	}
	if f.At.IsZero() {
		return fmt.Errorf("fact occured at is required")
	}
	return nil
}

type InstanceFact struct {
	BaseFact
	ID          string `json:"id"`
	ChallengeID int64  `json:"challenge_id"`
	TeamID      *int64 `json:"team_id"`
}

func (f *InstanceFact) Key() []byte {
	return []byte(f.ID)
}

func (f *InstanceFact) Validate() error {
	if err := f.BaseFact.Validate(); err != nil {
		return err
	}
	if f.ID == "" {
		return fmt.Errorf("instance fact ID is required")
	}
	return nil
}
