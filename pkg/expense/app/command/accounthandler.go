package command

import (
	"errors"
	"fmt"
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
	tx Transaction
}

func (h *RenameAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*RenameAccount)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.tx.Critical(updateAccountLockName, func(r DomainRegistry) error {
		return r.AccountService().Rename(domain.AccountID{UUID: cmd.ID}, cmd.Title)
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
		return r.AccountService().Delete(domain.AccountID{UUID: cmd.ID})
	})
}

func (h *DeleteAccountHandler) GetType() command.Type {
	return deleteAccountType
}

func NewCreateAccountHandler(transaction Transaction) *CreateAccountHandler {
	return &CreateAccountHandler{transaction}
}

func NewRenameAccountHandler(transaction Transaction) *RenameAccountHandler {
	return &RenameAccountHandler{transaction}
}

func NewDeleteAccountHandler(transaction Transaction) *DeleteAccountHandler {
	return &DeleteAccountHandler{transaction}
}
