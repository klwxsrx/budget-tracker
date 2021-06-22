package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commonStoredEvent "github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	domainEvent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type accountListRepository struct {
	dispatcher   domainEvent.Dispatcher
	store        commonStoredEvent.Store
	deserializer storedevent.Deserializer
}

func (r *accountListRepository) Update(list *domain.AccountList) error {
	err := r.dispatcher.Dispatch(list.GetChanges())
	if err != nil {
		return fmt.Errorf("can't update aggregate: %v", err)
	}
	return nil
}

func (r *accountListRepository) FindByID(id domain.BudgetID) (*domain.AccountList, error) { // TODO: repo base implementation
	state := &domain.AccountListState{}
	storedEvents, err := r.store.GetByAggregateID(domainEvent.AggregateID{UUID: id.UUID}, commonStoredEvent.ID{UUID: uuid.Nil})
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %v", err)
	}
	if len(storedEvents) == 0 {
		return nil, nil
	}
	for _, storedEvent := range storedEvents {
		event, err := r.deserializer.Deserialize(storedEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize events: %v", err)
		}
		err = state.Apply(event)
		if err != nil {
			return nil, fmt.Errorf("failed to create accountListState: %v", err)
		}
	}
	return domain.LoadAccountList(state), nil
}

func NewAccountRepository(
	dispatcher domainEvent.Dispatcher,
	store commonStoredEvent.Store,
	deserializer storedevent.Deserializer,
) domain.AccountListRepository {
	return &accountListRepository{dispatcher, store, deserializer}
}
