package mysql

import (
	"fmt"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

const batchSize = 500

type unsentEventProvider struct {
	store  storedevent.Store
	client Client
}

func (u *unsentEventProvider) GetBatch() ([]*storedevent.StoredEvent, error) {
	ids, err := u.getUnsentIds()
	if err != nil {
		return nil, fmt.Errorf("can't get last sent id: %v", err)
	}
	if len(ids) == 0 {
		return nil, nil
	}
	return u.store.GetByIDs(ids)
}

func (u *unsentEventProvider) Ack(id storedevent.ID) error {
	_, err := u.client.Exec("DELETE FROM unsent_event WHERE id = UUID_TO_BIN(?)", id.UUID)
	return err
}

func (u *unsentEventProvider) getUnsentIds() ([]storedevent.ID, error) {
	var ids []storedevent.ID
	err := u.client.Select(&ids, "SELECT BIN_TO_UUID(id) FROM unsent_event ORDER BY id ASC LIMIT ?", batchSize)
	return ids, err
}

func NewUnsentEventProvider(store storedevent.Store, client Client) storedevent.UnsentEventProvider {
	return &unsentEventProvider{store, client}
}
