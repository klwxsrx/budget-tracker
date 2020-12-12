package storedevent

import (
	"fmt"
	appEvent "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

type storeEventHandler struct {
	store Store
}

func (eh *storeEventHandler) Handle(e event.Event) error {
	err := eh.store.Append(e)
	if err != nil {
		return fmt.Errorf("can't append event to store: %v", err)
	}
	return nil
}

func NewStoreEventHandler(es Store) appEvent.Handler {
	return &storeEventHandler{es}
}
