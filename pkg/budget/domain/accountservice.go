package domain

import (
	"errors"
	"github.com/google/uuid"
)

type AccountListRepository interface {
	FindByID(id BudgetID) (*AccountList, error)
	Update(list *AccountList) error
}

var (
	ErrorAccountListAlreadyExists = errors.New("account list already exists")
	ErrorAccountListDoesNotExist  = errors.New("account list does not exists")
)

type AccountListService interface {
	Create(listID BudgetID) error
	Add(listID BudgetID, title string, initialBalance MoneyAmount) (AccountID, error)
	Reorder(listID BudgetID, id AccountID, position int) error
	Rename(listID BudgetID, id AccountID, title string) error
	Activate(listID BudgetID, id AccountID) error
	Cancel(ListID BudgetID, id AccountID) error
	Delete(ListID BudgetID, id AccountID) error
}

type accountService struct {
	repo AccountListRepository
}

func (service *accountService) Create(listID BudgetID) error {
	acc, err := service.repo.FindByID(listID)
	if err != nil {
		return err
	}
	if acc != nil {
		return ErrorAccountListAlreadyExists
	}
	list := NewAccountList(listID)
	return service.repo.Update(list)
}

func (service *accountService) Add(listID BudgetID, title string, initialBalance MoneyAmount) (AccountID, error) {
	acc, err := service.repo.FindByID(listID)
	if err != nil {
		return AccountID{uuid.Nil}, err
	}
	if acc == nil {
		return AccountID{uuid.Nil}, ErrorAccountListDoesNotExist
	}
	return acc.Add(title, initialBalance)
}

func (service *accountService) Reorder(listID BudgetID, id AccountID, position int) error {
	acc, err := service.repo.FindByID(listID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrorAccountListDoesNotExist
	}
	return acc.Reorder(id, position)
}

func (service *accountService) Rename(listID BudgetID, id AccountID, title string) error {
	acc, err := service.repo.FindByID(listID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrorAccountListDoesNotExist
	}
	return acc.Rename(id, title)
}

func (service *accountService) Activate(listID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(listID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrorAccountListDoesNotExist
	}
	return acc.Activate(id)
}

func (service *accountService) Cancel(listID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(listID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrorAccountListDoesNotExist
	}
	return acc.Cancel(id)
}

func (service *accountService) Delete(listID BudgetID, id AccountID) error {
	acc, err := service.repo.FindByID(listID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrorAccountListDoesNotExist
	}
	return acc.Delete(id)
}

func NewAccountListService(repo AccountListRepository) AccountListService {
	return &accountService{repo}
}
