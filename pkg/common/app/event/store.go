package event

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"time"
)

type StoredEventID int

type StoredEvent struct {
	ID            StoredEventID       `db:"id"`
	AggregateID   event.AggregateID   `db:"aggregate_id"`
	AggregateName event.AggregateName `db:"aggregate_name"`
	Type          event.Type          `db:"event_type"`
	EventData     []byte              `db:"event_data"`
	CreatedAt     time.Time           `db:"created_at"`
}

type Store interface {
	LastID() (StoredEventID, error)
	Get(fromID StoredEventID) ([]*StoredEvent, error)
	GetByAggregateID(id event.AggregateID, fromID StoredEventID) ([]*StoredEvent, error)
	GetByAggregateName(name event.AggregateName, fromID StoredEventID) ([]*StoredEvent, error)
	Append(e event.Event) error
}
