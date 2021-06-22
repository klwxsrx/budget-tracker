package storedevent

import commonEvent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"

type Serializer interface {
	Serialize(event commonEvent.Event) (string, error)
}

type Deserializer interface {
	Deserialize(event *StoredEvent) (commonEvent.Event, error)
}
