package mysql

import (
	"github.com/google/uuid"
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

func (s *store) Get(id domain.AggregateID) ([]*app.StoredEvent, error) {
	return selectEvents(s.db, []string{
		"aggregate_id = UUID_TO_BIN(?)",
	}, id.String())
}

func (s *store) GetFromID(id domain.AggregateID, fromID app.StoredEventID) ([]*app.StoredEvent, error) {
	var surrogateId int
	err := s.db.Get(&surrogateId, "SELECT surrogate_id FROM event WHERE id = UUID_TO_BIN(?)", fromID.String())
	if err != nil {
		return nil, err
	}
	return selectEvents(s.db, []string{
		"aggregate_id = UUID_TO_BIN(?)",
		"surrogate_id > ?",
	}, id.String(), surrogateId)
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

func selectEvents(db mysql.Client, conditions []string, args ...interface{}) ([]*app.StoredEvent, error) {
	var events []*app.StoredEvent
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

func NewStore(client mysql.Client, serializer app.Serializer) app.Store {
	return &store{client, serializer}
}
