package persistence

import (
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	domain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

type eventStoreEventHandler struct {
	store event.Store
}

func (eh *eventStoreEventHandler) Handle(e domain.Event) error {
	err := eh.store.Append(e)
	if err != nil {
		return fmt.Errorf("can't append event to store, %v", err)
	}
	return nil
}

func NewEventStoreEventHandler(es event.Store) event.Handler {
	return &eventStoreEventHandler{es}
}
