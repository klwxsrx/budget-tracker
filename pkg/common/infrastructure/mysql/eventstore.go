package mysql

import (
	appEvent "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
	domain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"strings"
	"time"
)

type store struct {
	db         Client
	serializer appEvent.Serializer
}

func (s *store) LastID() (storedevent.ID, error) {
	var id storedevent.ID
	err := s.db.Get(&id, "SELECT IFNULL(MAX(id), 0) FROM event")
	return id, err
}

func (s *store) Get(fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return selectEvents(s.db, []string{
		"id > ?",
	}, fromID)
}

func (s *store) GetByAggregateID(id domain.AggregateID, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return selectEvents(s.db, []string{
		"aggregate_id = UUID_TO_BIN(?)",
		"id > ?",
	}, id.String(), fromID)
}

func (s *store) GetByAggregateName(name domain.AggregateName, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return selectEvents(s.db, []string{
		"aggregate_name = ?",
		"id > ?",
	}, string(name), fromID)
}

func (s *store) Append(e domain.Event) error {
	eventData, err := s.serializer.Serialize(e)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(
		"INSERT INTO event"+
			"(aggregate_id, aggregate_name, event_type, event_data, created_at)"+
			"VALUES (UUID_TO_BIN(?), ?, ?, ?, ?)",
		e.GetAggregateID().UUID,
		e.GetAggregateName(),
		e.GetType(),
		eventData,
		time.Now(),
	)
	return err
}

func selectEvents(db Client, conditions []string, args ...interface{}) ([]*storedevent.StoredEvent, error) {
	var events []*storedevent.StoredEvent
	err := db.Select(&events,
		"SELECT "+
			"id, "+
			"BIN_TO_UUID(aggregate_id) AS aggregate_id, "+
			"aggregate_name, "+
			"event_type, "+
			"event_data, "+
			"created_at "+
			"FROM event "+
			"WHERE "+strings.Join(conditions, " AND "), args...)
	return events, err
}

func NewStore(client Client, serializer appEvent.Serializer) storedevent.Store {
	return &store{client, serializer}
}
