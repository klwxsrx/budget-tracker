package persistence

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	commonDomain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	domain "github.com/klwxsrx/expense-tracker/pkg/expense/domain/account"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/account/serialization"
)

type repository struct {
	dispatcher   event.Dispatcher
	store        event.Store
	deserializer serialization.EventDeserializer
}

func (r *repository) NextID() domain.ID {
	return domain.ID{UUID: uuid.New()}
}

func (r *repository) Update(a *domain.Account) error {
	err := r.dispatcher.Dispatch(a.GetChanges())
	if err != nil {
		return fmt.Errorf("can't update aggregate: %v", err)
	}
	return nil
}

func (r *repository) GetByID(id domain.ID) (*domain.Account, error) {
	state := &domain.State{}
	storedEvents, err := r.store.Get(commonDomain.AggregateID{UUID: id.UUID})
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
	return domain.Create(state), nil
}

func (r *repository) Exists(spec domain.Specification) (bool, error) {
	storedEvents, err := r.store.GetByName(domain.AggregateName)
	if err != nil {
		return false, fmt.Errorf("failed to get events of agregates %s: %v", domain.AggregateName, err)
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

func (r *repository) buildAccountsFromEvents(events []*event.StoredEvent) ([]*domain.Account, error) {
	states := make(map[commonDomain.AggregateID]*domain.State)
	for _, storedEvent := range events {
		domainEvent, err := r.deserializer.Deserialize(storedEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize events, %v", err)
		}
		state, exists := states[domainEvent.GetAggregateID()]
		if !exists {
			state = &domain.State{}
			states[domainEvent.GetAggregateID()] = state
		}
		err = state.Apply(domainEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to apply account state, %v", err)
		}
	}
	var result []*domain.Account
	for _, state := range states {
		result = append(result, domain.Create(state))
	}
	return result, nil
}

func NewAccountRepository() domain.Repository {
	return &repository{}
}
