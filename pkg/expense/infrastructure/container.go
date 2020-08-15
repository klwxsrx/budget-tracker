package infrastructure

import (
	commandApp "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	eventApp "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
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
	transaction  command.Transaction
	bus          commandApp.Bus
	dispatcher   eventApp.Dispatcher
	serializer   eventApp.Serializer
	deserializer serialization.EventDeserializer
}

func (c *container) CommandBus() commandApp.Bus {
	return c.bus
}

func eventStoreEventHandler(client mysql.Client, es eventApp.Serializer) eventApp.Handler {
	eventStore := eventMysql.NewStore(client, es)
	return event.NewStoreEventHandler(eventStore)
}

func eventDispatcher(client mysql.Client, es eventApp.Serializer) eventApp.Dispatcher {
	d := eventApp.NewDispatcher()
	d.Subscribe(eventStoreEventHandler(client, es))
	return registerEventHandlers(d)
}

func eventSerializer() eventApp.Serializer {
	return eventSerialization.NewSerializer()
}

func eventDeserializer() serialization.EventDeserializer {
	return serialization.NewEventDeserializer()
}

func transaction(
	client mysql.TransactionalClient,
	ed eventApp.Dispatcher,
	es eventApp.Serializer,
	de serialization.EventDeserializer,
) command.Transaction {
	return mysqlInfrastructure.NewTransaction(client, ed, es, de)
}

func bus(tx command.Transaction) commandApp.Bus {
	b := commandApp.NewBusRegistry()
	return registerCommandHandlers(b, tx)
}

func registerEventHandlers(d eventApp.Dispatcher) eventApp.Dispatcher {
	// TODO: add event handlers
	return d
}

func registerCommandHandlers(b commandApp.BusRegistry, tx command.Transaction) commandApp.Bus {
	_ = b.Register(command.NewCreateAccountHandler(tx))
	_ = b.Register(command.NewRenameAccountHandler(tx))
	_ = b.Register(command.NewDeleteAccountHandler(tx))
	return b
}

func NewContainer(client mysql.TransactionalClient) Container {
	es := eventSerializer()
	de := eventDeserializer()
	ed := eventDispatcher(client, es)
	tx := transaction(client, ed, es, de)
	return &container{client, tx, bus(tx), ed, es, de}
}
