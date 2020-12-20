package domain

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

const (
	EventTypeAccountCreated      = "created"
	EventTypeAccountTitleChanged = "changed"
	EventTypeAccountDeleted      = "deleted"
)

type AccountCreatedEvent struct {
	ID             AccountID `json:"id"`
	Title          string    `json:"title"`
	Currency       Currency  `json:"currency"`
	InitialBalance int       `json:"initial_balance"`
}

func (e *AccountCreatedEvent) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountCreatedEvent) GetAggregateName() event.AggregateName {
	return AccountAggregateName
}

func (e *AccountCreatedEvent) GetType() event.Type {
	return EventTypeAccountCreated
}

type AccountTitleChangedEvent struct {
	ID    AccountID `json:"id"`
	Title string    `json:"title"`
}

func (e *AccountTitleChangedEvent) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountTitleChangedEvent) GetAggregateName() event.AggregateName {
	return AccountAggregateName
}

func (e *AccountTitleChangedEvent) GetType() event.Type {
	return EventTypeAccountTitleChanged
}

type AccountDeletedEvent struct {
	ID AccountID `json:"id"`
}

func (e *AccountDeletedEvent) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountDeletedEvent) GetAggregateName() event.AggregateName {
	return AccountAggregateName
}

func (e *AccountDeletedEvent) GetType() event.Type {
	return EventTypeAccountDeleted
}
