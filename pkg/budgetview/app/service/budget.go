package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
)

type BudgetService struct {
	unitOfWork UnitOfWork
}

func (s *BudgetService) HandleBudgetCreated(id uuid.UUID, title, currency string) error {
	return s.unitOfWork.Execute(func(r RepositoryProvider) error {
		budget := &model.Budget{
			ID:       id,
			Title:    title,
			Currency: currency,
		}
		err := r.BudgetRepository().Create(budget)
		if err != nil && !errors.Is(err, model.ErrBudgetAlreadyExists) {
			return err
		}
		r.RealtimeService().BudgetCreated(budget)
		return nil
	})
}

func NewBudgetService(u UnitOfWork) *BudgetService {
	return &BudgetService{unitOfWork: u}
}
