package user

import "time"

type Rank int

const (
	RankAdmin   Rank = 1
	RankCaptain Rank = 2
	RankPlayer  Rank = 3
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	TeamID    *int64
	Rank      Rank
	IsBanned  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
