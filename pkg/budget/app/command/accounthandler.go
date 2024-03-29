package command

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/command"
)

const updateAccountLockName = "update_account_lock_"

type CreateAccountListHandler struct {
	command.Base
	unitOfWork service.UnitOfWork
}

func (h *CreateAccountListHandler) Execute(c command.Command) error {
	cmd, ok := c.(*CreateAccountList)
	if !ok {
		return errInvalidCommandType
	}
	return h.unitOfWork.Critical(updateAccountLockName+cmd.BudgetID.String(), func(r service.DomainRegistry) error {
		return r.AccountListService().Create(domain.BudgetID{UUID: cmd.BudgetID})
	})
}

type AddAccountHandler struct {
	command.Base
	unitOfWork service.UnitOfWork
}

func (h *AddAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*AddAccount)
	if !ok {
		return errInvalidCommandType
	}
	return h.unitOfWork.Critical(updateAccountLockName+cmd.BudgetID.String(), func(r service.DomainRegistry) error {
		_, err := r.AccountListService().Add(domain.BudgetID{UUID: cmd.BudgetID}, cmd.Title, domain.MoneyAmount(cmd.InitialBalance))
		return err
	})
}

type ReorderAccountHandler struct {
	command.Base
	unitOfWork service.UnitOfWork
}

func (h *ReorderAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*ReorderAccount)
	if !ok {
		return errInvalidCommandType
	}
	return h.unitOfWork.Critical(updateAccountLockName+cmd.BudgetID.String(), func(r service.DomainRegistry) error {
		return r.AccountListService().Reorder(domain.BudgetID{UUID: cmd.BudgetID}, domain.AccountID{UUID: cmd.AccountID}, cmd.Position)
	})
}

type RenameAccountHandler struct {
	command.Base
	unitOfWork service.UnitOfWork
}

func (h *RenameAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*RenameAccount)
	if !ok {
		return errInvalidCommandType
	}
	return h.unitOfWork.Critical(updateAccountLockName+cmd.BudgetID.String(), func(r service.DomainRegistry) error {
		return r.AccountListService().Rename(domain.BudgetID{UUID: cmd.BudgetID}, domain.AccountID{UUID: cmd.AccountID}, cmd.Title)
	})
}

type ActivateAccountHandler struct {
	command.Base
	unitOfWork service.UnitOfWork
}

func (h *ActivateAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*ActivateAccount)
	if !ok {
		return errInvalidCommandType
	}
	return h.unitOfWork.Critical(updateAccountLockName+cmd.BudgetID.String(), func(r service.DomainRegistry) error {
		return r.AccountListService().Activate(domain.BudgetID{UUID: cmd.BudgetID}, domain.AccountID{UUID: cmd.AccountID})
	})
}

type CancelAccountHandler struct {
	command.Base
	unitOfWork service.UnitOfWork
}

func (h *CancelAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*CancelAccount)
	if !ok {
		return errInvalidCommandType
	}
	return h.unitOfWork.Critical(updateAccountLockName+cmd.BudgetID.String(), func(r service.DomainRegistry) error {
		return r.AccountListService().Cancel(domain.BudgetID{UUID: cmd.BudgetID}, domain.AccountID{UUID: cmd.AccountID})
	})
}

type DeleteAccountHandler struct {
	command.Base
	unitOfWork service.UnitOfWork
}

func (h *DeleteAccountHandler) Execute(c command.Command) error {
	cmd, ok := c.(*DeleteAccount)
	if !ok {
		return errInvalidCommandType
	}
	return h.unitOfWork.Critical(updateAccountLockName+cmd.BudgetID.String(), func(r service.DomainRegistry) error {
		return r.AccountListService().Delete(domain.BudgetID{UUID: cmd.BudgetID}, domain.AccountID{UUID: cmd.AccountID})
	})
}

func NewAccountCreateListHandler(unitOfWork service.UnitOfWork) command.Handler {
	return &CreateAccountListHandler{
		Base:       command.Base{CommandType: typeAccountCreateList},
		unitOfWork: unitOfWork,
	}
}

func NewAccountAddHandler(unitOfWork service.UnitOfWork) command.Handler {
	return &AddAccountHandler{
		Base:       command.Base{CommandType: typeAccountAdd},
		unitOfWork: unitOfWork,
	}
}

func NewAccountReorderHandler(unitOfWork service.UnitOfWork) command.Handler {
	return &ReorderAccountHandler{
		Base:       command.Base{CommandType: typeAccountReorder},
		unitOfWork: unitOfWork,
	}
}

func NewAccountRenameHandler(unitOfWork service.UnitOfWork) command.Handler {
	return &RenameAccountHandler{
		Base:       command.Base{CommandType: typeAccountRename},
		unitOfWork: unitOfWork,
	}
}

func NewAccountActivateHandler(unitOfWork service.UnitOfWork) command.Handler {
	return &ActivateAccountHandler{
		Base:       command.Base{CommandType: typeAccountActivate},
		unitOfWork: unitOfWork,
	}
}

func NewAccountCancelHandler(unitOfWork service.UnitOfWork) command.Handler {
	return &CancelAccountHandler{
		Base:       command.Base{CommandType: typeAccountCancel},
		unitOfWork: unitOfWork,
	}
}

func NewAccountDeleteHandler(unitOfWork service.UnitOfWork) command.Handler {
	return &DeleteAccountHandler{
		Base:       command.Base{CommandType: typeAccountDelete},
		unitOfWork: unitOfWork,
	}
}
