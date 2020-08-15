package command

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

var updateAccountLockName = "update_account_lock"

type CreateAccountHandler struct {
	tx Transaction
}

func (h *CreateAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*CreateAccount)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.tx.Critical(updateAccountLockName, func(r DomainRegistry) error {
		return r.AccountService().Create(cmd.Title, domain.Currency(cmd.Currency), cmd.InitialBalance)
	})
}

func (h *CreateAccountHandler) GetType() command.Type {
	return createAccountType
}

type RenameAccountHandler struct {
	tx Transaction
}

func (h *RenameAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*RenameAccount)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.tx.Critical(updateAccountLockName, func(r DomainRegistry) error {
		id, err := uuid.Parse(cmd.ID)
		if err != nil {
			return err
		}
		return r.AccountService().Rename(domain.AccountID{UUID: id}, cmd.Title)
	})
}

func (h *RenameAccountHandler) GetType() command.Type {
	return renameAccountType
}

type DeleteAccountHandler struct {
	tx Transaction
}

func (h *DeleteAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*DeleteAccount)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.tx.Critical(updateAccountLockName, func(r DomainRegistry) error {
		id, err := uuid.Parse(cmd.ID)
		if err != nil {
			return err
		}
		return r.AccountService().Delete(domain.AccountID{UUID: id})
	})
}

func (h *DeleteAccountHandler) GetType() command.Type {
	return deleteAccountType
}
