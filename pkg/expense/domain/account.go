package domain

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"strings"
)

const (
	accountTitleMaxLength = 100
	accountAggregateName  = "expense.account"
)

var (
	InvalidAccountTitle   = errors.New("invalid Title")
	AlreadyDeletedAccount = errors.New("account is already deleted")
	UnknownEventType      = errors.New("unknown event type")
)

type AccountID struct {
	uuid.UUID
}

type AccountStatus int

const (
	activeAccount AccountStatus = iota
	deletedAccount
)

type Currency string

type AccountState struct {
	ID             AccountID
	Status         AccountStatus
	Title          string
	Currency       Currency
	InitialBalance int
}

func (a *AccountState) Apply(e event.Event) error {
	switch e.GetType() {
	case accountCreatedEvent:
		return a.applyAccountCreatedEvent(e)
	case accountTitleChangedEvent:
		return a.applyTitleChangedEvent(e)
	case accountDeletedEvent:
		return a.applyAccountDeletedEvent(e)
	default:
		return UnknownEventType
	}
}

func (a *AccountState) applyAccountCreatedEvent(e event.Event) error {
	ev, ok := e.(*AccountCreated)
	if !ok {
		return UnknownEventType
	}

	a.ID = ev.ID
	a.Title = ev.Title
	a.Currency = ev.Currency
	a.InitialBalance = ev.InitialBalance
	return nil
}

func (a *AccountState) applyTitleChangedEvent(e event.Event) error {
	ev, ok := e.(*AccountTitleChanged)
	if !ok {
		return UnknownEventType
	}

	a.Title = ev.Title
	return nil
}

func (a *AccountState) applyAccountDeletedEvent(e event.Event) error {
	_, ok := e.(*AccountDeleted)
	if !ok {
		return UnknownEventType
	}
	a.Status = deletedAccount
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
	a.applyChanges(&AccountTitleChanged{a.state.ID, t})
	return nil
}

func (a *Account) Delete() error {
	if a.state.Status == deletedAccount {
		return AlreadyDeletedAccount
	}
	a.applyChanges(&AccountDeleted{a.state.ID})
	return nil
}

func (a *Account) GetChanges() []event.Event {
	return a.changes
}

func (a *Account) applyChanges(e event.Event) {
	_ = a.state.Apply(e)
	a.changes = append(a.changes, e)
}

func validateAccountTitle(title string) error {
	if len(title) == 0 || len(title) > accountTitleMaxLength {
		return InvalidAccountTitle
	}
	return nil
}

func NewAccount(id AccountID, title string, currency Currency, initialBalance int) (*Account, error) {
	title = strings.TrimSpace(title)
	if err := validateAccountTitle(title); err != nil {
		return nil, err
	}

	state := &AccountState{id, activeAccount, title, currency, initialBalance}
	events := []event.Event{&AccountCreated{id, title, currency, initialBalance}}
	return &Account{state, events}, nil
}

func CreateAccount(s *AccountState) *Account {
	return &Account{s, make([]event.Event, 0)}
}
