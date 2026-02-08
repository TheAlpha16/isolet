package postgres

import (
	"github.com/TheAlpha16/isolet/herald/pkg/facts"
)

type table string

const (
	tableInstances table = "instances"
	tableEndpoints table = "endpoints"
)

var tableHandlerMap = map[table]func(data map[string]any, factChan chan<- facts.Fact, eventType eventType){
	tableInstances: instancesHandler,
	tableEndpoints: endpointsHandler,
}

var actionMap = map[eventType]facts.Action{
	eventTypeCreate: facts.ActionCreated,
	eventTypeUpdate: facts.ActionUpdated,
	eventTypeDelete: facts.ActionDeleted,
}
