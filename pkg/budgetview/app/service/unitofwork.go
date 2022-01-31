package service

import "github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"

type UnitOfWork interface {
	Execute(f func(r RepositoryProvider) error) error
}

type RepositoryProvider interface {
	BudgetRepository() model.BudgetRepository
}
