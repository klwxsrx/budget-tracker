package messaging

import "github.com/klwxsrx/budget-tracker/pkg/common/domain"

type DomainEventSerializer interface {
	Serialize(event domain.Event) (string, error)
}

type DomainEventDeserializer interface {
	Deserialize(eventType string, eventData []byte) (domain.Event, error)
}
