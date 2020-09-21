package infrastructure

import (
	commandApp "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	eventApp "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event"
	eventMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/mysql"
	eventSerialization "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/serialization"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	mysqlInfrastructure "github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
)

type Container interface {
	CommandBus() commandApp.Bus
}

type container struct {
	client       mysql.TransactionalClient
	unitOfWork   command.UnitOfWork
	bus          commandApp.Bus
	dispatcher   eventApp.Dispatcher
	serializer   eventApp.Serializer
	deserializer serialization.EventDeserializer
	logger       logger.Logger
}

func (c *container) CommandBus() commandApp.Bus {
	return c.bus
}

func eventStoreEventHandler(client mysql.Client, serializer eventApp.Serializer) eventApp.Handler {
	eventStore := eventMysql.NewStore(client, serializer)
	return event.NewStoreEventHandler(eventStore)
}

func eventDispatcher(client mysql.Client, serializer eventApp.Serializer) eventApp.Dispatcher {
	dispatcher := eventApp.NewDispatcher()
	dispatcher.Subscribe(eventStoreEventHandler(client, serializer)) // TODO: subscribe within DomainRegistry
	return registerEventHandlers(dispatcher)
}

func eventSerializer() eventApp.Serializer {
	return eventSerialization.NewSerializer()
}

func eventDeserializer() serialization.EventDeserializer {
	return serialization.NewEventDeserializer()
}

func unitOfWork(
	client mysql.TransactionalClient,
	dispatcher eventApp.Dispatcher,
	serializer eventApp.Serializer,
	deserializer serialization.EventDeserializer,
) command.UnitOfWork {
	return mysqlInfrastructure.NewUnitOfWork(client, dispatcher, serializer, deserializer)
}

func bus(unitOfWork command.UnitOfWork, logger logger.Logger) commandApp.Bus {
	bus := commandApp.NewBusRegistry(logger)
	return registerCommandHandlers(bus, unitOfWork)
}

func registerEventHandlers(dispatcher eventApp.Dispatcher) eventApp.Dispatcher {
	// TODO: add event handlers
	return dispatcher
}

func registerCommandHandlers(bus commandApp.BusRegistry, unitOfWork command.UnitOfWork) commandApp.Bus {
	_ = bus.Register(command.NewCreateAccountHandler(unitOfWork))
	_ = bus.Register(command.NewRenameAccountHandler(unitOfWork))
	_ = bus.Register(command.NewDeleteAccountHandler(unitOfWork))
	return bus
}

func NewContainer(client mysql.TransactionalClient, logger logger.Logger) Container {
	serializer := eventSerializer()
	deserializer := eventDeserializer()
	dispatcher := eventDispatcher(client, serializer)
	unitOfWork := unitOfWork(client, dispatcher, serializer, deserializer)
	return &container{client, unitOfWork, bus(unitOfWork, logger), dispatcher, serializer, deserializer, logger}
}
