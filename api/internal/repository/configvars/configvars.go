package configvars

import (
	"context"

	cvDom "github.com/TheAlpha16/isolet/api/internal/domain/configvars"

	"gorm.io/gorm"
)

type cvRepo struct {
	db *gorm.DB
}

func (cvRepo *cvRepo) Refresh(ctx context.Context) (map[string]any, error) {
	var vars []ConfigVars

	if err := cvRepo.db.WithContext(ctx).Model(ConfigVars{}).Find(&vars).Error; err != nil {
		return nil, err
	}

	result := make(map[string]any)
	for _, v := range vars {
		result[v.Key] = v.Value
	}

	return result, nil
}

func New(db *gorm.DB) cvDom.Repository {
	return &cvRepo{
		db: db,
	}
}
