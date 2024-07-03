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
	ChallID  int	  `json:"challid"`
	Level 	 int 	  `json:"level"`
	Name 	 string   `json:"name"`
	Prompt 	 string   `json:"prompt"`
	Category string   `json:"category"`
	Type 	 string   `json:"type"`
	Points 	 int	  `json:"points"`
	Files 	 []string `json:"files"`
	Hints 	 []string `json:"hints"`
	Solves	 int	  `json:"solves"`
	Author 	 string   `json:"author"`
	Visible  bool	  `json:"visible"`
	Tags 	 []string `json:"tags"`
	Port 	 int   	  `json:"port"`
	Subd 	 string   `json:"subd"`
	CPU 	 int 	  `json:"cpu"`
	Memory 	 int	  `json:"memory"`
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