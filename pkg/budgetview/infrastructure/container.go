package infrastructure

import (
	budgetviewapplogger "github.com/klwxsrx/budget-tracker/pkg/budgetview/app/logger"
	budgetviewappmessaging "github.com/klwxsrx/budget-tracker/pkg/budgetview/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/query"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	commoninfrastructuremysql "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/pulsar"
)

const (
	eventHandlerName = "budget_view_event_handler"
)

type Container interface {
	AccountQueryService() query.AccountQueryService
	BudgetQueryService() query.BudgetQueryService
	Stop()
}

type container struct {
	eventMessageConsumer *pulsar.MessageConsumer
	accountQueryService  query.AccountQueryService
	budgetQueryService   query.BudgetQueryService
}

func (c *container) AccountQueryService() query.AccountQueryService {
	return c.accountQueryService
}

func (c *container) BudgetQueryService() query.BudgetQueryService {
	return c.budgetQueryService
}

func (c *container) Stop() {
	c.eventMessageConsumer.Stop()
}

func eventMessageHandler(unitOfWork service.UnitOfWork) messaging.MessageHandler {
	budgetService := service.NewBudgetService(unitOfWork)

	handler := messaging.NewCompositeTypedMessageHandler()
	handler.SubscribeTyped(budgetviewappmessaging.NewBudgetCreatedMessageHandler(budgetService))

	return handler
}

func NewContainer(
	mysqlClient commoninfrastructuremysql.TransactionalClient,
	pulsarConn pulsar.Connection,
	logger log.Logger,
) (Container, error) {
	unitOfWork := mysql.NewUnitOfWork(mysqlClient)

	eventMessageConsumer, err := pulsar.NewMessageConsumer(
		pulsar.EventTopicsPattern,
		eventHandlerName,
		true,
		eventMessageHandler(unitOfWork),
		pulsarConn,
		logger,
	)
	if err != nil {
		return nil, err
	}

	return &container{
		eventMessageConsumer: eventMessageConsumer,
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
