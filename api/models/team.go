package models

type Team struct {
	TeamID         int64    `gorm:"primaryKey;autoIncrement;column:teamid" json:"teamid"`
	TeamName       string   `gorm:"unique;not null;column:teamname" json:"teamname"`
	Captain        int64    `gorm:"not null" json:"captain"`
	Members        []User   `gorm:"column:members" json:"members"`
	Password       string   `gorm:"not null" json:"-"`
	Cost           int64    `gorm:"not null" json:"cost"`
	LastSubmission int64    `gorm:"column:last_submission" json:"last_submission"`
	Score          int64    `gorm:"column:score" json:"score"`
	Rank           int64    `gorm:"column:rank" json:"rank"`
	Submissions    []Sublog `gorm:"column:submissions;serializer:json" json:"submissions"`
}
