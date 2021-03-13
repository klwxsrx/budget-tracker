package storedevent

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/persistence"
)

type Bus interface {
	Dispatch(event *StoredEvent) error
}

type UnsentEventProvider interface {
	GetBatch() ([]*StoredEvent, error)
	SetOffset(id ID) error
}

type UnsentEventBusHandler struct {
	eventProvider UnsentEventProvider
	bus           Bus
	sync          persistence.Synchronization
}

func (handler *UnsentEventBusHandler) ProcessUnsentEvents() error {
	return handler.sync.CriticalSection("process_unsent_events", func() error {
		events, err := handler.eventProvider.GetBatch()
		for events != nil {
			for _, event := range events {
				err := handler.bus.Dispatch(event)
				if err != nil {
					return err
				}
				err = handler.eventProvider.SetOffset(event.ID)
				if err != nil {
					return err
				}
			}
			events, err = handler.eventProvider.GetBatch()
		}
		return err
	})
}
