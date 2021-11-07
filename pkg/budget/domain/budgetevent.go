package domain

import "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"

const (
	EventTypeBudgetCreated = budgetAggregateName + ".created"
)

type BudgetCreatedEvent struct {
	event.Base
	Title    string
	Currency string
}

func NewEventBudgetCreated(id BudgetID, title, currency string) event.Event {
	return &BudgetCreatedEvent{
		Base: event.Base{
			EventAggregateID:   id.UUID,
			EventAggregateName: budgetAggregateName,
			EventType:          EventTypeBudgetCreated,
		},
		Title:    title,
		Currency: currency,
	}
}
