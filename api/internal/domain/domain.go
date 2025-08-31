package domain

import "time"

type BaseEntity struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (be *BaseEntity) UpdateTime() {
	if be.CreatedAt.IsZero() {
		be.CreatedAt = time.Now()
	}
	be.UpdatedAt = time.Now()
}
