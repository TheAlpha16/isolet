package models

import (
	"time"

	"github.com/google/uuid"
)

type TeamNameKey struct{}

type Token struct {
	TID       int64     `gorm:"column:tid;primaryKey;autoIncrement" json:"-"`
	Token     uuid.UUID `gorm:"column:token;unique;not null;default:gen_random_uuid()" json:"token"`
	UserID    int64     `gorm:"column:userid;not null" json:"-"`
	User	  User      `gorm:"foreignKey:UserID;references:UserID" json:"-"`
	Timestamp time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"-"`
}
