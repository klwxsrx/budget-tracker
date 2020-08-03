package persistence

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/serialization"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"strings"
	"time"
)

type eventStore struct {
	db mysql.Client
	serialization.EventSerializer
}

func (es *eventStore) Get(id event.AggregateID) ([]*messaging.StoredEvent, error) {
	return selectEvents(es.db, []string{
		"aggregate_id = UUID_TO_BIN(?)",
	}, id.String())
}

func (es *eventStore) GetFromID(id event.AggregateID, fromID messaging.StoredEventID) ([]*messaging.StoredEvent, error) {
	var surrogateId int
	err := es.db.Get(&surrogateId, "SELECT surrogate_id FROM event WHERE id = UUID_TO_BIN(?)", fromID.String())
	if err != nil {
		return nil, err
	}
	return selectEvents(es.db, []string{
		"aggregate_id = UUID_TO_BIN(?)",
		"surrogate_id > ?",
	}, id.String(), surrogateId)
}

func (es *eventStore) GetByName(name event.AggregateName) ([]*messaging.StoredEvent, error) {
	return selectEvents(es.db, []string{
		"aggregate_name = ?",
	}, string(name))
}

func (es *eventStore) Append(e event.Event) error {
	eventData, err := es.EventSerializer.Serialize(e)
	if err != nil {
		return err
	}
	_, err = es.db.Exec(
		"INSERT INTO event"+
			"(id, type, aggregate_id, aggregate_name, event_data, created_at)"+
			"VALUES (UUID_TO_BIN(?), ?, UUID_TO_BIN(?), ?, ?, ?)",
		uuid.New(),
		e.GetType(),
		e.GetAggregateID().UUID,
		e.GetAggregateName(),
		eventData,
		time.Now(),
	)
	return err
}

func selectEvents(db mysql.Client, conditions []string, args ...interface{}) ([]*messaging.StoredEvent, error) {
	var events []*messaging.StoredEvent
	err := db.Select(&events,
		"SELECT "+
			"BIN_TO_UUID(id) AS id, "+
			"type, "+
			"BIN_TO_UUID(aggregate_id) AS aggregate_id, "+
			"aggregate_name, "+
			"event_data, "+
			"created_at "+
			"FROM event "+
			"WHERE "+strings.Join(conditions, " AND "), args...)
	return events, err
}

func NewEventStore(client mysql.Client, serializer serialization.EventSerializer) messaging.EventStore {
	return &eventStore{client, serializer}
}
