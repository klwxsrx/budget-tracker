package domain

import "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"

const (
	EventTypeAccountListCreated = "list_created"
	EventTypeAccountCreated     = "created"
	EventTypeAccountReordered   = "reordered"
	EventTypeAccountRenamed     = "renamed"
	EventTypeAccountActivated   = "activated"
	EventTypeAccountCancelled   = "cancelled"
	EventTypeAccountDeleted     = "deleted"
)

type AccountListCreatedEvent struct {
	event.Base
}

type AccountCreatedEvent struct {
	event.Base
	AccountID      AccountID
	Title          string
	Currency       Currency
	InitialBalance int
}

type AccountReorderedEvent struct {
	event.Base
	AccountID AccountID
	Position  int
}

type AccountRenamedEvent struct {
	event.Base
	AccountID AccountID
	Title     string
}

type AccountActivatedEvent struct {
	event.Base
	AccountID AccountID
}

type AccountCancelledEvent struct {
	event.Base
	AccountID AccountID
}

type AccountDeletedEvent struct {
	event.Base
	AccountID AccountID
}

func NewEventAccountListCreated(id BudgetID) event.Event {
	return &AccountListCreatedEvent{event.Base{
		EventAggregateID:   id.UUID,
		EventAggregateName: accountListAggregateName,
		EventType:          EventTypeAccountListCreated,
	}}
}

func NewEventAccountCreated(id BudgetID, accountID AccountID, title string, currency Currency, initialBalance int) event.Event {
	return &AccountCreatedEvent{
		Base: event.Base{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountCreated,
		},
		AccountID:      accountID,
		Title:          title,
		Currency:       currency,
		InitialBalance: initialBalance,
	}
}

func NewEventAccountReordered(id BudgetID, accountID AccountID, position int) event.Event {
	return &AccountReorderedEvent{
		Base: event.Base{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountReordered,
		},
		AccountID: accountID,
		Position:  position,
	}
}

func NewEventAccountRenamed(id BudgetID, accountID AccountID, title string) event.Event {
	return &AccountRenamedEvent{
		Base: event.Base{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountRenamed,
		},
		AccountID: accountID,
		Title:     title,
	}
}

func NewEventAccountActivated(id BudgetID, accountID AccountID) event.Event {
	return &AccountActivatedEvent{
		Base: event.Base{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountActivated,
		},
		AccountID: accountID,
	}
}

func NewEventAccountCancelled(id BudgetID, accountID AccountID) event.Event {
	return &AccountCancelledEvent{
		Base: event.Base{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountCancelled,
		},
		AccountID: accountID,
	}
}

func NewEventAccountDeleted(id BudgetID, accountID AccountID) event.Event {
	return &AccountDeletedEvent{
		Base: event.Base{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountDeleted,
		},
		AccountID: accountID,
	}
}
