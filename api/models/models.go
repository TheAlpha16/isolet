package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	UserID   uint   `gorm:"primaryKey;autoIncrement" json:"userid"`
	Email    string `gorm:"unique;not null" json:"email" form:"email"`
	Username string `gorm:"unique;not null" json:"username" form:"username"`
	Password string `gorm:"not null" json:"password" form:"password"`
	Rank     int    `gorm:"default:3" json:"rank"`
	TeamID   int    `gorm:"default:-1" json:"teamid"`
}

type ToVerify struct {
	VID       uint      `gorm:"primaryKey;autoIncrement" json:"vid"`
	Email     string    `gorm:"unique;not null" json:"email" form:"email"`
	Username  string    `gorm:"unique;not null" json:"username" form:"username"`
	Password  string    `gorm:"not null" json:"password" form:"password"`
	Confirm  string `gorm:"-" json:"confirm" form:"confirm"`
	Timestamp time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"timestamp"`
}

type VerifyClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

type Challenge struct {
	ChallID   int    `gorm:"primaryKey" json:"challid"`
	Level     int    `gorm:"not null" json:"level"`
	Name      string `gorm:"not null;unique" json:"name"`
	Prompt    string `gorm:"-" json:"prompt"`
	Category  Category `gorm:"foreignKey:CategoryID" json:"category"`
	CategoryID int    `gorm:"not null" json:"-"`
	Flag      string `gorm:"-" json:"flag"`
	Type      string `gorm:"type:chall_type;default:static" json:"type"`
	Points    int    `gorm:"not null;default:100" json:"points"`
	Files     []string `gorm:"type:text[]" json:"files"`
	Hints     []Hint `gorm:"foreignKey:ChallID" json:"hints"`
	Solves    int    `gorm:"default:0" json:"solves"`
	Author    string `gorm:"default:anonymous" json:"author"`
	Visible   bool   `gorm:"default:false" json:"visible"`
	Tags      []string `gorm:"type:text[]" json:"tags"`
	Port      int    `gorm:"default:0" json:"port"`
	Subd      string `gorm:"default:localhost" json:"subd"`
	CPU       int    `gorm:"default:5" json:"cpu"`
	Memory    int    `gorm:"default:10" json:"memory"`
}

type Hint struct {
	HID     int    `gorm:"primaryKey" json:"hid"`
	ChallID int    `gorm:"not null" json:"challid"`
	Hint    string `gorm:"not null" json:"hint"`
	Cost    int    `gorm:"not null;default:0" json:"cost"`
	Visible bool   `gorm:"default:false" json:"visible"`
}

type Category struct {
	CategoryID   int    `gorm:"primaryKey" json:"category_id"`
	CategoryName string `gorm:"not null;unique" json:"category_name"`
}

type Instance struct {
	UserID 		int 	`json:"userid"`
	Level 	 	int		`json:"level"`
	Password 	string	`json:"password"`
	Port 		string	`json:"port"`
	Verified 	bool 	`json:"verified"`
	Hostname	string	`json:"hostname"`
	Deadline	int64	`json:"deadline"`
}

type Score struct {
	Username	string	`json:"username"`
	Score		string	`json:"score"`
}

type AccessDetails struct {
	Password	string	`json:"password"`
	Port		int32	`json:"port"`
	Hostname	string	`json:"hostname"`
	Deadline	int64	`json:"deadline"`
}

type ExtendDeadline struct {
	Deadline	int64	`json:"deadline"`
}