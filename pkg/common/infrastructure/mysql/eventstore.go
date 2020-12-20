package mysql

import (
	"fmt"
	appEvent "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
	domain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"strings"
	"time"
)

const batchSize = 1 // TODO: change to 1000

type store struct {
	db         Client
	serializer appEvent.Serializer
}

func (s *store) LastID() (storedevent.ID, error) {
	var id storedevent.ID
	err := s.db.Get(&id, "SELECT IFNULL(MAX(id), 0) FROM event")
	return id, err
}

func (s *store) GetBatch(fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return selectEvents(s.db, fromID, batchSize, nil)
}

func (s *store) GetByAggregateID(id domain.AggregateID, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return selectEvents(s.db, fromID, 0, []string{"aggregate_id = UUID_TO_BIN(?)"}, id.String())
}

func (s *store) GetByAggregateName(name domain.AggregateName, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return selectEvents(s.db, fromID, 0, []string{"aggregate_name = ?"}, string(name))
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

func selectEvents(db Client, fromID storedevent.ID, limit int, conditions []string, args ...interface{}) ([]*storedevent.StoredEvent, error) {
	if fromID > 0 {
		conditions = append(conditions, "id > ?")
		args = append(args, fromID)
	}

	whereCondition := ""
	if conditions != nil {
		whereCondition = "WHERE " + strings.Join(conditions, " AND ") + " "
	}

	query := "SELECT " +
		"id, " +
		"BIN_TO_UUID(aggregate_id) AS aggregate_id, " +
		"aggregate_name, " +
		"event_type, " +
		"event_data, " +
		"created_at " +
		"FROM event " +
		whereCondition +
		"ORDER BY id ASC"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %v", limit)
	}

	var events []*storedevent.StoredEvent
	err := db.Select(&events, query, args...)
	return events, err
}

func NewStore(client Client, serializer appEvent.Serializer) storedevent.Store {
	return &store{client, serializer}
}
