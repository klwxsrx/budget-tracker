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

func (e *AccountCreatedEvent) AggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountCreatedEvent) AggregateName() event.AggregateName {
	return AccountAggregateName
}

func (e *AccountCreatedEvent) Type() event.Type {
	return EventTypeAccountCreated
}

type AccountTitleChangedEvent struct {
	ID    AccountID `json:"id"`
	Title string    `json:"title"`
}

func (e *AccountTitleChangedEvent) AggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountTitleChangedEvent) AggregateName() event.AggregateName {
	return AccountAggregateName
}

func (e *AccountTitleChangedEvent) Type() event.Type {
	return EventTypeAccountTitleChanged
}

type AccountDeletedEvent struct {
	ID AccountID `json:"id"`
}

func (e *AccountDeletedEvent) AggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *AccountDeletedEvent) AggregateName() event.AggregateName {
	return AccountAggregateName
}

func (e *AccountDeletedEvent) Type() event.Type {
	return EventTypeAccountDeleted
}
