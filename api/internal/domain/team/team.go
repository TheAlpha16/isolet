package team

import (
	"github.com/TheAlpha16/isolet/api/internal/domain"
)

type Team struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CaptainID int64  `json:"captain_id"`
	Password  string `json:"-"`
	domain.BaseEntity
}
