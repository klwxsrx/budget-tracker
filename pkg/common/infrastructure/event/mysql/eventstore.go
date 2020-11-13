package mysql

import (
	app "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	domain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"strings"
	"time"
)

type store struct {
	db         mysql.Client
	serializer app.Serializer
}

func (s *store) LastID() (app.StoredEventID, error) {
	var id app.StoredEventID
	err := s.db.Get(&id, "SELECT IFNULL(MAX(id), 0) FROM event")
	return id, err
}

func (s *store) Get(id domain.AggregateID) ([]*app.StoredEvent, error) {
	return selectEvents(s.db, []string{
		"aggregate_id = UUID_TO_BIN(?)",
	}, id.String())
}

func (s *store) GetFromID(id domain.AggregateID, fromID app.StoredEventID) ([]*app.StoredEvent, error) {
	return selectEvents(s.db, []string{
		"aggregate_id = UUID_TO_BIN(?)",
		"id > ?",
	}, id.String(), fromID)
}

func (s *store) GetByName(name domain.AggregateName) ([]*app.StoredEvent, error) {
	return selectEvents(s.db, []string{
		"aggregate_name = ?",
	}, string(name))
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

func selectEvents(db mysql.Client, conditions []string, args ...interface{}) ([]*app.StoredEvent, error) {
	var events []*app.StoredEvent
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

func NewStore(client mysql.Client, serializer app.Serializer) app.Store {
	return &store{client, serializer}
}
