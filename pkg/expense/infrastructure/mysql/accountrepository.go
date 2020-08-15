package mysql

import (
	"fmt"
	"github.com/google/uuid"
	eventApp "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	eventDomain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
)

type accountRepository struct {
	dispatcher   eventApp.Dispatcher
	store        eventApp.Store
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
	storedEvents, err := r.store.Get(eventDomain.AggregateID{UUID: id.UUID})
	if err != nil {
		return nil, fmt.Errorf("failed to get events, %v", err)
	}
	if len(storedEvents) == 0 {
		return nil, nil
	}
	for _, storedEvent := range storedEvents {
		domainEvent, err := r.deserializer.Deserialize(storedEvent)
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

func (r *accountRepository) Exists(spec domain.AccountSpecification) (bool, error) {
	storedEvents, err := r.store.GetByName(domain.AccountAggregateName)
	if err != nil {
		return false, fmt.Errorf("failed to get events of agregates %s: %v", domain.AccountAggregateName, err)
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

func (r *accountRepository) buildAccountsFromEvents(events []*eventApp.StoredEvent) ([]*domain.Account, error) {
	states := make(map[eventDomain.AggregateID]*domain.AccountState)
	for _, storedEvent := range events {
		domainEvent, err := r.deserializer.Deserialize(storedEvent)
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
