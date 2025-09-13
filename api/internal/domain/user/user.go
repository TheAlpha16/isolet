package user

import (
	"fmt"

	"github.com/TheAlpha16/isolet/api/internal/domain"
)

const (
	emailCachePrefix    = "exists:email"
	usernameCachePrefix = "exists:username"
)

type User struct {
	ID       int64
	Username string
	Email    string
	Password string
	TeamID   *int64
	Role     Role
	IsBanned bool
	domain.BaseEntity
}

func EmailCacheKey(email string) string {
	return fmt.Sprintf("%s:%s", emailCachePrefix, email)
}

func UsernameCacheKey(username string) string {
	return fmt.Sprintf("%s:%s", usernameCachePrefix, username)
}
