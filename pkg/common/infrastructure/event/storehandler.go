package event

import (
	"fmt"
	app "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	domain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

type storeEventHandler struct {
	store app.Store
}

func (eh *storeEventHandler) Handle(e domain.Event) error {
	err := eh.store.Append(e)
	if err != nil {
		return fmt.Errorf("can't append event to store, %v", err)
	}
	return nil
}

func NewStoreEventHandler(es app.Store) app.Handler {
	return &storeEventHandler{es}
}
