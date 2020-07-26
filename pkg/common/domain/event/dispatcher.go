package event

import "github.com/google/uuid"

type AggregateID struct {
	uuid.UUID
}
type AggregateName string
type Type string

type Event interface {
	GetAggregateID() AggregateID
	GetAggregateName() AggregateName
	GetType() Type
}

type Handler interface {
	Handle(e Event) error
}

type Dispatcher interface {
	Dispatch(events []Event) error
	Subscribe(h Handler)
}
