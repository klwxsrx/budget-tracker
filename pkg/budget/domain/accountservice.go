package domain

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/domain"

	"github.com/google/uuid"
)

type AccountListRepository interface {
	FindByID(id BudgetID) (*AccountList, error)
	Update(list *AccountList) error
}

var (
	ErrAccountListAlreadyExists = domain.NewError("account list already exists")
	ErrAccountListDoesNotExist  = domain.NewError("account list does not exists")
)

type AccountListService interface {
	Create(budgetID BudgetID) error
	Add(budgetID BudgetID, title string, initialBalance MoneyAmount) (AccountID, error)
	Reorder(budgetID BudgetID, id AccountID, position int) error
	Rename(budgetID BudgetID, id AccountID, title string) error
	Activate(budgetID BudgetID, id AccountID) error
	Cancel(budgetID BudgetID, id AccountID) error
	Delete(budgetID BudgetID, id AccountID) error
}

type accountService struct {
	repo AccountListRepository
}

func (service *accountService) Create(budgetID BudgetID) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc != nil {
		return ErrAccountListAlreadyExists
	}
	list := NewAccountList(budgetID)
	return service.repo.Update(list)
}

func (service *accountService) Add(budgetID BudgetID, title string, initialBalance MoneyAmount) (AccountID, error) {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return AccountID{uuid.Nil}, err
	}
	if acc == nil {
		return AccountID{uuid.Nil}, ErrAccountListDoesNotExist
	}
	id, err := acc.Add(title, initialBalance)
	if err != nil {
		return AccountID{uuid.Nil}, err
	}
	err = service.repo.Update(acc)
	if err != nil {
		return AccountID{uuid.Nil}, err
	}
	return id, nil
}

func (service *accountService) Reorder(budgetID BudgetID, id AccountID, position int) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	err = acc.Reorder(id, position)
	if err != nil {
		return err
	}
	return service.repo.Update(acc)
}

func (service *accountService) Rename(budgetID BudgetID, id AccountID, title string) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	err = acc.Rename(id, title)
	if err != nil {
		return err
	}
	return service.repo.Update(acc)
}

func (service *accountService) Activate(budgetID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	err = acc.Activate(id)
	if err != nil {
		return err
	}
	return service.repo.Update(acc)
}

func (service *accountService) Cancel(budgetID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	err = acc.Cancel(id)
	if err != nil {
		return err
	}
	return service.repo.Update(acc)
}

func (service *accountService) Delete(budgetID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	err = acc.Delete(id)
	if err != nil {
		return err
	}
	return service.repo.Update(acc)
}

func NewAccountListService(repo AccountListRepository) AccountListService {
	return &accountService{repo}
}
