package postgres

import (
	"github.com/TheAlpha16/isolet/herald/pkg/facts"
)

type table string

const (
	tableInstances table = "instances"
)

var tableHandlerMap = map[table]func(data map[string]any, factChan chan<- facts.Fact){
	tableInstances: instancesHandler,
}

func instancesHandler(data map[string]any, factChan chan<- facts.Fact) {}
