package persistence

import (
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
	accountMessaging "github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/messaging"
)

type accountRepository struct {
	dispatcher   event.Dispatcher
	store        messaging.EventStore
	deserializer accountMessaging.EventDeserializer
}

func (ar *accountRepository) Update(a *domain.Account) error {
	err := ar.dispatcher.Dispatch(a.GetChanges())
	if err != nil {
		return fmt.Errorf("can't update aggregate: %v", err)
	}
	return nil
}

func (ar *accountRepository) GetByID(id domain.AccountID) (*domain.Account, error) {
	state := &domain.AccountState{}
	storedEvents, err := ar.store.Get(event.AggregateID{UUID: id.UUID})
	if err != nil {
		return nil, fmt.Errorf("failed to get events, %v", err)
	}
	if len(storedEvents) == 0 {
		return nil, nil
	}
	for _, storedEvent := range storedEvents {
		domainEvent, err := ar.deserializer.Deserialize(storedEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize events, %v", err)
		}
		err = state.Apply(domainEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to create accountState, %v", err)
		}
	}
	return domain.CreateAccount(state), nil
}

func (ar *accountRepository) Exists(spec domain.AccountSpecification) (bool, error) {
	storedEvents, err := ar.store.GetByName(domain.AccountAggregateName)
	if err != nil {
		return false, fmt.Errorf("failed to get events of agregates %s: %v", domain.AccountAggregateName, err)
	}
	accounts, err := ar.buildAccountsFromEvents(storedEvents)
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

func (ar *accountRepository) buildAccountsFromEvents(events []*messaging.StoredEvent) ([]*domain.Account, error) {
	states := make(map[event.AggregateID]*domain.AccountState)
	for _, storedEvent := range events {
		domainEvent, err := ar.deserializer.Deserialize(storedEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize events, %v", err)
		}
		state, exists := states[domainEvent.GetAggregateID()]
		if !exists {
			state = &domain.AccountState{}
			states[domainEvent.GetAggregateID()] = state
		}
		err = state.Apply(domainEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to apply account state, %v", err)
		}
	}
	var result []*domain.Account
	for _, state := range states {
		result = append(result, domain.CreateAccount(state))
	}
	return result, nil
}

func NewAccountRepository() domain.AccountRepository {
	return &accountRepository{}
}
