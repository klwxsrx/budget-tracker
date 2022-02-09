package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
)

type AccountService struct {
	unitOfWork UnitOfWork
}

func (s *AccountService) HandleAccountCreated(budgetID, accountID uuid.UUID, title string, initialBalance int) error {
	return s.unitOfWork.Execute(func(r RepositoryProvider) error {
		existedAccounts, err := r.AccountRepository().FindByBudgetID(budgetID)
		if err != nil {
			return err
		}
		for _, acc := range existedAccounts {
			if acc.AccountID == accountID {
				return nil
			}
		}

		acc := &model.Account{
			BudgetID:       budgetID,
			AccountID:      accountID,
			Title:          title,
			Status:         model.AccountStatusActive,
			InitialBalance: initialBalance,
			CurrentBalance: initialBalance,
			Position:       len(existedAccounts),
		}
		err = r.AccountRepository().Create(acc)
		if err != nil && !errors.Is(err, model.ErrAccountAlreadyExists) {
			return err
		}
		return nil
	})
}

func (s *AccountService) HandleAccountReordered(budgetID, accountID uuid.UUID, newPosition int) error {
	return s.unitOfWork.Execute(func(r RepositoryProvider) error {
		existedAccounts, err := r.AccountRepository().FindByBudgetID(budgetID)
		if err != nil {
			return err
		}

		specAcc := s.findAccInSlice(accountID, existedAccounts)
		if specAcc == nil {
			return fmt.Errorf("account from budget %v with id %v not found", budgetID, accountID)
		}

		lastPosition := specAcc.Position
		if lastPosition == newPosition {
			return nil
		}

		accountsToUpdate := make([]*model.Account, 0, len(existedAccounts))
		for _, acc := range existedAccounts {
			switch {
			case acc.Position == lastPosition:
				acc.Position = newPosition
			case lastPosition < newPosition && acc.Position > lastPosition && acc.Position <= newPosition:
				acc.Position--
			case lastPosition > newPosition && acc.Position < lastPosition && acc.Position >= newPosition:
				acc.Position++
			default:
				continue
			}
			accountsToUpdate = append(accountsToUpdate, acc)
		}

		for _, acc := range accountsToUpdate {
			err := r.AccountRepository().Update(acc)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *AccountService) HandleAccountRenamed(accountID uuid.UUID, title string) error {
	return s.unitOfWork.Execute(func(r RepositoryProvider) error {
		acc, err := r.AccountRepository().FindByID(accountID)
		if err != nil {
			return err
		}

		if acc.Title == title {
			return nil
		}
		acc.Title = title
		return r.AccountRepository().Update(acc)
	})
}

func (s *AccountService) HandleAccountActivated(accountID uuid.UUID) error {
	return s.unitOfWork.Execute(func(r RepositoryProvider) error {
		acc, err := r.AccountRepository().FindByID(accountID)
		if err != nil {
			return err
		}

		if acc.Status == model.AccountStatusActive {
			return nil
		}
		acc.Status = model.AccountStatusActive
		return r.AccountRepository().Update(acc)
	})
}

func (s *AccountService) HandleAccountCancelled(accountID uuid.UUID) error {
	return s.unitOfWork.Execute(func(r RepositoryProvider) error {
		acc, err := r.AccountRepository().FindByID(accountID)
		if err != nil {
			return err
		}

		if acc.Status == model.AccountStatusCancelled {
			return nil
		}
		acc.Status = model.AccountStatusCancelled
		return r.AccountRepository().Update(acc)
	})
}

func (s *AccountService) HandleAccountDeleted(budgetID, accountID uuid.UUID) error {
	return s.unitOfWork.Execute(func(r RepositoryProvider) error {
		existedAccounts, err := r.AccountRepository().FindByBudgetID(budgetID)
		if err != nil {
			return err
		}
		specAcc := s.findAccInSlice(accountID, existedAccounts)
		if specAcc == nil {
			return nil
		}

		for _, acc := range existedAccounts {
			if acc.Position > specAcc.Position {
				acc.Position--
				err := r.AccountRepository().Update(acc)
				if err != nil {
					return err
				}
			}
		}

		return r.AccountRepository().Delete(accountID)
	})
}

func (s *AccountService) findAccInSlice(id uuid.UUID, accounts []*model.Account) *model.Account {
	for _, acc := range accounts {
		if acc.AccountID == id {
			return acc
		}
	}
	return nil
}

func NewAccountService(u UnitOfWork) *AccountService {
	return &AccountService{unitOfWork: u}
}
