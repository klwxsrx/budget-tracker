package mysql

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type unsentEventStoreDecorator struct {
	client Client
	store  storedevent.Store
}

func (u *unsentEventStoreDecorator) GetByIDs(ids []storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return u.store.GetByIDs(ids)
}

func (u *unsentEventStoreDecorator) GetByAggregateID(id event.AggregateID, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return u.store.GetByAggregateID(id, fromID)
}

func (u *unsentEventStoreDecorator) GetByAggregateName(name string, fromID storedevent.ID) ([]*storedevent.StoredEvent, error) {
	return u.store.GetByAggregateName(name, fromID)
}

func (u *unsentEventStoreDecorator) Append(e event.Event) (storedevent.ID, error) {
	id, err := u.store.Append(e)
	if err != nil {
		return id, err
	}
	_, err = u.client.Exec("INSERT INTO unsent_event (id) VALUES (UUID_TO_BIN(?))", id.UUID)
	return id, err
}

func NewUnsentEventStoreDecorator(client Client, store storedevent.Store) storedevent.Store {
	return &unsentEventStoreDecorator{client, store}
}
