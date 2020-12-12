package infrastructure

import (
	commandCommon "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	commonEvent "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	eventMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/mysql"
	commonSerialization "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/serialization"
	commonMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/pulsar"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
)

type Container interface {
	CommandBus() commandCommon.Bus
}

type container struct {
	bus commandCommon.Bus
}

func (c *container) CommandBus() commandCommon.Bus {
	return c.bus
}

func registerCommandHandlers(bus commandCommon.BusRegistry, unitOfWork command.UnitOfWork) commandCommon.Bus {
	_ = bus.Register(command.NewCreateAccountHandler(unitOfWork))
	_ = bus.Register(command.NewRenameAccountHandler(unitOfWork))
	_ = bus.Register(command.NewDeleteAccountHandler(unitOfWork))
	return bus
}

func NewContainer(client commonMysql.TransactionalClient, broker pulsar.Connection, logger logger.Logger) (Container, error) {
	serializer := commonSerialization.NewSerializer()
	deserializer := serialization.NewEventDeserializer()
	unitOfWork := mysql.NewUnitOfWork(client, serializer, deserializer)

	_ = eventMysql.NewStore(client, serializer)             // TODO:
	var unsentEventProvider commonEvent.UnsentEventProvider // TODO:
	var eventBus commonEvent.Bus                            // TODO:
	storedEventHandler := commonEvent.NewStoredEventHandler(unsentEventProvider, eventBus, logger)
	notifyingUnitOfWork := event.NewStoredEventHandlingUnitOfWork(unitOfWork, storedEventHandler)

	bus := registerCommandHandlers(commandCommon.NewBusRegistry(logger), notifyingUnitOfWork)

	return &container{bus}, nil
}
