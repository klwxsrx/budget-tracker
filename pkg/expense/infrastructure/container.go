package infrastructure

import (
	commandCommon "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	commonSerialization "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/serialization"
	commonMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
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

func bus(unitOfWork command.UnitOfWork, logger logger.Logger) commandCommon.Bus {
	bus := commandCommon.NewBusRegistry(logger)
	return registerCommandHandlers(bus, unitOfWork)
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
	return &container{bus(unitOfWork, logger)}
}
