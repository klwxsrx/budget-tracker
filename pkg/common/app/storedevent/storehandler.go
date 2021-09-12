package storedevent

import (
	"fmt"

	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type storeEventHandler struct {
	store Store
}

func (eh *storeEventHandler) Handle(e event.Event) error {
	_, err := eh.store.Append(e)
	if err != nil {
		return fmt.Errorf("can't append event to store: %w", err)
	}
	return nil
}

func NewStoreEventHandler(es Store) event.Handler {
	return &storeEventHandler{es}
}
