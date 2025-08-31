package user

import (
	"github.com/TheAlpha16/isolet/api/internal/domain"
)

type User struct {
	ID         int64
	Username   string
	Email      string
	Password   string
	TeamID     *int64
	Role       Role
	IsVerified bool
	IsBanned   bool
	domain.BaseEntity
}
