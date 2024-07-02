package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	UserID	 int	  `json:"userid"`
	Email 	 string   `json:"email" form:"email"`
	Username string   `json:"username" form:"username"`
	Password string   `json:"password" form:"password"`
	Confirm  string   `json:"confirm" form:"confirm"`
	Rank	 int      `json:"rank"`
	TeamID	 int	  `json:"teamid"`
}

type Creds struct {
	Email 	 string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Username string `json:"username" form:"username"`
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