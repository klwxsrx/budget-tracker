package messaging

import (
	"encoding/json"
	"github.com/google/uuid"
)

type EventID struct {
	uuid.UUID
}
type EventType string
type AggregateID struct {
	uuid.UUID
}
type AggregateName string

type Event struct {
	Id            EventID
	Type          EventType
	AggregateID   AggregateID
	AggregateName AggregateName
	Payload       []byte // json
}

type EventStore interface {
	Get(id AggregateID) ([]*Event, error)
	GetFromID(id AggregateID, fromID EventID) ([]*Event, error)
	GetByName(t AggregateName) ([]*Event, error)
	Append(e *Event) error
}

func NewEvent(typ EventType, aggregateID AggregateID, name AggregateName, event interface{}) (*Event, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return &Event{EventID{uuid.New()}, typ, aggregateID, name, payload}, nil
}
