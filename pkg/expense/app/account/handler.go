package account

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app"
	domainCommon "github.com/klwxsrx/expense-tracker/pkg/expense/domain"
	domain "github.com/klwxsrx/expense-tracker/pkg/expense/domain/account"
)

var updateAccountLockName = "update_account_lock"

type CreateCommandHandler struct {
	tx app.Transaction
}

func (h *CreateCommandHandler) Execute(c command.Command) error {
	cmd, ok := c.(*CreateCommand)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.tx.Critical(updateAccountLockName, func(r app.DomainRegistry) error {
		return r.AccountService().Create(cmd.Title, domainCommon.Currency(cmd.Currency), cmd.InitialBalance)
	})
}

func (h *CreateCommandHandler) GetType() command.Type {
	return createCommandType
}

type RenameCommandHandler struct {
	tx app.Transaction
}

func (h *RenameCommandHandler) Execute(c command.Command) error {
	cmd, ok := c.(*RenameCommand)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.tx.Critical(updateAccountLockName, func(r app.DomainRegistry) error {
		id, err := uuid.Parse(cmd.ID)
		if err != nil {
			return err
		}
		return r.AccountService().Rename(domain.ID{UUID: id}, cmd.Title)
	})
}

func (h *RenameCommandHandler) GetType() command.Type {
	return renameCommandType
}

type DeleteCommandHandler struct {
	tx app.Transaction
}

func (h *DeleteCommandHandler) Execute(c command.Command) error {
	cmd, ok := c.(*DeleteCommand)
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %v", c.GetType()))
	}
	return h.tx.Critical(updateAccountLockName, func(r app.DomainRegistry) error {
		id, err := uuid.Parse(cmd.ID)
		if err != nil {
			return err
		}
		return r.AccountService().Delete(domain.ID{UUID: id})
	})
}

func (h *DeleteCommandHandler) GetType() command.Type {
	return deleteCommandType
}
