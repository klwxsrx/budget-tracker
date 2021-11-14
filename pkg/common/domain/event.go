package domain

import "github.com/google/uuid"

type AggregateID struct {
	uuid.UUID
}

type Event interface {
	AggregateID() AggregateID
	AggregateName() string
	Type() string
}

type BaseEvent struct {
	EventAggregateID   uuid.UUID
	EventAggregateName string
	EventType          string
}

func (e *BaseEvent) AggregateID() AggregateID {
	return AggregateID{e.EventAggregateID}
}

func (e *BaseEvent) AggregateName() string {
	return e.EventAggregateName
}

func (e *BaseEvent) Type() string {
	return e.EventType
}
