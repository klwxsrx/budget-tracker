package storedevent

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/persistence"
)

type Bus interface {
	Dispatch(event *StoredEvent) error
}

type UnsentEventProvider interface {
	GetBatch() ([]*StoredEvent, error)
	Ack(id ID) error
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
			if err != nil {
				return err
			}
			for _, e := range events {
				err = handler.bus.Dispatch(e)
				if err != nil {
					return err
				}
				err = handler.eventProvider.Ack(e.ID)
				if err != nil {
					return err
				}
			}
			events, err = handler.eventProvider.GetBatch()
		}
		return err
	})
}
