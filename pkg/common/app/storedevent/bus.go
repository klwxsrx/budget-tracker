package storedevent

import "github.com/klwxsrx/budget-tracker/pkg/common/app/persistence"

type Bus interface {
	Publish(event *StoredEvent) error
}

type UnsentEventProvider interface {
	GetBatch() ([]*StoredEvent, error)
	Ack(id ID) error
}

func NewUnsentEventHandler(
	unsentEventProvider UnsentEventProvider,
	eventBus Bus,
	sync persistence.Synchronization,
) *UnsentEventHandler {
	return &UnsentEventHandler{unsentEventProvider, eventBus, sync}
}

type UnsentEventHandler struct {
	eventProvider UnsentEventProvider
	bus           Bus
	sync          persistence.Synchronization
}

func (handler *UnsentEventHandler) ProcessUnsentEvents() error {
	return handler.sync.CriticalSection("process_unsent_events", func() error {
		events, err := handler.eventProvider.GetBatch()
		for events != nil {
			if err != nil {
				return err
			}
			for _, e := range events {
				err = handler.bus.Publish(e)
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
