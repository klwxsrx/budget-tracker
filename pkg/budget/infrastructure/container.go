package infrastructure

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/mysql"
	commonappcommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	commonappstoredevent "github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commoninfrastructuremysql "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/pulsar"
)

const (
	moduleName = "budget"
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
	_ = bus.Register(command.NewAccountAddHandler(unitOfWork))
	_ = bus.Register(command.NewAccountReorderHandler(unitOfWork))
	_ = bus.Register(command.NewAccountRenameHandler(unitOfWork))
	_ = bus.Register(command.NewAccountActivateHandler(unitOfWork))
	_ = bus.Register(command.NewAccountCancelHandler(unitOfWork))
	_ = bus.Register(command.NewAccountDeleteHandler(unitOfWork))
	return bus
}

func NewContainer(
	client commoninfrastructuremysql.TransactionalClient,
	broker pulsar.Connection,
	loggerImpl logger.Logger,
) (Container, error) {
	serializer := storedevent.NewSerializer()
	deserializer := storedevent.NewDeserializer()
	unitOfWork := mysql.NewUnitOfWork(client, serializer, deserializer)

	storedEventSerializer := messaging.NewStoredEventSerializer()
	eventBus, err := pulsar.NewEventBus(broker, moduleName, storedEventSerializer)
	if err != nil {
		return nil, err
	}

	sync := commoninfrastructuremysql.NewSynchronization(client)
	eventStore := commoninfrastructuremysql.NewEventStore(client, serializer)
	unsentEventProvider := commoninfrastructuremysql.NewUnsentEventProvider(eventStore, client)
	unsentEventHandler := commonappstoredevent.NewUnsentEventHandler(unsentEventProvider, eventBus, sync)
	unsentEventDispatcher := commonappstoredevent.NewUnsentEventDispatcher(unsentEventHandler, loggerImpl)
	dispatchingUnitOfWork := storedevent.NewDispatchingUnitOfWork(unitOfWork, unsentEventDispatcher)

	busRegistry := commonappcommand.NewBusRegistry(command.NewResultTranslator(), loggerImpl)
	bus := registerCommandHandlers(busRegistry, dispatchingUnitOfWork)

	return &container{bus, unsentEventDispatcher.Stop}, nil
}
