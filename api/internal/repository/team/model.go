package team

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/api/infra/database/postgres"
	"github.com/TheAlpha16/isolet/api/internal/domain"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/api/internal/repository/user"
)

type Team struct {
	postgres.BaseModel
	Name      string `gorm:"uniqueIndex;not null"`
	CaptainID int64  `gorm:"column:captain_id;not null"`
	Password  string `gorm:"not null"`

	Captain userDom.User `gorm:"foreignKey:CaptainID;references:ID"`
}

func (t *Team) ToDomain(ctx context.Context) (*teamDom.Team, error) {
	return &teamDom.Team{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(t.CreatedAt, 0),
			UpdatedAt: time.Unix(t.UpdatedAt, 0),
		},
		ID:        t.ID,
		Name:      t.Name,
		CaptainID: t.CaptainID,
		Password:  t.Password,
	}, nil
}

func NewTeamModel(team *teamDom.Team) (*Team, error) {
	return &Team{
		Name:      team.Name,
		CaptainID: team.CaptainID,
		Password:  team.Password,
	}, nil
}
