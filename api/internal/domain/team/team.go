package team

import (
	"github.com/TheAlpha16/isolet/api/internal/domain"
)

type Team struct {
	ID        int64
	Name      string
	CaptainID int64
	Password  string
	domain.BaseEntity
}
