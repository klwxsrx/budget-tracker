package domain

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

const (
	AccountCreatedEvent      = "expense.account_created"
	AccountTitleChangedEvent = "expense.account_title_changed"
	AccountDeletedEvent      = "expense.account_deleted"
)

type AccountCreated struct {
	ID             AccountID
	Title          string
	Currency       Currency
	InitialBalance int
}

func (e *AccountCreated) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountCreated) GetAggregateName() string {
	return accountAggregateName
}

func (e *AccountCreated) GetType() string {
	return AccountCreatedEvent
}

type AccountTitleChanged struct {
	ID    AccountID
	Title string
}

func (e *AccountTitleChanged) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountTitleChanged) GetAggregateName() string {
	return accountAggregateName
}

func (e *AccountTitleChanged) GetType() string {
	return AccountTitleChangedEvent
}

type AccountDeleted struct {
	ID AccountID
}

func (e *AccountDeleted) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountDeleted) GetAggregateName() string {
	return accountAggregateName
}

func (e *AccountDeleted) GetType() string {
	return AccountDeletedEvent
}
