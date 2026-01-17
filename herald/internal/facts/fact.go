package facts

import "time"

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
