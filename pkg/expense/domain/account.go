package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"strings"
)

const (
	AccountAggregateName  = "expense.account"
	accountTitleMaxLength = 100
)

var (
	ErrorInvalidAccountTitle   = errors.New("invalid title")
	ErrorAlreadyDeletedAccount = errors.New("account is already deleted")
	errorUnknownAccountEvent   = errors.New("unknown account event")
)

type AccountID struct {
	uuid.UUID
}

type AccountStatus int

const (
	AccountStatusActive AccountStatus = iota
	AccountStatusDeleted
)

type AccountState struct {
	ID             AccountID
	Status         AccountStatus
	Title          string
	InitialBalance MoneyAmount
}

func (a *AccountState) Apply(e event.Event) error {
	var err error
	switch e.GetType() {
	case EventTypeAccountCreated:
		err = a.applyCreatedEvent(e)
	case EventTypeAccountTitleChanged:
		err = a.applyTitleChangedEvent(e)
	case EventTypeAccountDeleted:
		err = a.applyDeletedEvent(e)
	default:
		err = errorUnknownAccountEvent
	}
	if err != nil {
		return fmt.Errorf("%v %v", err, e.GetType())
	}
	return nil
}

func (a *AccountState) applyCreatedEvent(e event.Event) error {
	ev, ok := e.(*AccountCreatedEvent)
	if !ok {
		return errorUnknownAccountEvent
	}

	a.ID = ev.ID
	a.Title = ev.Title
	a.InitialBalance = MoneyAmount{ev.InitialBalance, ev.Currency}
	return nil
}

func (a *AccountState) applyTitleChangedEvent(e event.Event) error {
	ev, ok := e.(*AccountTitleChangedEvent)
	if !ok {
		return errorUnknownAccountEvent
	}

	a.Title = ev.Title
	return nil
}

func (a *AccountState) applyDeletedEvent(e event.Event) error {
	_, ok := e.(*AccountDeletedEvent)
	if !ok {
		return errorUnknownAccountEvent
	}
	a.Status = AccountStatusDeleted
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
	if a.state.Status == AccountStatusDeleted {
		return ErrorAlreadyDeletedAccount
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
		return ErrorInvalidAccountTitle
	}
	return nil
}

func NewAccount(id AccountID, title string, initialBalance MoneyAmount) (*Account, error) {
	title = strings.TrimSpace(title)
	if err := validateAccountTitle(title); err != nil {
		return nil, err
	}

	state := &AccountState{id, AccountStatusActive, title, initialBalance}
	events := []event.Event{&AccountCreatedEvent{id, title, initialBalance.Currency, initialBalance.Amount}}
	return &Account{state, events}, nil
}

func CreateAccount(s *AccountState) *Account {
	return &Account{s, make([]event.Event, 0)}
}
