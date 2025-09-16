package challenge

import "github.com/TheAlpha16/isolet/api/internal/domain"

type Submission struct {
	ID          int64
	ChallengeID int64
	UserID      int64
	TeamID      int64
	Flag        string
	IsCorrect   bool
	IPAddress   string
	domain.ImmutableEntity
}

type Solve struct {
	ID           int64
	ChallengeID  int64
	TeamID       int64
	SubmissionID int64
	domain.ImmutableEntity
}
