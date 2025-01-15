package models

import (
	"time"

	"github.com/lib/pq"
)

type Challenge struct {
	ChallID      int            `gorm:"primaryKey;column:chall_id" json:"chall_id"`
	Name         string         `gorm:"not null;unique;column:chall_name" json:"name"`
	Prompt       string         `gorm:"type:text" json:"prompt"`
	Category     Category       `gorm:"foreignKey:CategoryID;references:CategoryID" json:"-"`
	CategoryID   int            `gorm:"not null;column:category_id" json:"-"`
	Flag         string         `gorm:"type:text" json:"-"`
	Type         string         `gorm:"type:chall_type;default:static" json:"type"`
	Points       int            `gorm:"not null;default:100" json:"points"`
	Files        pq.StringArray `gorm:"type:text[]" json:"files"`
	Requirements pq.Int64Array  `gorm:"type:integer[]" json:"-"`
	Hints        []Hint         `gorm:"foreignKey:ChallID" json:"hints"`
	Solves       int            `gorm:"default:0" json:"solves"`
	Author       string         `gorm:"default:anonymous" json:"author"`
	Visible      bool           `gorm:"default:false" json:"-"`
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
	Links        pq.StringArray `gorm:"type:text[]" json:"links"`
	Done         bool           `gorm:"column:done" json:"done"`
}

type ChallengeData struct {
	ChallID      int            `gorm:"column:chall_id" json:"chall_id"`
	Name         string         `gorm:"column:chall_name" json:"name"`
	Prompt       string         `gorm:"type:text" json:"prompt"`
	Type         string         `gorm:"type:chall_type" json:"type"`
	Points       int            `gorm:"column:points" json:"points"`
	Files        pq.StringArray `gorm:"column:files;type:text[]" json:"files"`
	Hints        []Hint         `gorm:"serializer:json" json:"hints"`
	Solves       int            `gorm:"column:solves" json:"solves"`
	Author       string         `gorm:"column:author" json:"author"`
	Tags         pq.StringArray `gorm:"column:tags;type:text[]" json:"tags"`
	Links        pq.StringArray `gorm:"column:links;type:text[]" json:"links"`
	CategoryName string         `gorm:"column:category_name" json:"-"`
	Deployment   string         `gorm:"type:deployment_type" json:"-"`
	Port         int            `gorm:"column:port" json:"-"`
	Subd         string         `gorm:"column:subd" json:"-"`
	Done         bool           `gorm:"column:done" json:"done"`
}

type Hint struct {
	HID      int    `gorm:"primaryKey;column:hid" json:"hid"`
	ChallID  int    `gorm:"not null;column:chall_id" json:"-"`
	Hint     string `gorm:"not null" json:"hint"`
	Cost     int    `gorm:"not null;default:0" json:"cost"`
	Visible  bool   `gorm:"default:false" json:"-"`
	Unlocked bool   `gorm:"-" json:"unlocked"`
}

type Category struct {
	CategoryID   int    `gorm:"primaryKey;column:category_id" json:"category_id"`
	CategoryName string `gorm:"not null;unique" json:"category_name"`
}

type UHint struct {
	HID       int       `gorm:"column:hid" json:"hid"`
	TeamID    int64     `gorm:"column:teamid" json:"teamid"`
	Timestamp time.Time `gorm:"column:timestamp" json:"timestamp"`
}

func (UHint) TableName() string {
	return "uhints"
}

func (ChallengeData) TableName() string {
	return "get_challenges"
}
