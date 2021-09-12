package storedevent

import (
	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"

	"time"
)

type ID struct {
	uuid.UUID
}

type SurrogateID int64

type StoredEvent struct {
	ID            ID                `db:"id"`
	SurrogateID   SurrogateID       `db:"surrogate_id"`
	AggregateID   event.AggregateID `db:"aggregate_id"`
	AggregateName string            `db:"aggregate_name"`
	Type          string            `db:"event_type"`
	EventData     []byte            `db:"event_data"`
	CreatedAt     time.Time         `db:"created_at"`
}

type Store interface {
	GetByIDs(ids []ID) ([]*StoredEvent, error)
	GetByAggregateID(id event.AggregateID, fromID ID) ([]*StoredEvent, error)
	GetByAggregateName(name string, fromID ID) ([]*StoredEvent, error)
	Append(e event.Event) (ID, error)
}
