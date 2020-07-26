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
	ID             AccountID `json:"id"`
	Title          string    `json:"title"`
	Currency       Currency  `json:"currency"`
	InitialBalance int       `json:"initial_balance"`
}

func (e *AccountCreated) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountCreated) GetAggregateName() event.AggregateName {
	return accountAggregateName
}

func (e *AccountCreated) GetType() event.Type {
	return AccountCreatedEvent
}

type AccountTitleChanged struct {
	ID    AccountID `json:"id"`
	Title string    `json:"title"`
}

func (e *AccountTitleChanged) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountTitleChanged) GetAggregateName() event.AggregateName {
	return accountAggregateName
}

func (e *AccountTitleChanged) GetType() event.Type {
	return AccountTitleChangedEvent
}

type AccountDeleted struct {
	ID AccountID `json:"id"`
}

func (e *AccountDeleted) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountDeleted) GetAggregateName() event.AggregateName {
	return accountAggregateName
}

func (e *AccountDeleted) GetType() event.Type {
	return AccountDeletedEvent
}
