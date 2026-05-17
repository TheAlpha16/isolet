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

// PostgreSQL enum types used by GORM models via gorm:"type:..." tags.
// GORM AutoMigrate does not emit CREATE TYPE, so we output them here
// before the table DDL so Atlas can track and diff them correctly.
const enumTypes = `
CREATE TYPE "challenge_type" AS ENUM ('static', 'dynamic', 'on-demand');
CREATE TYPE "user_role" AS ENUM ('admin', 'author', 'captain', 'player');
CREATE TYPE "protocol_type" AS ENUM ('http', 'https', 'nc', 'ssh');
CREATE TYPE "resource_name_type" AS ENUM ('cpu', 'memory', 'storage', 'ephemeral-storage');
CREATE TYPE "resource_type" AS ENUM ('request', 'limit');
`

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
	_, _ = io.WriteString(os.Stdout, enumTypes)
	_, _ = io.WriteString(os.Stdout, stmts)
}
