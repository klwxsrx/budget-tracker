package storedevent

import commondomainevent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"

type Serializer interface {
	Serialize(event commondomainevent.Event) (string, error)
}

type Deserializer interface {
	Deserialize(eventType string, eventData []byte) (commondomainevent.Event, error)
}
