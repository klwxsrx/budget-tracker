package domain

import (
	"errors"
	"fmt"
)

var (
	AccountTitleIsDuplicated = errors.New("account with this title is already exists")
	AccountIsNotExists       = errors.New("account is not exists")
)

type AccountService interface {
	CreateAccount(id AccountID, title string, currency Currency, initialBalance int) error
}

type accountService struct {
	repo AccountRepository
}

func (a *accountService) CreateAccount(id AccountID, title string, currency Currency, initialBalance int) error {
	exists, err := a.repo.Exists(&accountWithTitleSpecification{title})
	if err != nil {
		return fmt.Errorf("exists checking failed: %v", err)
	}
	if exists {
		return AccountTitleIsDuplicated
	}

	acc, err := NewAccount(id, title, currency, initialBalance)
	if err != nil {
		return err
	}

	err = a.repo.Update(acc)
	if err != nil {
		return fmt.Errorf("account creation failed: %v", err)
	}
	return nil
}

func (a *accountService) Rename(id AccountID, title string) error {
	exists, err := a.repo.Exists(&accountWithTitleSpecification{title})
	if err != nil {
		return fmt.Errorf("exists checking failed: %v", err)
	}
	if exists {
		return AccountTitleIsDuplicated
	}

	acc, err := a.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}
	if acc == nil {
		return AccountIsNotExists
	}
	err = acc.ChangeTitle(title)
	if err != nil {
		return err
	}
	err = a.repo.Update(acc)
	if err != nil {
		return fmt.Errorf("account renaming failed: %v", err)
	}
	return nil
}

func (a *accountService) Delete(id AccountID) error {
	acc, err := a.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}
	if acc == nil {
		return AccountIsNotExists
	}
	err = acc.Delete()
	if err != nil {
		return err
	}
	err = a.repo.Update(acc)
	if err != nil {
		return fmt.Errorf("account deletion failed: %v", err)
	}
	return nil
}

func NewAccountService(repo AccountRepository) AccountService {
	return &accountService{repo}
}
