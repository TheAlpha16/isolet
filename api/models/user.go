package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	UserID   int64  `gorm:"primaryKey;autoIncrement;column:userid" json:"userid"`
	Email    string `gorm:"unique;not null" json:"email" form:"email"`
	Username string `gorm:"unique;not null" json:"username" form:"username"`
	Password string `gorm:"not null" json:"-" form:"password"`
	Rank     int    `gorm:"default:3" json:"rank"`
	TeamID   int64  `gorm:"default:-1;column:teamid" json:"teamid"`
}

type ToVerify struct {
	VID       int64     `gorm:"primaryKey;autoIncrement;column:vid" json:"vid"`
	Email     string    `gorm:"unique;not null" json:"email" form:"email"`
	Username  string    `gorm:"unique;not null" json:"username" form:"username"`
	Password  string    `gorm:"not null" json:"-" form:"password"`
	Confirm   string    `gorm:"-" json:"confirm" form:"confirm"`
	Timestamp time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"timestamp"`
}

type VerifyClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func (ToVerify) TableName() string {
	return "toverify"
}
