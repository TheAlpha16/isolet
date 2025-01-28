package models

type Team struct {
	TeamID         int64    `gorm:"primaryKey;autoIncrement;column:teamid" json:"teamid"`
	TeamName       string   `gorm:"unique;not null;column:teamname" json:"teamname"`
	Captain        int64    `gorm:"not null" json:"captain"`
	Members        []User   `gorm:"foreignKey:TeamID;references:TeamID" json:"members"`
	Password       string   `gorm:"not null" json:"-"`
	Cost           int64    `gorm:"not null" json:"-"`
	LastSubmission int64    `gorm:"column:last_submission" json:"-"`
}
