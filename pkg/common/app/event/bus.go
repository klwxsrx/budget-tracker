package event

type Bus interface {
	Dispatch(event StoredEvent)
}

type UnsentEventProvider interface {
	GetBatch() ([]StoredEvent, error)
	SetOffset(id StoredEventID) error
}

type BusHandler struct {
	eventProvider UnsentEventProvider
	bus           Bus
}

func (bh *BusHandler) ProcessUnsentEvents() error { // TODO: add critical section interface and use it
	return nil // TODO:
}
