package infrastructure

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/event"
	budgetappmessaging "github.com/klwxsrx/budget-tracker/pkg/budget/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/mysql"
	commonappcommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	commonappstoredevent "github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commoninfrastructuremysql "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/pulsar"
)

const (
	moduleName                  = "budget"
	integrationEventHandlerName = moduleName + "_integration_event_handler"
)

type Container interface {
	CommandBus() commonappcommand.Bus
	Stop()
}

type container struct {
	bus      commonappcommand.Bus
	stopFunc func()
}

func (c *container) CommandBus() commonappcommand.Bus {
	return c.bus
}

func (c *container) Stop() {
	c.stopFunc()
}

func registerCommandHandlers(bus commonappcommand.BusRegistry, unitOfWork service.UnitOfWork) commonappcommand.Bus {
	_ = bus.Register(command.NewAccountCreateListHandler(unitOfWork))
	_ = bus.Register(command.NewAccountAddHandler(unitOfWork))
	_ = bus.Register(command.NewAccountReorderHandler(unitOfWork))
	_ = bus.Register(command.NewAccountRenameHandler(unitOfWork))
	_ = bus.Register(command.NewAccountActivateHandler(unitOfWork))
	_ = bus.Register(command.NewAccountCancelHandler(unitOfWork))
	_ = bus.Register(command.NewAccountDeleteHandler(unitOfWork))
	return bus
}

func integrationEventMessageHandler(bus commonappcommand.Bus) messaging.MessageHandler {
	deserializer := budgetappmessaging.NewDomainEventDeserializer()
	handler := messaging.NewCompositeTypedMessageHandler()
	handler.Subscribe(
		domain.EventTypeBudgetCreated,
		messaging.NewDomainEventMessageHandler(event.NewBudgetCreatedEventHandler(bus), deserializer),
	)
	return handler
}

func NewContainer(
	mysqlClient commoninfrastructuremysql.TransactionalClient,
	pulsarConn pulsar.Connection,
	logger log.Logger,
) (Container, error) {
	serializer := budgetappmessaging.NewDomainEventSerializer()
	deserializer := budgetappmessaging.NewDomainEventDeserializer()
	unitOfWork := mysql.NewUnitOfWork(mysqlClient, serializer, deserializer)

	eventBus, err := pulsar.NewEventBus(pulsarConn, moduleName)
	if err != nil {
		return nil, err
	}

	sync := commoninfrastructuremysql.NewSynchronization(mysqlClient)
	eventStore := commoninfrastructuremysql.NewEventStore(mysqlClient, serializer)
	unsentEventProvider := commoninfrastructuremysql.NewUnsentEventProvider(eventStore, mysqlClient)
	unsentEventHandler := commonappstoredevent.NewUnsentEventHandler(unsentEventProvider, eventBus, sync)
	unsentEventDispatcher := commonappstoredevent.NewUnsentEventDispatcher(unsentEventHandler, logger)
	dispatchingUnitOfWork := storedevent.NewDispatchingUnitOfWork(unitOfWork, unsentEventDispatcher)

	busRegistry := commonappcommand.NewBusRegistry(logger)
	bus := registerCommandHandlers(busRegistry, dispatchingUnitOfWork)

	integrationEventMessageConsumer, err := pulsar.NewMessageConsumer(
		pulsar.EventTopicsPattern,
		integrationEventHandlerName,
		false,
		integrationEventMessageHandler(bus),
		pulsarConn,
		logger,
	)
	if err != nil {
		return nil, err
	}

	unsentEventDispatcher.Start()
	stopFunc := func() {
		unsentEventDispatcher.Stop()
		integrationEventMessageConsumer.Stop()
	}
	return &container{bus, stopFunc}, nil
}
