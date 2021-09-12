package repository

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commondomainevent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type accountListRepository struct {
	dispatcher   commondomainevent.Dispatcher
	store        storedevent.Store
	deserializer storedevent.Deserializer
}

func (r *accountListRepository) Update(list *domain.AccountList) error {
	err := r.dispatcher.Dispatch(list.GetChanges())
	if err != nil {
		return fmt.Errorf("can't update aggregate: %w", err)
	}
	return nil
}

func (r *accountListRepository) FindByID(id domain.BudgetID) (*domain.AccountList, error) { // TODO: repo base implementation
	state := &domain.AccountListState{}
	storedEvents, err := r.store.GetByAggregateID(commondomainevent.AggregateID{UUID: id.UUID}, storedevent.ID{UUID: uuid.Nil})
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	if len(storedEvents) == 0 {
		return nil, nil
	}
	for _, storedEvent := range storedEvents {
		event, err := r.deserializer.Deserialize(storedEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize events: %w", err)
		}
		err = state.Apply(event)
		if err != nil {
			return nil, fmt.Errorf("failed to create accountListState: %w", err)
		}
	}
	return domain.LoadAccountList(state), nil
}

func NewAccountRepository(
	dispatcher commondomainevent.Dispatcher,
	store storedevent.Store,
	deserializer storedevent.Deserializer,
) domain.AccountListRepository {
	return &accountListRepository{dispatcher, store, deserializer}
}
