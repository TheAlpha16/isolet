package configvars

import (
	"time"
)

type Usecase interface {
	GetString(key ConfigKey[string]) string
	GetBool(key ConfigKey[bool]) bool
	GetInt(key ConfigKey[int]) int
	GetDuration(key ConfigKey[time.Duration]) time.Duration
}

type Repository interface{}
