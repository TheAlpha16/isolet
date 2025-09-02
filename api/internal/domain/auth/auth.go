package auth

type Session struct {
	UserID    int64  `json:"user_id"`
	Token     string `json:"-"`
	ExpiresAt int64  `json:"expires_at"`
}
