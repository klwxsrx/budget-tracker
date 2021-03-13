package mysql

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/serialization"
	appEvent "github.com/klwxsrx/budget-tracker/pkg/common/app/event"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	domainEvent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type accountRepository struct {
	dispatcher   appEvent.Dispatcher
	store        storedevent.Store
	deserializer serialization.EventDeserializer
}

func (r *accountRepository) NextID() domain.AccountID {
	return domain.AccountID{UUID: uuid.New()}
}

func (r *accountRepository) Update(a *domain.Account) error {
	err := r.dispatcher.Dispatch(a.GetChanges())
	if err != nil {
		return fmt.Errorf("can't update aggregate: %v", err)
	}
	return nil
}

func (r *accountRepository) GetByID(id domain.AccountID) (*domain.Account, error) {
	state := &domain.AccountState{}
	storedEvents, err := r.store.GetByAggregateID(domainEvent.AggregateID{UUID: id.UUID}, 0)
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
			return nil, fmt.Errorf("failed to create accountState: %v", err)
		}
	}
	return domain.CreateAccount(state), nil
}

func (r *accountRepository) Exists(spec domain.AccountSpecification) (bool, error) {
	storedEvents, err := r.store.GetByAggregateName(domain.AccountAggregateName, 0)
	if err != nil {
		return false, fmt.Errorf("failed to get events of agregates %v: %v", domain.AccountAggregateName, err)
	}
	accounts, err := r.buildAccountsFromEvents(storedEvents)
	if err != nil {
		return false, fmt.Errorf("failed to deserialize events of agregates: %v", err)
	}
	for _, acc := range accounts {
		if spec.IsSatisfied(acc) {
			return true, nil
		}
	}
	return false, nil
}

func (r *accountRepository) buildAccountsFromEvents(events []*storedevent.StoredEvent) ([]*domain.Account, error) {
	states := make(map[domainEvent.AggregateID]*domain.AccountState)
	for _, storedEvent := range events {
		event, err := r.deserializer.Deserialize(storedEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize events: %v", err)
		}
		state, exists := states[event.AggregateID()]
		if !exists {
			state = &domain.AccountState{}
			states[event.AggregateID()] = state
		}
		err = state.Apply(event)
		if err != nil {
			return nil, fmt.Errorf("failed to apply account state: %v", err)
		}
	}
	var result []*domain.Account
	for _, state := range states {
		result = append(result, domain.CreateAccount(state))
	}
	return result, nil
}

func NewAccountRepository(
	dispatcher appEvent.Dispatcher,
	store storedevent.Store,
	deserializer serialization.EventDeserializer,
) domain.AccountRepository {
	return &accountRepository{dispatcher, store, deserializer}
}
