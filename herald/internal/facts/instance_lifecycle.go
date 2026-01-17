package facts

import (
	"github.com/TheAlpha16/isolet/herald/utils/errors"

	"github.com/google/uuid"
)

const (
	InstanceExpiredFactType FactType = "InstanceExpired"
)

type InstanceLifecycleFact struct {
	BaseFact
	ID uuid.UUID
}

func (f *InstanceLifecycleFact) Key() string {
	return f.ID.String()
}

func (f *InstanceLifecycleFact) Validate() error {
	if err := f.BaseFact.Validate(); err != nil {
		return err
	}

	if f.ID == uuid.Nil {
		return errors.Raise(errors.ErrFactInvalidKey, "InstanceLifecycleFact has invalid ID", nil)
	}

	return nil
}
