package models

type Sublog struct {
	SID       int64  `gorm:"primaryKey;autoIncrement;column:sid" json:"-"`
	ChallID   int    `gorm:"not null;column:chall_id" json:"chall_id"`
	UserID    int64  `gorm:"not null;column:userid" json:"userid"`
	TeamID    int64  `gorm:"not null;column:teamid" json:"teamid"`
	Flag      string `gorm:"not null" json:"-"`
	Correct   bool   `gorm:"not null" json:"correct"`
	IP        string `gorm:"not null;column:ip" json:"-"`
	Timestamp string `gorm:"not null;column:timestamp" json:"timestamp"`
	Points    int    `gorm:"column:points" json:"points"`
}

type ScoreBoard struct {
	PageCount int    `json:"page_count"`
	Page      int    `json:"page"`
	Scores    []Team `json:"scores"`
}
