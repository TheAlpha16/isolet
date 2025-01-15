package models

import (
	"github.com/lib/pq"
)

type Team struct {
	TeamID         int64         `gorm:"primaryKey;autoIncrement;column:teamid" json:"teamid"`
	TeamName       string        `gorm:"unique;not null;column:teamname" json:"teamname"`
	Captain        int64         `gorm:"not null" json:"captain"`
	Members        pq.Int64Array `gorm:"type:integer[]" json:"members"`
	Password       string        `gorm:"not null" json:"password"`
	LastSubmission int64         `gorm:"column:last_submission" json:"last_submission"`
}
