package mysql

import (
	"fmt"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

type unsentEventProvider struct {
	store  storedevent.Store
	client Client
}

func (u *unsentEventProvider) GetBatch() ([]*storedevent.StoredEvent, error) {
	lastSentID, err := u.getLastSentEventID()
	if err != nil {
		return nil, fmt.Errorf("can't get last sent id: %v", err)
	}
	lastID, err := u.store.LastID()
	if err != nil {
		return nil, fmt.Errorf("can't get last id from store: %v", err)
	}
	if lastID == lastSentID {
		return nil, nil
	}

	return u.store.GetBatch(lastSentID)
}

func (u *unsentEventProvider) SetOffset(id storedevent.ID) error {
	_, err := u.client.Exec("INSERT INTO last_notified_event (id, event_id)"+
		"VALUES (1, ?) ON DUPLICATE KEY UPDATE event_id = VALUES(event_id)", id)
	return err
}

func (u *unsentEventProvider) getLastSentEventID() (storedevent.ID, error) {
	var id storedevent.ID
	err := u.client.Get(&id, "SELECT IFNULL(MAX(event_id), 0) FROM last_notified_event")
	return id, err
}

func NewUnsentEventProvider(store storedevent.Store, client Client) storedevent.UnsentEventProvider {
	return &unsentEventProvider{store, client}
}
