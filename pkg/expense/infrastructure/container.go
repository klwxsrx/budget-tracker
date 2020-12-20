package infrastructure

import (
	"context"
	commonCommand "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	commandStoreEvent "github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
	commonMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/pulsar"
	commonSerialization "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/serialization"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/storedevent"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
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

func NewContainer(client commonMysql.TransactionalClient, broker pulsar.Connection, logger logger.Logger, ctx context.Context) (Container, error) {
	serializer := commonSerialization.NewSerializer()
	deserializer := serialization.NewEventDeserializer()
	unitOfWork := mysql.NewUnitOfWork(client, serializer, deserializer)

	eventStore := commonMysql.NewStore(client, serializer)
	unsentEventProvider := commonMysql.NewUnsentEventProvider(eventStore, client)
	eventBus, err := pulsar.NewEventBus(broker, ctx)
	if err != nil {
		return nil, err
	}
	sync := commonMysql.NewSynchronization(client)
	storedEventHandler := commandStoreEvent.NewHandler(unsentEventProvider, eventBus, sync, logger)
	storedEventHandlingUnitOfWork := storedevent.NewHandlingUnitOfWork(unitOfWork, storedEventHandler)

	bus := registerCommandHandlers(commonCommand.NewBusRegistry(logger), storedEventHandlingUnitOfWork)

	return &container{bus}, nil
}
