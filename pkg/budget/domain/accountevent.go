package domain

import "github.com/klwxsrx/budget-tracker/pkg/common/domain"

const (
	EventTypeAccountListCreated = accountListAggregateName + ".created"
	EventTypeAccountCreated     = accountListAggregateName + ".account_created"
	EventTypeAccountReordered   = accountListAggregateName + ".account_reordered"
	EventTypeAccountRenamed     = accountListAggregateName + ".account_renamed"
	EventTypeAccountActivated   = accountListAggregateName + ".account_activated"
	EventTypeAccountCancelled   = accountListAggregateName + ".account_cancelled"
	EventTypeAccountDeleted     = accountListAggregateName + ".account_deleted"
)

type AccountListCreatedEvent struct {
	domain.BaseEvent
}

type AccountCreatedEvent struct {
	domain.BaseEvent
	AccountID      AccountID
	Title          string
	InitialBalance int
}

type AccountReorderedEvent struct {
	domain.BaseEvent
	AccountID AccountID
	Position  int
}

type AccountRenamedEvent struct {
	domain.BaseEvent
	AccountID AccountID
	Title     string
}

type AccountActivatedEvent struct {
	domain.BaseEvent
	AccountID AccountID
}

type AccountCancelledEvent struct {
	domain.BaseEvent
	AccountID AccountID
}

type AccountDeletedEvent struct {
	domain.BaseEvent
	AccountID AccountID
}

func NewEventAccountListCreated(id BudgetID) domain.Event {
	return &AccountListCreatedEvent{domain.BaseEvent{
		EventAggregateID:   id.UUID,
		EventAggregateName: accountListAggregateName,
		EventType:          EventTypeAccountListCreated,
	}}
}

func NewEventAccountCreated(id BudgetID, accountID AccountID, title string, initialBalance MoneyAmount) domain.Event {
	return &AccountCreatedEvent{
		BaseEvent: domain.BaseEvent{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountCreated,
		},
		AccountID:      accountID,
		Title:          title,
		InitialBalance: int(initialBalance),
	}
}

func NewEventAccountReordered(id BudgetID, accountID AccountID, position int) domain.Event {
	return &AccountReorderedEvent{
		BaseEvent: domain.BaseEvent{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountReordered,
		},
		AccountID: accountID,
		Position:  position,
	}
}

func NewEventAccountRenamed(id BudgetID, accountID AccountID, title string) domain.Event {
	return &AccountRenamedEvent{
		BaseEvent: domain.BaseEvent{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountRenamed,
		},
		AccountID: accountID,
		Title:     title,
	}
}

func NewEventAccountActivated(id BudgetID, accountID AccountID) domain.Event {
	return &AccountActivatedEvent{
		BaseEvent: domain.BaseEvent{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountActivated,
		},
		AccountID: accountID,
	}
}

func NewEventAccountCancelled(id BudgetID, accountID AccountID) domain.Event {
	return &AccountCancelledEvent{
		BaseEvent: domain.BaseEvent{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountCancelled,
		},
		AccountID: accountID,
	}
}

func NewEventAccountDeleted(id BudgetID, accountID AccountID) domain.Event {
	return &AccountDeletedEvent{
		BaseEvent: domain.BaseEvent{
			EventAggregateID:   id.UUID,
			EventAggregateName: accountListAggregateName,
			EventType:          EventTypeAccountDeleted,
		},
		AccountID: accountID,
	}
}
