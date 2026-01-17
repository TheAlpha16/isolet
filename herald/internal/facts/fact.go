package facts

import (
	"fmt"
	"time"

	"github.com/TheAlpha16/isolet/herald/utils/errors"
)

type FactType string

type Fact interface {
	FactType() FactType
	Key() string
	OccuredAt() time.Time
	Validate() error
}

type BaseFact struct {
	Type FactType
	At   time.Time
}

func (f *BaseFact) FactType() FactType {
	return f.Type
}

func (f *BaseFact) OccuredAt() time.Time {
	return f.At
}

func (f *BaseFact) Validate() error {
	if _, ok := allowedTypes[f.FactType()]; !ok {
		return errors.Raise(errors.ErrFactInvalidType, fmt.Sprintf("fact has invalid type: %s", f.FactType()), nil)
	}

	if f.OccuredAt().IsZero() {
		return errors.Raise(errors.ErrFactInvalidOccuredAt, "", nil)
	}

	return nil
}

var allowedTypes = map[FactType]struct{}{
	InstanceExpiredFactType: {},
}
