package facts

import (
	"fmt"
	"strconv"

	"github.com/TheAlpha16/isolet/herald/utils/errors"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

const (
	NotificationFactType FactType = "Notification"
)

type Action string
type EntityName string
type NotificationSeverity string

const (
	EntityInstance EntityName = "instance"
)

const (
	ActionCreated Action = "created"
	ActionUpdated Action = "updated"
	ActionDeleted Action = "deleted"
)

const (
	SeverityInfo    NotificationSeverity = "info"
	SeverityWarning NotificationSeverity = "warning"
	SeveritySuccess NotificationSeverity = "success"
)

type Entity struct {
	Name EntityName      `json:"name"`
	ID   int64           `json:"id"`
	Data json.RawMessage `json:"data,omitempty"`
}

type Notification struct {
	BaseFact
	Entity   *Entity              `json:"entity,omitempty"`
	Action   *Action              `json:"action,omitempty"`
	Message  *string              `json:"message,omitempty"`
	Severity NotificationSeverity `json:"severity,omitempty"`
	TeamIDs  []int64              `json:"team_ids,omitempty"`
}

func (f *Notification) Key() []byte {
	var key string

	// Broadcast messages go to same partition for ordering, and we can use entity info to further partition if available
	if len(f.TeamIDs) == 0 {
		key = "broadcast"
		if f.Entity != nil {
			key += ":" + string(f.Entity.Name) + ":" + strconv.FormatInt(f.Entity.ID, 10)
		}
	}

	// Single team - partition by team for ordering
	if len(f.TeamIDs) == 1 {
		key = fmt.Sprintf("team:%d", f.TeamIDs[0])
	}

	// Multiple teams - partition by entity if available, otherwise use a random key to distribute
	if key == "" {
		if f.Entity != nil {
			key = string(f.Entity.Name) + ":" + strconv.FormatInt(f.Entity.ID, 10)
		} else {
			key = uuid.NewString()
		}
	}

	return []byte(key)
}

func (f *Notification) Validate() error {
	if err := f.BaseFact.Validate(); err != nil {
		return err
	}

	if f.Message == nil && f.Entity == nil {
		return errors.Raise(errors.ErrFactInvalidField, "Notification must have either Message or Entity", nil)
	}

	if f.Entity != nil && f.Action == nil {
		return errors.Raise(errors.ErrFactInvalidField, "Notification with Entity must have Action", nil)
	}

	if f.Action != nil && f.Entity == nil {
		return errors.Raise(errors.ErrFactInvalidField, "Notification with Action must have Entity", nil)
	}

	if f.Entity != nil {
		if f.Entity.ID <= 0 {
			return errors.Raise(errors.ErrFactInvalidField, "Notification Entity has invalid ID", nil)
		}

		switch f.Entity.Name {
		case EntityInstance:
			// valid entity
		default:
			return errors.Raise(errors.ErrFactInvalidField, "Notification has invalid Entity Name", nil)
		}
	}

	if f.Action != nil {
		switch *f.Action {
		case ActionCreated, ActionUpdated, ActionDeleted:
			// valid action
		default:
			return errors.Raise(errors.ErrFactInvalidField, "Notification has invalid Action", nil)
		}
	}

	if f.Message != nil && f.Severity == "" {
		f.Severity = SeverityInfo
	}

	if len(f.TeamIDs) > 0 {
		for _, id := range f.TeamIDs {
			if id <= 0 {
				return errors.Raise(errors.ErrFactInvalidField, "Notification has invalid TeamID", nil)
			}
		}
	}

	return nil
}

func (f *Notification) Marshal() ([]byte, error) {
	val, err := json.Marshal(f)
	if err != nil {
		return nil, errors.Raise(errors.ErrFactMarshalFailed, "failed to marshal Notification", err)
	}
	return val, nil
}
