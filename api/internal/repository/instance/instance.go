package instance

import (
	instanceDom "github.com/TheAlpha16/isolet/api/internal/domain/instance"
	"gorm.io/gorm"
)

type instanceRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) instanceDom.Repository {
	return &instanceRepo{
		db: db,
	}
}
