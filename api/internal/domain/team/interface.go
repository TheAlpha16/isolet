package team

import (
	"context"

	authDom "github.com/TheAlpha16/isolet/api/internal/domain/auth"
)

type Usecase interface {
	Create(ctx context.Context, input *CreateInput) (*authDom.Session, error)
	Join(ctx context.Context, input *JoinInput) (*authDom.Session, error)
	GenerateInvite(ctx context.Context) (*GenerateInviteOutput, error)
	AcceptInvite(ctx context.Context, inviteToken string) (*authDom.Session, error)
	GetByID(ctx context.Context, id int64) (*Team, error)
	GetNameByIDs(ctx context.Context, teamIDs []int64) (map[int64]string, error)
}

type Repository interface {
	Create(ctx context.Context, team *Team) (*Team, error)
	GetByName(ctx context.Context, name string) (*Team, error)
	GetByID(ctx context.Context, id int64) (*Team, error)
	GetNameByIDs(ctx context.Context, teamIDs []int64) (map[int64]string, error)
}
