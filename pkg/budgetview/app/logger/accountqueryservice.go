package logger

import (
	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/query"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
)

type accountQueryServiceDecorator struct {
	queryService query.AccountQueryService
	logger       log.Logger
}

func (d *accountQueryServiceDecorator) ListAccounts(budgetID uuid.UUID) ([]query.Account, error) {
	result, err := d.queryService.ListAccounts(budgetID)
	if err != nil {
		d.logger.WithError(err).Error("failed to list accounts")
	}
	return result, err
}

func NewAccountQueryServiceDecorator(
	queryService query.AccountQueryService,
	logger log.Logger,
) query.AccountQueryService {
	return &accountQueryServiceDecorator{queryService, logger}
}
