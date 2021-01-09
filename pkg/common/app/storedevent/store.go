package storedevent

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"time"
)

type ID int64

type StoredEvent struct {
	ID            ID                  `db:"id"`
	AggregateID   event.AggregateID   `db:"aggregate_id"`
	AggregateName event.AggregateName `db:"aggregate_name"`
	Type          event.Type          `db:"event_type"`
	EventData     []byte              `db:"event_data"`
	CreatedAt     time.Time           `db:"created_at"`
}

type Store interface {
	LastID() (ID, error)
	GetBatch(fromID ID) ([]*StoredEvent, error)
	GetByAggregateID(id event.AggregateID, fromID ID) ([]*StoredEvent, error)
	GetByAggregateName(name event.AggregateName, fromID ID) ([]*StoredEvent, error)
	Append(e event.Event) error
}
