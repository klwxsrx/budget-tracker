package domain

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"strings"
)

const (
	AccountAggregateName  = "expense.account"
	accountTitleMaxLength = 100
)

var (
	InvalidAccountTitleError   = errors.New("invalid title")
	AlreadyDeletedAccountError = errors.New("account is already deleted")
	UnknownAccountEventError   = errors.New("unknown event type")
)

type AccountID struct {
	uuid.UUID
}

type AccountStatus int

const (
	AccountActiveStatus AccountStatus = iota
	AccountDeletedStatus
)

type AccountState struct {
	ID             AccountID
	Status         AccountStatus
	Title          string
	Currency       Currency
	InitialBalance int
}

func (a *AccountState) Apply(e event.Event) error {
	switch e.GetType() {
	case AccountCreatedEventType:
		return a.applyCreatedEvent(e)
	case AccountTitleChangedEventType:
		return a.applyTitleChangedEvent(e)
	case AccountDeletedEventType:
		return a.applyDeletedEvent(e)
	default:
		return UnknownAccountEventError
	}
}

func (a *AccountState) applyCreatedEvent(e event.Event) error {
	ev, ok := e.(*AccountCreatedEvent)
	if !ok {
		return UnknownAccountEventError
	}

	a.ID = ev.ID
	a.Title = ev.Title
	a.Currency = ev.Currency
	a.InitialBalance = ev.InitialBalance
	return nil
}

func (a *AccountState) applyTitleChangedEvent(e event.Event) error {
	ev, ok := e.(*AccountTitleChangedEvent)
	if !ok {
		return UnknownAccountEventError
	}

	a.Title = ev.Title
	return nil
}

func (a *AccountState) applyDeletedEvent(e event.Event) error {
	_, ok := e.(*AccountDeletedEvent)
	if !ok {
		return UnknownAccountEventError
	}
	a.Status = AccountDeletedStatus
	return nil
}

type Account struct {
	state   *AccountState
	changes []event.Event
}

func (a *Account) ChangeTitle(t string) error {
	if a.state.Title == t {
		return nil
	}
	if err := validateAccountTitle(t); err != nil {
		return err
	}
	a.applyChange(&AccountTitleChangedEvent{a.state.ID, t})
	return nil
}

func (a *Account) Delete() error {
	if a.state.Status == AccountDeletedStatus {
		return AlreadyDeletedAccountError
	}
	a.applyChange(&AccountDeletedEvent{a.state.ID})
	return nil
}

func (a *Account) GetChanges() []event.Event {
	return a.changes
}

func (a *Account) applyChange(e event.Event) {
	_ = a.state.Apply(e)
	a.changes = append(a.changes, e)
}

func validateAccountTitle(title string) error {
	if len(title) == 0 || len(title) > accountTitleMaxLength {
		return InvalidAccountTitleError
	}
	return nil
}

func NewAccount(id AccountID, title string, currency Currency, initialBalance int) (*Account, error) {
	title = strings.TrimSpace(title)
	if err := validateAccountTitle(title); err != nil {
		return nil, err
	}

	state := &AccountState{id, AccountActiveStatus, title, currency, initialBalance}
	events := []event.Event{&AccountCreatedEvent{id, title, currency, initialBalance}}
	return &Account{state, events}, nil
}

func CreateAccount(s *AccountState) *Account {
	return &Account{s, make([]event.Event, 0)}
}
