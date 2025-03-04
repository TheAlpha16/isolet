package models

type Config struct {
	Key   string      `gorm:"primaryKey;not null;unique" json:"key" form:"key"`
	Value string      `gorm:"not null" json:"value" form:"value"`
}
