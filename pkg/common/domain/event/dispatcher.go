package event

import "github.com/google/uuid"

type AggregateID uuid.UUID

type Event interface {
	GetAggregateID() AggregateID
	GetAggregateName() string
	GetType() string
}

type Handler interface {
	Handle(e Event) error
}

type Dispatcher interface {
	Dispatch(events []Event) error
	Subscribe(h Handler)
}
