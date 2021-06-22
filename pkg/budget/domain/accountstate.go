package domain

import (
	"errors"
	"fmt"
	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

var errorUnknownAccountEventType = errors.New("unknown account event")

type AccountState struct {
	ID             AccountID
	Status         AccountStatus
	Title          string
	InitialBalance MoneyAmount
}

func (state *AccountState) GetID() AccountID {
	return state.ID
}

func (state *AccountState) GetStatus() AccountStatus {
	return state.Status
}

func (state *AccountState) GetTitle() string {
	return state.Title
}

func (state *AccountState) GetInitialBalance() MoneyAmount {
	return state.InitialBalance
}

type AccountListState struct {
	ID       BudgetID
	accounts []*AccountState
}

func (state *AccountListState) Apply(e event.Event) error {
	var err error
	switch e.Type() {
	case EventTypeAccountListCreated:
		err = state.applyListCreated(e)
	case EventTypeAccountCreated:
		err = state.applyCreated(e)
	case EventTypeAccountReordered:
		err = state.applyReordered(e)
	case EventTypeAccountRenamed:
		err = state.applyRenamed(e)
	case EventTypeAccountActivated:
		err = state.applyActivated(e)
	case EventTypeAccountCancelled:
		err = state.applyCancelled(e)
	case EventTypeAccountDeleted:
		err = state.applyDeleted(e)
	default:
		err = errorUnknownAccountEventType
	}
	if err != nil {
		return fmt.Errorf("%w %v", err, e.Type())
	}
	return err
}

func (state *AccountListState) applyListCreated(event event.Event) error {
	createdEvent, ok := event.(*AccountListCreatedEvent)
	if !ok {
		return errorUnknownAccountEventType
	}
	state.ID = BudgetID{createdEvent.AggregateId}
	return nil
}

func (state *AccountListState) applyCreated(event event.Event) error {
	createdEvent, ok := event.(*AccountCreatedEvent)
	if !ok {
		return errorUnknownAccountEventType
	}
	account := &AccountState{
		ID:             createdEvent.AccountID,
		Status:         AccountStatusActive,
		Title:          createdEvent.Title,
		InitialBalance: MoneyAmount{createdEvent.InitialBalance, createdEvent.Currency},
	}
	state.accounts = append(state.accounts, account)
	return nil
}

func (state *AccountListState) applyReordered(event event.Event) error {
	reorderedEvent, ok := event.(*AccountReorderedEvent)
	if !ok {
		return errorUnknownAccountEventType
	}

	beforeIndex := state.findAccountIndex(reorderedEvent.AccountID)
	afterIndex := reorderedEvent.Position
	tmp := state.accounts[beforeIndex]

	if beforeIndex < afterIndex {
		// shift items between two indexes to the right for 1 position
		copy(state.accounts[beforeIndex:], state.accounts[beforeIndex+1:afterIndex+1])
	} else { // to the left for 1 position
		copy(state.accounts[afterIndex+1:], state.accounts[afterIndex:beforeIndex])
	}
	state.accounts[afterIndex] = tmp

	return nil
}

func (state *AccountListState) applyRenamed(event event.Event) error {
	renamedEvent, ok := event.(*AccountRenamedEvent)
	if !ok {
		return errorUnknownAccountEventType
	}
	for _, acc := range state.accounts {
		if acc.ID == renamedEvent.AccountID {
			acc.Title = renamedEvent.Title
		}
	}
	return nil
}

func (state *AccountListState) applyActivated(event event.Event) error {
	activatedEvent, ok := event.(*AccountActivatedEvent)
	if !ok {
		return errorUnknownAccountEventType
	}
	for _, acc := range state.accounts {
		if acc.ID == activatedEvent.AccountID {
			acc.Status = AccountStatusActive
		}
	}
	return nil
}

func (state *AccountListState) applyCancelled(event event.Event) error {
	cancelledEvent, ok := event.(*AccountCancelledEvent)
	if !ok {
		return errorUnknownAccountEventType
	}
	for _, acc := range state.accounts {
		if acc.ID == cancelledEvent.AccountID {
			acc.Status = AccountStatusCancelled
		}
	}
	return nil
}

func (state *AccountListState) applyDeleted(event event.Event) error {
	deletedEvent, ok := event.(*AccountDeletedEvent)
	if !ok {
		return errorUnknownAccountEventType
	}
	index := state.findAccountIndex(deletedEvent.AccountID)
	state.accounts = append(state.accounts[:index], state.accounts[index+1:]...)
	return nil
}

func (state *AccountListState) findAccountIndex(id AccountID) int {
	for index, acc := range state.accounts {
		if acc.ID == id {
			return index
		}
	}
	return 0
}
