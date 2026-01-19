package user

import (
	"fmt"

	"github.com/TheAlpha16/isolet/oracle/internal/domain"
)

const (
	emailCachePrefix    = "exists:email"
	usernameCachePrefix = "exists:username"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	TeamID   *int64 `json:"-"`
	Role     Role   `json:"role"`
	IsBanned bool   `json:"-"`
	domain.BaseEntity
}

func EmailCacheKey(email string) string {
	return fmt.Sprintf("%s:%s", emailCachePrefix, email)
}

func UsernameCacheKey(username string) string {
	return fmt.Sprintf("%s:%s", usernameCachePrefix, username)
}
