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
	return acc.Add(title, initialBalance)
}

func (service *accountService) Reorder(budgetID BudgetID, id AccountID, position int) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	return acc.Reorder(id, position)
}

func (service *accountService) Rename(budgetID BudgetID, id AccountID, title string) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	return acc.Rename(id, title)
}

func (service *accountService) Activate(budgetID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	return acc.Activate(id)
}

func (service *accountService) Cancel(budgetID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	return acc.Cancel(id)
}

func (service *accountService) Delete(budgetID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(budgetID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountListDoesNotExist
	}
	return acc.Delete(id)
}

func NewAccountListService(repo AccountListRepository) AccountListService {
	return &accountService{repo}
}
