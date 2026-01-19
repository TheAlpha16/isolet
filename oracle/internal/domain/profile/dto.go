package profile

import (
	challengeDom "github.com/TheAlpha16/isolet/oracle/internal/domain/challenge"
	teamDom "github.com/TheAlpha16/isolet/oracle/internal/domain/team"
	userDom "github.com/TheAlpha16/isolet/oracle/internal/domain/user"
)

type Me struct {
	User userDom.User  `json:"user"`
	Team *teamDom.Team `json:"team,omitempty"`
}

type Team struct {
	teamDom.Team
	Score       int                           `json:"score"`
	Rank        *int                          `json:"rank,omitempty"`
	Members     []*userDom.User               `json:"members"`
	Submissions []*challengeDom.SubmissionDTO `json:"submissions"`
}
