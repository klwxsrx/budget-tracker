package event

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"time"
)

type StoredEventID struct {
	uuid.UUID
}

type StoredEvent struct {
	Id            StoredEventID       `db:"id"`
	Type          event.Type          `db:"type"`
	AggregateID   event.AggregateID   `db:"aggregate_id"`
	AggregateName event.AggregateName `db:"aggregate_name"`
	EventData     []byte              `db:"event_data"`
	CreatedAt     time.Time           `db:"created_at"`
}

type Store interface {
	Get(id event.AggregateID) ([]*StoredEvent, error)
	GetFromID(id event.AggregateID, fromID StoredEventID) ([]*StoredEvent, error)
	GetByName(name event.AggregateName) ([]*StoredEvent, error)
	Append(e event.Event) error
}
