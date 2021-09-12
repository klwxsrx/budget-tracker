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
	EventAggregateID   uuid.UUID
	EventAggregateName string
	EventType          string
}

func (e *Base) AggregateID() AggregateID {
	return AggregateID{e.EventAggregateID}
}

func (e *Base) AggregateName() string {
	return e.EventAggregateName
}

func (e *Base) Type() string {
	return e.EventType
}
