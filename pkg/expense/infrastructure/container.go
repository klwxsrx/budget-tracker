package infrastructure

import (
	commandCommon "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/amqp"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/messaging"
	eventMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/mysql"
	commonSerialization "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/serialization"
	commonMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
)

type Container interface {
	CommandBus() commandCommon.Bus
	EventNotifierChannel() amqp.Channel
}

type container struct {
	bus           commandCommon.Bus
	eventNotifier amqp.Channel
}

func (c *container) CommandBus() commandCommon.Bus {
	return c.bus
}

func (c *container) EventNotifierChannel() amqp.Channel {
	return c.eventNotifier
}

func registerCommandHandlers(bus commandCommon.BusRegistry, unitOfWork command.UnitOfWork) commandCommon.Bus {
	_ = bus.Register(command.NewCreateAccountHandler(unitOfWork))
	_ = bus.Register(command.NewRenameAccountHandler(unitOfWork))
	_ = bus.Register(command.NewDeleteAccountHandler(unitOfWork))
	return bus
}

func NewContainer(client commonMysql.TransactionalClient, logger logger.Logger) Container {
	serializer := commonSerialization.NewSerializer()
	deserializer := serialization.NewEventDeserializer()
	unitOfWork := mysql.NewUnitOfWork(client, serializer, deserializer)

	eventStore := eventMysql.NewStore(client, serializer)
	storedEventNotifier := messaging.NewStoredEventNotifier(eventStore, client)
	notificationDispatcher := event.NewStoredEventNotificationDispatcher(storedEventNotifier, logger)
	notifyingUnitOfWork := mysql.NewNotifyingUnitOfWork(unitOfWork, notificationDispatcher)

	bus := registerCommandHandlers(commandCommon.NewBusRegistry(logger), notifyingUnitOfWork)

	return &container{bus, storedEventNotifier}
}
