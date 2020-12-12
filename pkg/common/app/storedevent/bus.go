package storedevent

type Bus interface {
	Dispatch(event StoredEvent)
}

type UnsentEventProvider interface {
	GetBatch() ([]StoredEvent, error)
	SetOffset(id ID) error
}

type UnsentEventBusHandler struct {
	eventProvider UnsentEventProvider
	bus           Bus
}

func (handler *UnsentEventBusHandler) ProcessUnsentEvents() error { // TODO: add critical section interface and use it
	return nil // TODO:
}
