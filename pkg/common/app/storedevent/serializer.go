package storedevent

import "github.com/klwxsrx/budget-tracker/pkg/common/domain"

type Serializer interface {
	Serialize(event domain.Event) (string, error)
}

type Deserializer interface {
	Deserialize(eventType string, eventData []byte) (domain.Event, error)
}
