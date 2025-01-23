package models

type Team struct {
	TeamID         int64  `gorm:"primaryKey;autoIncrement;column:teamid" json:"teamid"`
	TeamName       string `gorm:"unique;not null;column:teamname" json:"teamname"`
	Captain        int64  `gorm:"not null" json:"captain"`
	Members        []User `gorm:"-" json:"members"`
	Password       string `gorm:"not null" json:"password"`
	Cost           int64  `gorm:"not null" json:"cost"`
	LastSubmission int64  `gorm:"column:last_submission" json:"last_submission"`
}
