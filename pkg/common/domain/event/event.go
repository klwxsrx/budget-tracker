package event

import "github.com/google/uuid"

type AggregateID struct {
	uuid.UUID
}

type Event interface {
	AggregateID() AggregateID
	AggregateName() string
	Type() string
}

type Base struct {
	AggregateId uuid.UUID
	Name        string
	EventType   string
}

func (e *Base) AggregateID() AggregateID {
	return AggregateID{e.AggregateId}
}

func (e *Base) AggregateName() string {
	return e.Name
}

func (e *Base) Type() string {
	return e.EventType
}
