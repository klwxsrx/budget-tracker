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
	storedEvents, err := ar.store.Get(messaging.AggregateID{UUID: id.UUID})
	if err != nil {
		return nil, fmt.Errorf("failed to get events, %v", err)
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

func (ar *accountRepository) Exists(title string) (bool, error) {
	return false, nil // TODO:
}

func NewAccountRepository() domain.AccountRepository {
	return &accountRepository{}
}
