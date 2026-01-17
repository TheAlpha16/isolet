package facts

import (
	"github.com/TheAlpha16/isolet/herald/utils/errors"

	"github.com/google/uuid"
)

const (
	InstanceExpiredFactType FactType = "InstanceExpired"
)

type InstanceFact struct {
	BaseFact
	ID uuid.UUID
}

func (f *InstanceFact) Key() string {
	return f.ID.String()
}

func (f *InstanceFact) Validate() error {
	if err := f.BaseFact.Validate(); err != nil {
		return err
	}

	if f.ID == uuid.Nil {
		return errors.Raise(errors.ErrFactInvalidKey, "InstanceFact has invalid ID", nil)
	}

	return nil
}
