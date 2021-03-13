package command

import (
	"github.com/google/uuid"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/command"
)

const (
	typeCreateAccount = "account.create"
	typeRenameAccount = "account.rename"
	typeDeleteAccount = "account.delete"
)

type CreateAccount struct {
	Title          string
	Currency       string
	InitialBalance int
}

func (c *CreateAccount) Type() command.Type {
	return typeCreateAccount
}

type RenameAccount struct {
	ID    uuid.UUID
	Title string
}

func (c *RenameAccount) Type() command.Type {
	return typeRenameAccount
}

type DeleteAccount struct {
	ID uuid.UUID
}

func (c *DeleteAccount) Type() command.Type {
	return typeDeleteAccount
}
