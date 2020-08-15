package account

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
	"strings"
)

const (
	AggregateName  = "expense.account"
	titleMaxLength = 100
)

var (
	InvalidTitleError     = errors.New("invalid Title")
	AlreadyDeletedError   = errors.New("account is already deleted")
	UnknownEventTypeError = errors.New("unknown event type")
)

type ID struct {
	uuid.UUID
}

type Status int

const (
	activeStatus Status = iota
	deletedStatus
)

type State struct {
	ID             ID
	Status         Status
	Title          string
	Currency       domain.Currency
	InitialBalance int
}

func (a *State) Apply(e event.Event) error {
	switch e.GetType() {
	case CreatedEventType:
		return a.applyCreatedEvent(e)
	case TitleChangedEventType:
		return a.applyTitleChangedEvent(e)
	case DeletedEventType:
		return a.applyDeletedEvent(e)
	default:
		return UnknownEventTypeError
	}
}

func (a *State) applyCreatedEvent(e event.Event) error {
	ev, ok := e.(*CreatedEvent)
	if !ok {
		return UnknownEventTypeError
	}

	a.ID = ev.ID
	a.Title = ev.Title
	a.Currency = ev.Currency
	a.InitialBalance = ev.InitialBalance
	return nil
}

func (a *State) applyTitleChangedEvent(e event.Event) error {
	ev, ok := e.(*TitleChangedEvent)
	if !ok {
		return UnknownEventTypeError
	}

	a.Title = ev.Title
	return nil
}

func (a *State) applyDeletedEvent(e event.Event) error {
	_, ok := e.(*DeletedEvent)
	if !ok {
		return UnknownEventTypeError
	}
	a.Status = deletedStatus
	return nil
}

type Account struct {
	state   *State
	changes []event.Event
}

func (a *Account) ChangeTitle(t string) error {
	if a.state.Title == t {
		return nil
	}
	if err := validateTitle(t); err != nil {
		return err
	}
	a.applyChanges(&TitleChangedEvent{a.state.ID, t})
	return nil
}

func (a *Account) Delete() error {
	if a.state.Status == deletedStatus {
		return AlreadyDeletedError
	}
	a.applyChanges(&DeletedEvent{a.state.ID})
	return nil
}

func (a *Account) GetChanges() []event.Event {
	return a.changes
}

func (a *Account) applyChanges(e event.Event) {
	_ = a.state.Apply(e)
	a.changes = append(a.changes, e)
}

func validateTitle(title string) error {
	if len(title) == 0 || len(title) > titleMaxLength {
		return InvalidTitleError
	}
	return nil
}

func New(id ID, title string, currency domain.Currency, initialBalance int) (*Account, error) {
	title = strings.TrimSpace(title)
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	state := &State{id, activeStatus, title, currency, initialBalance}
	events := []event.Event{&CreatedEvent{id, title, currency, initialBalance}}
	return &Account{state, events}, nil
}

func Create(s *State) *Account {
	return &Account{s, make([]event.Event, 0)}
}
