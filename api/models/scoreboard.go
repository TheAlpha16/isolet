package models

import "encoding/json"

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
	PageCount int              `json:"page_count"`
	Page      int              `json:"page"`
	Scores    []ScoreBoardTeam `json:"scores"`
}

type ScoreBoardTeam struct {
	TeamID      int64   `gorm:"column:teamid" json:"teamid"`
	TeamName    string  `gorm:"column:teamname" json:"teamname"`
	Rank        int64   `gorm:"column:rank" json:"rank"`
	Score       int64   `gorm:"column:score" json:"score"`
	Submissions []Solve `gorm:"column:submissions;serializer:json" json:"submissions"`
}

type Solve struct {
	Timestamp string `gorm:"column:timestamp" json:"timestamp"`
	Points    int    `gorm:"column:points" json:"points"`
}

func (s ScoreBoardTeam) MarshalJSON() ([]byte, error) {
	type Alias ScoreBoardTeam
	if s.Submissions == nil {
		s.Submissions = []Solve{}
	}
	return json.Marshal((Alias)(s))
}

// type TopScore struct {
// 	TeamID   int64  `gorm:"column:teamid" json:"teamid"`
// 	TeamName string `gorm:"column:teamname" json:"teamname"`
// 	Rank     int64  `gorm:"column:rank" json:"rank"`
// }
