package infrastructure

import (
	"context"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/mysql"
	commonCommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	commonStoredEvent "github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commonMysql "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/pulsar"
)

type Container interface {
	CommandBus() commonCommand.Bus
}

type container struct {
	bus commonCommand.Bus
}

func (c *container) CommandBus() commonCommand.Bus {
	return c.bus
}

func registerCommandHandlers(bus commonCommand.BusRegistry, unitOfWork command.UnitOfWork) commonCommand.Bus {
	_ = bus.Register(command.NewCreateAccountHandler(unitOfWork))
	_ = bus.Register(command.NewRenameAccountHandler(unitOfWork))
	_ = bus.Register(command.NewDeleteAccountHandler(unitOfWork))
	return bus
}

func NewContainer(
	client commonMysql.TransactionalClient,
	broker pulsar.Connection,
	logger logger.Logger,
	ctx context.Context,
) (Container, error) {
	serializer := storedevent.NewSerializer()
	deserializer := storedevent.NewDeserializer()
	unitOfWork := mysql.NewUnitOfWork(client, serializer, deserializer)

	eventStore := commonMysql.NewStore(client, serializer)
	unsentEventProvider := commonMysql.NewUnsentEventProvider(eventStore, client)
	storedEventSerializer := messaging.NewStoredEventSerializer()
	eventBus, err := pulsar.NewEventBus(broker, storedEventSerializer, ctx)
	if err != nil {
		return nil, err
	}

	sync := commonMysql.NewSynchronization(client)
	storedEventHandler := commonStoredEvent.NewHandler(unsentEventProvider, eventBus, sync, logger, ctx)
	storedEventHandlingUnitOfWork := storedevent.NewHandlingUnitOfWork(unitOfWork, storedEventHandler)

	busRegistry := commonCommand.NewBusRegistry(command.ResultMap, logger)
	bus := registerCommandHandlers(busRegistry, storedEventHandlingUnitOfWork)

	return &container{bus}, nil
}
