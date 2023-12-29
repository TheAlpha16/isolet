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
}

type Creds struct {
	Email 	 string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type VerifyClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

type Challenge struct {
	ChallID  int 	  `json:"chall_id"`
	Level 	 int 	  `json:"level"`
	Name 	 string   `json:"name"`
	Prompt 	 string   `json:"prompt"`
	Solves	 int	  `json:"solves"`
	Tags 	 []string `json:"tags"`
}

type Instance struct {
	UserID 		int 	`json:"userid"`
	Level 	 	int		`json:"level"`
	Password 	string	`json:"password"`
	Port 		string	`json:"port"`
	Verified 	bool 	`json:"verified"`
	Hostname	string	`json:"hostname"`
}

type Score struct {
	Username	string	`json:"username"`
	Score		string	`json:"score"`
}

type AccessDetails struct {
	Password	string	`json:"password"`
	Port		int32	`json:"port"`
	Hostname	string	`json:"hostname"`
}