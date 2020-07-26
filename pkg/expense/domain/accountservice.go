package domain

import (
	"errors"
	"fmt"
)

var (
	AccountTitleIsDuplicated = errors.New("account with this title is already exists")
)

type AccountService interface {
	CreateAccount(id AccountID, title string, currency Currency, initialBalance int) error
}

type accountService struct {
	repo AccountRepository
}

func (a *accountService) CreateAccount(id AccountID, title string, currency Currency, initialBalance int) error {
	exists, err := a.repo.Exists(title)
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

func NewAccountService(repo AccountRepository) AccountService {
	return &accountService{repo}
}
