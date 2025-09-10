package configvars

import (
	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
)

type ConfigVars struct {
	postgres.BaseModel
	Key   string `gorm:"uniqueIndex;not null"`
	Value string `gorm:"not null"`
}

func (cv *ConfigVars) TableName() string {
	return "config_vars"
}
