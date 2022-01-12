package infrastructure

import (
	budgetviewapplogger "github.com/klwxsrx/budget-tracker/pkg/budgetview/app/logger"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/query"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
	commoninfrastructuremysql "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/pulsar"
)

type Container interface {
	AccountQueryService() query.AccountQueryService
	BudgetQueryService() query.BudgetQueryService
	Stop()
}

type container struct {
	accountQueryService query.AccountQueryService
	budgetQueryService  query.BudgetQueryService
}

func (c *container) AccountQueryService() query.AccountQueryService {
	return c.accountQueryService
}

func (c *container) BudgetQueryService() query.BudgetQueryService {
	return c.budgetQueryService
}

func (c *container) Stop() {
	// TODO: stop
}

func NewContainer(
	mysqlClient commoninfrastructuremysql.TransactionalClient,
	pulsarConn pulsar.Connection,
	logger log.Logger,
) (Container, error) {
	return &container{
		accountQueryService: budgetviewapplogger.NewAccountQueryServiceDecorator(
			mysql.NewAccountQueryService(mysqlClient),
			logger,
		),
		budgetQueryService: budgetviewapplogger.NewBudgetQueryServiceDecorator(
			mysql.NewBudgetQueryService(mysqlClient),
			logger,
		),
	}, nil
}
