package service

import (
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/realtime"
)

type UnitOfWork interface {
	Execute(f func(r RepositoryProvider) error) error
}

type RepositoryProvider interface {
	BudgetRepository() model.BudgetRepository
	AccountRepository() model.AccountRepository
	RealtimeService() realtime.Service
}
