package domain

import "github.com/klwxsrx/budget-tracker/pkg/common/domain"

const (
	EventTypeBudgetCreated = budgetAggregateName + ".created"
)

type BudgetCreatedEvent struct {
	domain.BaseEvent
	Title    string
	Currency string
}

func NewEventBudgetCreated(id BudgetID, title, currency string) domain.Event {
	return &BudgetCreatedEvent{
		BaseEvent: domain.BaseEvent{
			EventAggregateID:   id.UUID,
			EventAggregateName: budgetAggregateName,
			EventType:          EventTypeBudgetCreated,
		},
		Title:    title,
		Currency: currency,
	}
}
