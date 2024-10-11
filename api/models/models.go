package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

type User struct {
	UserID   int64  `gorm:"primaryKey;autoIncrement;column:userid" json:"userid"`
	Email    string `gorm:"unique;not null" json:"email" form:"email"`
	Username string `gorm:"unique;not null" json:"username" form:"username"`
	Password string `gorm:"not null" json:"password" form:"password"`
	Rank     int    `gorm:"default:3" json:"rank"`
	TeamID   int64  `gorm:"default:-1;column:teamid" json:"teamid"`
}

type ToVerify struct {
	VID       int64     `gorm:"primaryKey;autoIncrement;column:vid" json:"vid"`
	Email     string    `gorm:"unique;not null" json:"email" form:"email"`
	Username  string    `gorm:"unique;not null" json:"username" form:"username"`
	Password  string    `gorm:"not null" json:"password" form:"password"`
	Confirm   string    `gorm:"-" json:"confirm" form:"confirm"`
	Timestamp time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"timestamp"`
}

type VerifyClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

type Challenge struct {
	ChallID      int            `gorm:"primaryKey;column:chall_id" json:"chall_id"`
	Name         string         `gorm:"not null;unique;column:chall_name" json:"name"`
	Prompt       string         `gorm:"type:text" json:"prompt"`
	Category     Category       `gorm:"foreignKey:CategoryID;references:CategoryID" json:"-"`
	CategoryID   int            `gorm:"not null;column:category_id" json:"-"`
	Flag         string         `gorm:"-" json:"-"`
	Type         string         `gorm:"type:chall_type;default:static" json:"type"`
	Points       int            `gorm:"not null;default:100" json:"points"`
	Files        pq.StringArray `gorm:"type:text[]" json:"files"`
	Requirements pq.StringArray `gorm:"type:text[]" json:"-"`
	Hints        []Hint         `gorm:"foreignKey:ChallID" json:"hints"`
	Solves       int            `gorm:"default:0" json:"solves"`
	Author       string         `gorm:"default:anonymous" json:"author"`
	Visible      bool           `gorm:"default:false" json:"-"`
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
	Links        pq.StringArray `gorm:"type:text[]" json:"links"`
	Done         bool           `gorm:"-" json:"done"`
}

type Image struct {
	IID        int    `gorm:"primaryKey;autoIncrement;column:iid" json:"iid"`
	ChallID    int    `gorm:"not null;column:chall_id" json:"chall_id"`
	Registry   string `gorm:"not null" json:"registry"`
	Image      string `gorm:"not null" json:"image"`
	Deployment string `gorm:"type:deployment_type;default:http" json:"deployment"`
	Port       int    `gorm:"default:80" json:"port"`
	Subd       string `gorm:"default:localhost" json:"subd"`
	CPU        int    `gorm:"default:5;column:cpu" json:"-"`
	Memory     int    `gorm:"default:10;column:mem" json:"-"`
}

type Hint struct {
	HID      int    `gorm:"primaryKey;column:hid" json:"hid"`
	ChallID  int    `gorm:"not null;column:chall_id" json:"chall_id"`
	Hint     string `gorm:"not null" json:"hint"`
	Cost     int    `gorm:"not null;default:0" json:"cost"`
	Visible  bool   `gorm:"default:false" json:"-"`
	Unlocked bool   `gorm:"-" json:"unlocked"`
}

type Category struct {
	CategoryID   int    `gorm:"primaryKey;column:category_id" json:"category_id"`
	CategoryName string `gorm:"not null;unique" json:"category_name"`
}

type Team struct {
	TeamID   int64          `gorm:"primaryKey;autoIncrement;column:teamid" json:"teamid"`
	TeamName string         `gorm:"unique;not null;column:teamname" json:"teamname"`
	Captain  int64          `gorm:"not null" json:"captain"`
	Members  pq.Int64Array  `gorm:"type:integer[]" json:"members"`
	Password string         `gorm:"not null" json:"password"`
	Solved   pq.StringArray `gorm:"type:text[]" json:"solved"`
	UHints   pq.Int64Array  `gorm:"type:integer[];column:uhints" json:"uhints"`
}

type Flag struct {
	FlagID   int64  `gorm:"primaryKey;autoIncrement;column:flagid" json:"-"`
	TeamID   int64  `gorm:"not null;column:teamid" json:"-"`
	ChallID  int    `gorm:"not null;column:chall_id" json:"chall_id"`
	Flag     string `gorm:"not null" json:"-"`
	Password string `gorm:"type:text" json:"password"`
	Port     int    `gorm:"type:integer" json:"port"`
	Hostname string `gorm:"type:text" json:"hostname"`
	Deadline int64  `gorm:"type:integer" json:"deadline"`
	Extended int    `gorm:"type:integer;default:1" json:"-"`
}

type Score struct {
	Username string `json:"username"`
	Score    string `json:"score"`
}

type AccessDetails struct {
	Password string `json:"password"`
	Port     int32  `json:"port"`
	Hostname string `json:"hostname"`
	Deadline int64  `json:"deadline"`
}

type ExtendDeadline struct {
	Deadline int64 `json:"deadline"`
}

func (ToVerify) TableName() string {
	return "toverify"
}
