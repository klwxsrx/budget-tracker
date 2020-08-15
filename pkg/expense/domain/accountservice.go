package domain

import (
	"errors"
	"fmt"
)

var (
	AccountTitleIsDuplicatedError = errors.New("account with this title is already exists")
	AccountIsNotExistsError       = errors.New("account is not exists")
)

type AccountService interface {
	Create(title string, initialBalance MoneyAmount) error
	Rename(id AccountID, title string) error
	Delete(id AccountID) error
}

type accountService struct {
	repo AccountRepository
}

func (s *accountService) Create(title string, initialBalance MoneyAmount) error {
	exists, err := s.repo.Exists(&accountTitleSpecification{title})
	if err != nil {
		return fmt.Errorf("exists checking failed: %v", err)
	}
	if exists {
		return AccountTitleIsDuplicatedError
	}

	acc, err := NewAccount(s.repo.NextID(), title, initialBalance)
	if err != nil {
		return err
	}

	err = s.repo.Update(acc)
	if err != nil {
		return fmt.Errorf("account creation failed: %v", err)
	}
	return nil
}

func (s *accountService) Rename(id AccountID, title string) error {
	exists, err := s.repo.Exists(&accountTitleSpecification{title})
	if err != nil {
		return fmt.Errorf("exists checking failed: %v", err)
	}
	if exists {
		return AccountTitleIsDuplicatedError
	}

	acc, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}
	if acc == nil {
		return AccountIsNotExistsError
	}
	err = acc.ChangeTitle(title)
	if err != nil {
		return err
	}
	err = s.repo.Update(acc)
	if err != nil {
		return fmt.Errorf("account renaming failed: %v", err)
	}
	return nil
}

func (s *accountService) Delete(id AccountID) error {
	acc, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}
	if acc == nil {
		return AccountIsNotExistsError
	}
	err = acc.Delete()
	if err != nil {
		return err
	}
	err = s.repo.Update(acc)
	if err != nil {
		return fmt.Errorf("account deletion failed: %v", err)
	}
	return nil
}

func NewAccountService(repo AccountRepository) AccountService {
	return &accountService{repo}
}
