package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	challengeRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/challenge"
	cvRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/configvars"
	instanceRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/instance"
	manifestRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/manifest"
	teamRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/team"
	userRepo "github.com/TheAlpha16/isolet/oracle/internal/repository/user"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&cvRepo.ConfigVars{},
		&userRepo.User{},
		&teamRepo.Team{},
		&challengeRepo.Category{},
		&challengeRepo.Challenge{},
		&challengeRepo.Hint{},
		&challengeRepo.Submission{},
		&challengeRepo.Solve{},
		&challengeRepo.UnlockedHint{},
		&instanceRepo.Instance{},
		&instanceRepo.Endpoint{},
		&manifestRepo.Manifest{},
		&manifestRepo.Resource{},
		&manifestRepo.EndpointSpec{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
