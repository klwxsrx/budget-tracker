package infrastructure

import (
	commonCommand "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
	commonMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/pulsar"
	commonSerialization "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/serialization"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/event"
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

func NewContainer(client commonMysql.TransactionalClient, broker pulsar.Connection, logger logger.Logger) (Container, error) {
	serializer := commonSerialization.NewSerializer()
	deserializer := serialization.NewEventDeserializer()
	unitOfWork := mysql.NewUnitOfWork(client, serializer, deserializer)

	_ = commonMysql.NewStore(client, serializer)            // TODO:
	var unsentEventProvider storedevent.UnsentEventProvider // TODO:
	var eventBus storedevent.Bus                            // TODO:
	storedEventHandler := storedevent.NewHandler(unsentEventProvider, eventBus, logger)
	notifyingUnitOfWork := event.NewStoredEventHandlingUnitOfWork(unitOfWork, storedEventHandler)

	bus := registerCommandHandlers(commonCommand.NewBusRegistry(logger), notifyingUnitOfWork)

	return &container{bus}, nil
}
