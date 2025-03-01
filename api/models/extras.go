package models

import (
	"time"

	"github.com/google/uuid"
)

type TeamNameKey struct{}

type Token struct {
	TID    int64     `gorm:"column:tid;primaryKey;autoIncrement" json:"-"`
	Token  uuid.UUID `gorm:"column:token;unique;not null;default:gen_random_uuid()" json:"token"`
	Type   string    `gorm:"column:type;not null;type:token_type" json:"-"`
	UserID int64     `gorm:"column:userid;not null" json:"-"`
	User   User      `gorm:"foreignKey:UserID;references:UserID" json:"-"`
	Expiry time.Time `gorm:"column:expiry" json:"-"`
}

type Config struct {
	Key   string `gorm:"column:key;primaryKey"`
	Value string `gorm:"column:value"`
}

func (Config) TableName() string {
	return "config"
}
