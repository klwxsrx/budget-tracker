package command

import (
	"errors"
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

var updateAccountLockName = "update_account_lock"

type CreateAccountHandler struct {
	uw UnitOfWork
}

func (h *CreateAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*CreateAccount)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.uw.Critical(updateAccountLockName, func(r DomainRegistry) error {
		initialBalance, err := domain.NewMoneyAmount(cmd.InitialBalance, domain.Currency(cmd.Currency))
		if err != nil {
			return err
		}
		return r.AccountService().Create(cmd.Title, initialBalance)
	})
}

func (h *CreateAccountHandler) GetType() command.Type {
	return createAccountType
}

type RenameAccountHandler struct {
	uw UnitOfWork
}

func (h *RenameAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*RenameAccount)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.uw.Critical(updateAccountLockName, func(r DomainRegistry) error {
		return r.AccountService().Rename(domain.AccountID{UUID: cmd.ID}, cmd.Title)
	})
}

func (h *RenameAccountHandler) GetType() command.Type {
	return renameAccountType
}

type DeleteAccountHandler struct {
	uw UnitOfWork
}

func (h *DeleteAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*DeleteAccount)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.uw.Critical(updateAccountLockName, func(r DomainRegistry) error {
		return r.AccountService().Delete(domain.AccountID{UUID: cmd.ID})
	})
}

func (h *DeleteAccountHandler) GetType() command.Type {
	return deleteAccountType
}

func NewCreateAccountHandler(uw UnitOfWork) *CreateAccountHandler {
	return &CreateAccountHandler{uw}
}

func NewRenameAccountHandler(uw UnitOfWork) *RenameAccountHandler {
	return &RenameAccountHandler{uw}
}

func NewDeleteAccountHandler(uw UnitOfWork) *DeleteAccountHandler {
	return &DeleteAccountHandler{uw}
}
