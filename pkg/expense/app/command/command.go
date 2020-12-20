package command

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/command"
)

const (
	typeCreateAccount = "expense.account.create"
	typeRenameAccount = "expense.account.rename"
	typeDeleteAccount = "expense.account.delete"
)

type CreateAccount struct {
	Title          string
	Currency       string
	InitialBalance int
}

func (c *CreateAccount) GetType() command.Type {
	return typeCreateAccount
}

type RenameAccount struct {
	ID    uuid.UUID
	Title string
}

func (c *RenameAccount) GetType() command.Type {
	return typeRenameAccount
}

type DeleteAccount struct {
	ID uuid.UUID
}

func (c *DeleteAccount) GetType() command.Type {
	return typeDeleteAccount
}
