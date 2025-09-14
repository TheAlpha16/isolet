package team

import (
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"

	"gorm.io/gorm"
)

type teamRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) teamDom.Repository {
	return &teamRepo{
		db: db,
	}
}
