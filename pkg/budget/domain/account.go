package domain

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
	"strings"
)

const (
	accountListAggregateName = "account_list"
	accountTitleMaxLength    = 100
)

var (
	ErrorAccountDoesNotExist   = errors.New("account does not exist")
	ErrorAccountInvalidTitle   = errors.New("invalid title")
	ErrorAccountDuplicateTitle = errors.New("account with this title already exist")
)

type AccountID struct {
	uuid.UUID
}

type AccountStatus int

const (
	AccountStatusActive AccountStatus = iota
	AccountStatusCancelled
)

type Account interface {
	GetID() AccountID
	GetStatus() AccountStatus
	GetTitle() string
	GetInitialBalance() MoneyAmount
}

type AccountList struct {
	state   *AccountListState
	changes []event.Event
}

func (list *AccountList) GetID() BudgetID {
	return list.state.ID
}

func (list *AccountList) Add(title string, initialBalance MoneyAmount) (AccountID, error) {
	title, err := list.validateTitle(title)
	if err != nil {
		return AccountID{}, err
	}
	err = list.assertAccountWithTitleNotExist(title)
	if err != nil {
		return AccountID{}, err
	}
	accountID := AccountID{uuid.New()}
	err = list.applyChange(NewEventAccountCreated(
		list.state.ID,
		accountID,
		title,
		initialBalance.Currency,
		initialBalance.Amount,
	))
	return accountID, err
}

func (list *AccountList) Reorder(id AccountID, position int) error {
	account := list.findAccount(id)
	if account == nil {
		return ErrorAccountDoesNotExist
	}
	if position >= len(list.state.accounts) {
		position = len(list.state.accounts) - 1
	}
	if position < 0 {
		position = 0
	}
	if list.state.accounts[position].ID == id {
		return nil
	}
	return list.applyChange(NewEventAccountReordered(list.state.ID, id, position))
}

func (list *AccountList) Rename(id AccountID, title string) error {
	title, err := list.validateTitle(title)
	if err != nil {
		return err
	}
	account := list.findAccount(id)
	if account == nil {
		return ErrorAccountDoesNotExist
	}
	if account.GetTitle() == title {
		return nil
	}
	err = list.assertAccountWithTitleNotExist(title)
	if err != nil {
		return err
	}
	return list.applyChange(NewEventAccountRenamed(list.state.ID, id, title))
}

func (list *AccountList) Activate(id AccountID) error {
	account := list.findAccount(id)
	if account == nil {
		return ErrorAccountDoesNotExist
	}
	if account.GetStatus() == AccountStatusActive {
		return nil
	}
	return list.applyChange(NewEventAccountActivated(list.state.ID, id))
}

func (list *AccountList) Cancel(id AccountID) error {
	account := list.findAccount(id)
	if account == nil {
		return ErrorAccountDoesNotExist
	}
	if account.GetStatus() == AccountStatusCancelled {
		return nil
	}
	return list.applyChange(NewEventAccountCancelled(list.state.ID, id))
}

func (list *AccountList) Delete(id AccountID) error {
	account := list.findAccount(id)
	if account == nil {
		return ErrorAccountDoesNotExist
	}
	return list.applyChange(NewEventAccountDeleted(list.state.ID, id))
}

func (list *AccountList) GetChanges() []event.Event {
	return list.changes
}

func (list *AccountList) findAccount(id AccountID) Account {
	for _, acc := range list.state.accounts {
		if acc.ID == id {
			return acc
		}
	}
	return nil
}

func (list *AccountList) getAccounts() []Account {
	accounts := make([]Account, 0, len(list.state.accounts))
	for _, acc := range list.state.accounts {
		accounts = append(accounts, acc)
	}
	return accounts
}

func (list *AccountList) validateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)
	if len(title) == 0 || len(title) > accountTitleMaxLength {
		return title, ErrorAccountInvalidTitle
	}
	return title, nil
}

func (list *AccountList) assertAccountWithTitleNotExist(title string) error {
	for _, acc := range list.state.accounts {
		if acc.Title == title {
			return ErrorAccountDuplicateTitle
		}
	}
	return nil
}

func (list *AccountList) applyChange(event event.Event) error {
	err := list.state.Apply(event)
	if err != nil {
		return err
	}
	list.changes = append(list.changes, event)
	return nil
}

func NewAccountList(id BudgetID) *AccountList {
	list := LoadAccountList(&AccountListState{})
	_ = list.applyChange(NewEventAccountListCreated(id))
	return list
}

func LoadAccountList(state *AccountListState) *AccountList {
	return &AccountList{state, nil}
}
