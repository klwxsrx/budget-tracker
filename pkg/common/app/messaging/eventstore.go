package messaging

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"time"
)

type StoredEventID struct {
	uuid.UUID
}

type StoredEvent struct {
	Id            StoredEventID
	Type          event.Type
	AggregateID   event.AggregateID
	AggregateName event.AggregateName
	EventData     []byte
	CreatedAt     time.Time
}

type EventStore interface {
	Get(id event.AggregateID) ([]*StoredEvent, error)
	GetFromID(id event.AggregateID, fromID StoredEventID) ([]*StoredEvent, error)
	GetByName(t event.AggregateName) ([]*StoredEvent, error)
	Append(e event.Event) error
}

func NewStoredEvent(
	id StoredEventID,
	typ event.Type,
	aggregateID event.AggregateID,
	aggregateName event.AggregateName,
	eventData []byte,
	createdAt time.Time,
) *StoredEvent {
	return &StoredEvent{id, typ, aggregateID, aggregateName, eventData, createdAt}
}
