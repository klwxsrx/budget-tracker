package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/common/domain"
)

type store struct {
	db         Client
	serializer storedevent.Serializer
}

func (s *store) GetByIDs(ids []storedevent.ID) ([]*storedevent.StoredEvent, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	idsStr := make([]string, 0, len(ids))
	for _, id := range ids {
		idsStr = append(idsStr, "UUID_TO_BIN('"+id.String()+"')")
	}
	return selectEvents(s.db, storedevent.ID{UUID: uuid.Nil}, []string{"id IN (" + strings.Join(idsStr, ",") + ")"})
}

func (s *store) GetByAggregate(id domain.AggregateID, name string, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return selectEvents(s.db, fromID, []string{"aggregate_id = UUID_TO_BIN(?)", "aggregate_name = ?"}, id.String(), name)
}

func (s *store) Append(e domain.Event) (storedevent.ID, error) {
	id := uuid.New()
	eventData, err := s.serializer.Serialize(e)
	if err != nil {
		return storedevent.ID{}, err
	}
	_, err = s.db.Exec(
		"INSERT INTO event "+
			"(id, aggregate_id, aggregate_name, event_type, event_data, created_at)"+
			"VALUES (UUID_TO_BIN(?), UUID_TO_BIN(?), ?, ?, ?, ?)",
		id,
		e.AggregateID().UUID,
		e.AggregateName(),
		e.Type(),
		eventData,
		time.Now(),
	)
	return storedevent.ID{UUID: id}, err
}

func selectEvents(db Client, fromID storedevent.ID, conditions []string, args ...interface{}) ([]*storedevent.StoredEvent, error) {
	if fromID.UUID != uuid.Nil {
		var id int64
		err := db.Get(&id, "SELECT surrogate_id FROM event WHERE id = UUID_TO_BIN(?)", fromID.UUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get id by event %s: %w", fromID.String(), err)
		}

		conditions = append(conditions, "surrogate_id > ?")
		args = append(args, id)
	}

	whereCondition := ""
	if conditions != nil {
		whereCondition = "WHERE " + strings.Join(conditions, " AND ") + " "
	}

	query := "SELECT " +
		"BIN_TO_UUID(id) AS id, " +
		"BIN_TO_UUID(aggregate_id) AS aggregate_id, " +
		"aggregate_name, " +
		"event_type, " +
		"event_data, " +
		"created_at " +
		"FROM event " +
		whereCondition +
		"ORDER BY surrogate_id ASC"

	var events []*storedevent.StoredEvent
	err := db.Select(&events, query, args...)
	return events, err
}

func NewEventStore(client Client, serializer storedevent.Serializer) storedevent.Store {
	return &store{client, serializer}
}
