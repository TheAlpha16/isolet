package postgres

import (
	"strconv"

	"github.com/TheAlpha16/isolet/herald/utils/errors"
)

func parseInstanceRow(data map[string]any) (*InstanceRow, error) {
	instance := InstanceRow{}

	if idStr, ok := data["id"].(string); ok {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, errors.Raise(errors.ErrPostgresInvalidValue, "received invalid id value for instance", err)
		}
		instance.ID = id
	}

	if challengeIDStr, ok := data["challenge_id"].(string); ok {
		challengeID, err := strconv.ParseInt(challengeIDStr, 10, 64)
		if err != nil {
			return nil, errors.Raise(errors.ErrPostgresInvalidValue, "received invalid challenge_id value for instance", err)
		}
		instance.ChallengeID = challengeID
	}

	if teamIDStr, ok := data["team_id"].(string); ok {
		teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
		if err != nil {
			return nil, errors.Raise(errors.ErrPostgresInvalidValue, "received invalid team_id value for instance", err)
		}
		instance.TeamID = &teamID
	}

	if expiresAtStr, ok := data["expires_at"].(string); ok {
		expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
		if err != nil {
			return nil, errors.Raise(errors.ErrPostgresInvalidValue, "received invalid expires_at value for instance", err)
		}
		instance.ExpiresAt = expiresAt
	}

	return &instance, nil
}
