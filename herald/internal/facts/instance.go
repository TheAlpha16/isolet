package facts

import (
	"github.com/TheAlpha16/isolet/herald/utils/errors"
)

const (
	InstanceExpiredFactType FactType = "InstanceExpired"
)

type InstanceFact struct {
	BaseFact
	ID string
}

func (f *InstanceFact) Key() string {
	return f.ID
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
