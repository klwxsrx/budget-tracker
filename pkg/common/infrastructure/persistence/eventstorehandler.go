package persistence

import (
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

type eventStoreEventHandler struct {
	store messaging.EventStore
}

func (eh *eventStoreEventHandler) Handle(e event.Event) error {
	ev, err := messaging.NewEvent(
		messaging.EventType(e.GetType()),
		messaging.AggregateID{UUID: e.GetAggregateID().UUID},
		messaging.AggregateName(e.GetAggregateName()),
		e,
	)
	if err != nil {
		return fmt.Errorf("can't create message from domain event, %v", err)
	}
	err = eh.store.Append(ev)
	if err != nil {
		return fmt.Errorf("can't append event to store, %v", err)
	}
	return nil
}

func NewEventStoreEventHandler(es messaging.EventStore) event.Handler {
	return &eventStoreEventHandler{es}
}
