package team

import (
	"context"
	"time"

	"github.com/TheAlpha16/isolet/oracle/infra/database/postgres"
	"github.com/TheAlpha16/isolet/oracle/internal/domain"
	teamDom "github.com/TheAlpha16/isolet/oracle/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/oracle/internal/domain/user"
	userRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/user"
)

type Team struct {
	postgres.BaseModel
	Name      string `gorm:"uniqueIndex;not null"`
	CaptainID int64  `gorm:"column:captain_id;not null"`
	Password  string `gorm:"not null"`

	Captain userRepo.User   `gorm:"foreignKey:CaptainID;references:ID"`
	Members []userRepo.User `gorm:"foreignKey:TeamID;references:ID"`
}

func (t *Team) ToDomain(ctx context.Context) (*teamDom.Team, error) {
	members := make([]*userDom.User, len(t.Members))
	for i, m := range t.Members {
		user, err := m.ToDomain(ctx)
		if err != nil {
			return nil, err
		}
		members[i] = user
	}

	return &teamDom.Team{
		BaseEntity: domain.BaseEntity{
			CreatedAt: time.Unix(t.CreatedAt, 0),
			UpdatedAt: time.Unix(t.UpdatedAt, 0),
		},
		ID:        t.ID,
		Name:      t.Name,
		CaptainID: t.CaptainID,
		Password:  t.Password,
		Members:   members,
	}, nil
}

func NewTeamModel(t *teamDom.Team) (*Team, error) {
	team := &Team{
		Name:      t.Name,
		CaptainID: t.CaptainID,
		Password:  t.Password,
	}

	if t.ID != 0 {
		team.ID = t.ID
	}

	return team, nil
}
