package command

import (
	"errors"
	"github.com/google/uuid"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/command"
)

const (
	typeAccountAdd      = "account.add"
	typeAccountReorder  = "account.reorder"
	typeAccountRename   = "account.rename"
	typeAccountActivate = "account.activate"
	typeAccountCancel   = "account.cancel"
	typeAccountDelete   = "account.delete"
)

var errorInvalidCommandType = errors.New("invalid command type")

type AddAccount struct {
	command.Base
	ListID         uuid.UUID
	Title          string
	Currency       string
	InitialBalance int
}

type ReorderAccount struct {
	command.Base
	ListID    uuid.UUID
	AccountID uuid.UUID
	Position  int
}

type RenameAccount struct {
	command.Base
	ListID    uuid.UUID
	AccountID uuid.UUID
	Title     string
}

type ActivateAccount struct {
	command.Base
	ListID    uuid.UUID
	AccountID uuid.UUID
}

type CancelAccount struct {
	command.Base
	ListID    uuid.UUID
	AccountID uuid.UUID
}

type DeleteAccount struct {
	command.Base
	ListID    uuid.UUID
	AccountID uuid.UUID
}

func NewAccountAdd(listID uuid.UUID, Title, Currency string, InitialBalance int) command.Command {
	return &AddAccount{
		Base:           command.Base{CommandType: typeAccountAdd},
		ListID:         listID,
		Title:          Title,
		Currency:       Currency,
		InitialBalance: InitialBalance,
	}
}

func NewAccountReorder(listID, accountID uuid.UUID, position int) command.Command {
	return &ReorderAccount{
		Base:      command.Base{CommandType: typeAccountReorder},
		ListID:    listID,
		AccountID: accountID,
		Position:  position,
	}
}

func NewAccountRename(listID, accountID uuid.UUID, title string) command.Command {
	return &RenameAccount{
		Base:      command.Base{CommandType: typeAccountRename},
		ListID:    listID,
		AccountID: accountID,
		Title:     title,
	}
}

func NewAccountActivate(listID, accountID uuid.UUID) command.Command {
	return &ActivateAccount{
		Base:      command.Base{CommandType: typeAccountActivate},
		ListID:    listID,
		AccountID: accountID,
	}
}

func NewAccountCancel(listID, accountID uuid.UUID) command.Command {
	return &CancelAccount{
		Base:      command.Base{CommandType: typeAccountCancel},
		ListID:    listID,
		AccountID: accountID,
	}
}

func NewAccountDelete(listID, accountID uuid.UUID) command.Command {
	return &DeleteAccount{
		Base:      command.Base{CommandType: typeAccountDelete},
		ListID:    listID,
		AccountID: accountID,
	}
}
