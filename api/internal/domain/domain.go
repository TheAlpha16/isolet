package domain

import "time"

type BaseEntity struct {
	CreatedAt time.Time `msgpack:"created_at"`
	UpdatedAt time.Time `msgpack:"updated_at"`
}

func (be *BaseEntity) UpdateTime() {
	if be.CreatedAt.IsZero() {
		be.CreatedAt = time.Now()
	}
	be.UpdatedAt = time.Now()
}
