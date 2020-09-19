package event

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"time"
)

type StoredEventID int

type StoredEvent struct {
	Id            StoredEventID       `db:"id"`
	AggregateID   event.AggregateID   `db:"aggregate_id"`
	AggregateName event.AggregateName `db:"aggregate_name"`
	Type          event.Type          `db:"event_type"`
	EventData     []byte              `db:"event_data"`
	CreatedAt     time.Time           `db:"created_at"`
}

type Store interface {
	Get(id event.AggregateID) ([]*StoredEvent, error)
	GetFromID(id event.AggregateID, fromID StoredEventID) ([]*StoredEvent, error)
	GetByName(name event.AggregateName) ([]*StoredEvent, error)
	Append(e event.Event) error
}
