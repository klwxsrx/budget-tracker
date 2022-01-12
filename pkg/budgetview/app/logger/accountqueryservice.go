package logger

import (
	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/query"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
)

type accountQueryServiceDecorator struct {
	queryService query.AccountQueryService
	loggerImpl   logger.Logger
}

func (d *accountQueryServiceDecorator) ListAccounts(budgetID uuid.UUID) ([]query.Account, error) {
	result, err := d.queryService.ListAccounts(budgetID)
	if err != nil {
		d.loggerImpl.WithError(err).Error("failed to list accounts")
	}
	return result, err
}

func NewAccountQueryServiceDecorator(
	queryService query.AccountQueryService,
	loggerImpl logger.Logger,
) query.AccountQueryService {
	return &accountQueryServiceDecorator{queryService, loggerImpl}
}
