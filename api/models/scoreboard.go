package models

import (
	"time"
)

type Sublog struct {
	SID     int64     `gorm:"primaryKey;autoIncrement;column:sid" json:"-"`
	ChallID int       `gorm:"not null;column:chall_id" json:"chall_id"`
	UserID  int64     `gorm:"not null;column:userid" json:"userid"`
	TeamID  int64     `gorm:"not null;column:teamid" json:"teamid"`
	Flag    string    `gorm:"not null" json:"flag"`
	Correct bool      `gorm:"not null" json:"correct"`
	IP      string    `gorm:"not null;column:ip" json:"-"`
	SubTime time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;column:subtime" json:"subtime"`
}

type Score struct {
	Rank     int64  `gorm:"column:rank" json:"rank"`
	TeamID   int64  `gorm:"column:teamid" json:"teamid"`
	TeamName string `gorm:"column:teamname" json:"teamname"`
	Score    int    `gorm:"column:score" json:"score"`
}

type ScoreBoard struct {
	PageCount int     `json:"page_count"`
	Page      int     `json:"page"`
	Scores    []Score `json:"scores"`
}

type Solve struct {
	ChallID   int       `gorm:"column:chall_id" json:"chall_id"`
	TeamID    int64     `gorm:"column:teamid" json:"teamid"`
	Timestamp time.Time `gorm:"column:timestamp" json:"timestamp"`
}

type PointData struct {
	Points    int    `gorm:"column:points" json:"points"`
	Timestamp string `gorm:"column:timestamp" json:"timestamp"`
}

type ScoreGraph struct {
	TeamName    string      `gorm:"column:teamname" json:"teamname"`
	TeamID      int64       `gorm:"column:teamid" json:"teamid"`
	Rank        int         `gorm:"column:rank" json:"rank"`
	Submissions []PointData `gorm:"column:submissions;serializer:json" json:"submissions"`
}
