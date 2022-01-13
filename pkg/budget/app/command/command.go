package command

import (
	"errors"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/command"
)

const (
	typeAccountCreateList = "account.create_list"
	typeAccountAdd        = "account.add"
	typeAccountReorder    = "account.reorder"
	typeAccountRename     = "account.rename"
	typeAccountActivate   = "account.activate"
	typeAccountCancel     = "account.cancel"
	typeAccountDelete     = "account.delete"
)

var errInvalidCommandType = errors.New("invalid command type")

type CreateAccountList struct {
	command.Base
	BudgetID uuid.UUID
}

type AddAccount struct {
	command.Base
	BudgetID       uuid.UUID
	Title          string
	InitialBalance int
}

type ReorderAccount struct {
	command.Base
	BudgetID  uuid.UUID
	AccountID uuid.UUID
	Position  int
}

type RenameAccount struct {
	command.Base
	BudgetID  uuid.UUID
	AccountID uuid.UUID
	Title     string
}

type ActivateAccount struct {
	command.Base
	BudgetID  uuid.UUID
	AccountID uuid.UUID
}

type CancelAccount struct {
	command.Base
	BudgetID  uuid.UUID
	AccountID uuid.UUID
}

type DeleteAccount struct {
	command.Base
	BudgetID  uuid.UUID
	AccountID uuid.UUID
}

func NewAccountCreateList(budgetID uuid.UUID) command.Command {
	return &CreateAccountList{
		Base:     command.Base{CommandType: typeAccountCreateList},
		BudgetID: budgetID,
	}
}

func NewAccountAdd(budgetID uuid.UUID, title string, initialBalance int) command.Command {
	return &AddAccount{
		Base:           command.Base{CommandType: typeAccountAdd},
		BudgetID:       budgetID,
		Title:          title,
		InitialBalance: initialBalance,
	}
}

func NewAccountReorder(budgetID, accountID uuid.UUID, position int) command.Command {
	return &ReorderAccount{
		Base:      command.Base{CommandType: typeAccountReorder},
		BudgetID:  budgetID,
		AccountID: accountID,
		Position:  position,
	}
}

func NewAccountRename(budgetID, accountID uuid.UUID, title string) command.Command {
	return &RenameAccount{
		Base:      command.Base{CommandType: typeAccountRename},
		BudgetID:  budgetID,
		AccountID: accountID,
		Title:     title,
	}
}

func NewAccountActivate(budgetID, accountID uuid.UUID) command.Command {
	return &ActivateAccount{
		Base:      command.Base{CommandType: typeAccountActivate},
		BudgetID:  budgetID,
		AccountID: accountID,
	}
}

func NewAccountCancel(budgetID, accountID uuid.UUID) command.Command {
	return &CancelAccount{
		Base:      command.Base{CommandType: typeAccountCancel},
		BudgetID:  budgetID,
		AccountID: accountID,
	}
}

func NewAccountDelete(budgetID, accountID uuid.UUID) command.Command {
	return &DeleteAccount{
		Base:      command.Base{CommandType: typeAccountDelete},
		BudgetID:  budgetID,
		AccountID: accountID,
	}
}
