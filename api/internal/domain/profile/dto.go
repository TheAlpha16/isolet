package profile

import (
	scoreDom "github.com/TheAlpha16/isolet/api/internal/domain/score"
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type Me struct {
	User userDom.User  `json:"user"`
	Team *teamDom.Team `json:"team,omitempty"`
}

type Team struct {
	teamDom.Team
	Records []*scoreDom.ScoreRecord `json:"records"`
}
