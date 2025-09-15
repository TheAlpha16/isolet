package profile

import (
	teamDom "github.com/TheAlpha16/isolet/api/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/api/internal/domain/user"
)

type Me struct {
	User userDom.User  `json:"user"`
	Team *teamDom.Team `json:"team,omitempty"`
}
