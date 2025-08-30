package user

type Rank int

const (
	RankAdmin   Rank = 1
	RankCaptain Rank = 2
	RankPlayer  Rank = 3
)

type User struct {
	ID       int
	Username string
	Email    string
	Password string
	TeamID   *int
	Rank     Rank
	IsBanned bool
}
