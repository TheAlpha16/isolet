package instance

import "github.com/TheAlpha16/isolet/api/internal/domain"

type Instance struct {
	ID          int64
	ChallengeID int64
	TeamID      int64
	ExpiresAt   int64
	domain.BaseEntity
}
