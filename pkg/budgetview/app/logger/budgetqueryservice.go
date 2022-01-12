package logger

import (
	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/query"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
)

type budgetQueryServiceDecorator struct {
	queryService query.BudgetQueryService
	loggerImpl   logger.Logger
}

func (d *budgetQueryServiceDecorator) ListBudgets() ([]query.Budget, error) {
	result, err := d.queryService.ListBudgets()
	if err != nil {
		d.loggerImpl.WithError(err).Error("failed to list budgets")
	}
	return result, err
}

func (d *budgetQueryServiceDecorator) ExistByIDs(ids []uuid.UUID) (bool, error) {
	result, err := d.queryService.ExistByIDs(ids)
	if err != nil {
		d.loggerImpl.WithError(err).Error("failed to check budget existence by ids")
	}
	return result, err
}

func NewBudgetQueryServiceDecorator(
	queryService query.BudgetQueryService,
	loggerImpl logger.Logger,
) query.BudgetQueryService {
	return &budgetQueryServiceDecorator{queryService, loggerImpl}
}
