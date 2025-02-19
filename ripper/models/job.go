package models

import "context"

type Job struct {
	TeamID       int64
	ChallID      int
	InstanceName string
	Ctx          context.Context
}
