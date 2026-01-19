package domain

import "time"

type BaseEntity struct {
	CreatedAt time.Time `json:"-" msgpack:"created_at"`
	UpdatedAt time.Time `json:"-" msgpack:"updated_at"`
}

func (be *BaseEntity) UpdateTime() {
	if be.CreatedAt.IsZero() {
		be.CreatedAt = time.Now()
	}
	be.UpdatedAt = time.Now()
}

type ImmutableEntity struct {
	CreatedAt time.Time `json:"-" msgpack:"created_at"`
}

func (ie *ImmutableEntity) UpdateTime() {
	if ie.CreatedAt.IsZero() {
		ie.CreatedAt = time.Now()
	}
}
