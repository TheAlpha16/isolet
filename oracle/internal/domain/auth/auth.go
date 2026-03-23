package auth

type Session struct {
	UserID        int64  `json:"user_id"`
	TeamID        *int64 `json:"team_id,omitempty"`
	Token         string `json:"-"`
	RealtimeToken string `json:"-"`
	ExpiresAt     int64  `json:"expires_at"`
}
