package infrastructure

import (
	"context"

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

type Container interface {
	CommandBus() commonappcommand.Bus
}

type container struct {
	bus commonappcommand.Bus
}

func (c *container) CommandBus() commonappcommand.Bus {
	return c.bus
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
	ctx context.Context,
	client commoninfrastructuremysql.TransactionalClient,
	broker pulsar.Connection,
	loggerImpl logger.Logger,
) (Container, error) {
	serializer := storedevent.NewSerializer()
	deserializer := storedevent.NewDeserializer()
	unitOfWork := mysql.NewUnitOfWork(client, serializer, deserializer)

	eventStore := commoninfrastructuremysql.NewStore(client, serializer)
	unsentEventProvider := commoninfrastructuremysql.NewUnsentEventProvider(eventStore, client)
	storedEventSerializer := messaging.NewStoredEventSerializer()
	eventBus, err := pulsar.NewEventBus(ctx, broker, storedEventSerializer)
	if err != nil {
		return nil, err
	}

	sync := commoninfrastructuremysql.NewSynchronization(client)
	storedEventHandler := commonappstoredevent.NewHandler(ctx, unsentEventProvider, eventBus, sync, loggerImpl)
	storedEventHandlingUnitOfWork := storedevent.NewHandlingUnitOfWork(unitOfWork, storedEventHandler)

	busRegistry := commonappcommand.NewBusRegistry(command.NewTranslator(), loggerImpl)
	bus := registerCommandHandlers(busRegistry, storedEventHandlingUnitOfWork)

	return &container{bus}, nil
}
