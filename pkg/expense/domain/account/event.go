package account

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

const (
	CreatedEventType      = "expense.account_created"
	TitleChangedEventType = "expense.account_title_changed"
	DeletedEventType      = "expense.account_deleted"
)

type CreatedEvent struct {
	ID             ID              `json:"id"`
	Title          string          `json:"title"`
	Currency       domain.Currency `json:"currency"`
	InitialBalance int             `json:"initial_balance"`
}

func (e *CreatedEvent) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *CreatedEvent) GetAggregateName() event.AggregateName {
	return AggregateName
}

func (e *CreatedEvent) GetType() event.Type {
	return CreatedEventType
}

type TitleChangedEvent struct {
	ID    ID     `json:"id"`
	Title string `json:"title"`
}

func (e *TitleChangedEvent) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *TitleChangedEvent) GetAggregateName() event.AggregateName {
	return AggregateName
}

func (e *TitleChangedEvent) GetType() event.Type {
	return TitleChangedEventType
}

type DeletedEvent struct {
	ID ID `json:"id"`
}

func (e *DeletedEvent) GetAggregateID() event.AggregateID {
	return event.AggregateID{UUID: e.ID.UUID}
}

func (e *DeletedEvent) GetAggregateName() event.AggregateName {
	return AggregateName
}

func (e *DeletedEvent) GetType() event.Type {
	return DeletedEventType
}
