package team

import (
	"github.com/TheAlpha16/isolet/api/internal/domain"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type Team struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	CaptainID int64           `json:"captain_id"`
	Password  string          `json:"-"`
	Members   []*userDom.User `json:"-"`
	domain.BaseEntity
}
