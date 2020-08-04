package account

import (
	"errors"
	"fmt"
)

var (
	TitleIsDuplicated = errors.New("account with this title is already exists")
	IsNotExists       = errors.New("account is not exists")
)

type Service interface {
	Create(id ID, title string, currency Currency, initialBalance int) error
}

type service struct {
	repo Repository
}

func (s *service) Create(id ID, title string, currency Currency, initialBalance int) error {
	exists, err := s.repo.Exists(&titleSpecification{title})
	if err != nil {
		return fmt.Errorf("exists checking failed: %v", err)
	}
	if exists {
		return TitleIsDuplicated
	}

	acc, err := New(id, title, currency, initialBalance)
	if err != nil {
		return err
	}

	err = s.repo.Update(acc)
	if err != nil {
		return fmt.Errorf("account creation failed: %v", err)
	}
	return nil
}

func (s *service) Rename(id ID, title string) error {
	exists, err := s.repo.Exists(&titleSpecification{title})
	if err != nil {
		return fmt.Errorf("exists checking failed: %v", err)
	}
	if exists {
		return TitleIsDuplicated
	}

	acc, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}
	if acc == nil {
		return IsNotExists
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

func (s *service) Delete(id ID) error {
	acc, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get account: %v", err)
	}
	if acc == nil {
		return IsNotExists
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

func NewService(repo Repository) Service {
	return &service{repo}
}
