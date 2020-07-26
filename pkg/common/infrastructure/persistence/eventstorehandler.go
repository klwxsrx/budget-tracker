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
	err := eh.store.Append(e)
	if err != nil {
		return fmt.Errorf("can't append event to store, %v", err)
	}
	return nil
}

func NewEventStoreEventHandler(es messaging.EventStore) event.Handler {
	return &eventStoreEventHandler{es}
}
