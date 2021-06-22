package mysql

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type unsentEventHandlingStore struct {
	client Client
	store  storedevent.Store
}

func (u *unsentEventHandlingStore) GetByIDs(ids []storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return u.store.GetByIDs(ids)
}

func (u *unsentEventHandlingStore) GetByAggregateID(id event.AggregateID, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return u.store.GetByAggregateID(id, fromID)
}

func (u *unsentEventHandlingStore) GetByAggregateName(name string, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return u.store.GetByAggregateName(name, fromID)
}

func (u *unsentEventHandlingStore) Append(e event.Event) (storedevent.ID, error) {
	id, err := u.store.Append(e)
	if err != nil {
		return id, err
	}
	_, err = u.client.Exec("INSERT INTO unsent_event (id) VALUES (UUID_TO_BIN(?))", id.UUID)
	return id, err
}

func NewUnsentEventHandlingStore(client Client, store storedevent.Store) storedevent.Store {
	return &unsentEventHandlingStore{client, store}
}
