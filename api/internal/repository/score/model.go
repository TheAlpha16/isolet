package score

type SubmissionStatsDTO struct {
	ChallengeID    int64 `gorm:"column:challenge_id"`
	CorrectCount   int64 `gorm:"column:correct_count"`
	IncorrectCount int64 `gorm:"column:incorrect_count"`
}

type ChallengeSolveCountDTO struct {
	ChallengeID int64 `gorm:"column:challenge_id"`
	SolveCount  int64 `gorm:"column:solve_count"`
}

type ScoreboardRow struct {
	TeamID int64 `gorm:"column:team_id"`
	Score  int64 `gorm:"column:score"`
}
