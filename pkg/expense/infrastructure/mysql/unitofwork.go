package mysql

import (
	"fmt"
	eventApp "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
)

type unitOfWork struct {
	client       mysql.TransactionalClient
	dispatcher   eventApp.Dispatcher
	serializer   eventApp.Serializer
	deserializer serialization.EventDeserializer
}

func (uw *unitOfWork) Execute(f func(r command.DomainRegistry) error) error {
	tx, err := uw.client.Begin()
	if err != nil {
		return fmt.Errorf("can't begin new transaction: %v", err)
	}
	registry := newDomainRegistry(tx, uw.dispatcher, uw.serializer, uw.deserializer)
	err = f(registry)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (uw *unitOfWork) Critical(lock string, f func(r command.DomainRegistry) error) error {
	dbLock := mysql.NewLock(uw.client, lock)
	err := dbLock.Get()
	if err != nil {
		return fmt.Errorf("can't create lock: %v", err)
	}
	defer dbLock.Release()
	return uw.Execute(f)
}

func NewUnitOfWork(
	client mysql.TransactionalClient,
	dispatcher eventApp.Dispatcher,
	serializer eventApp.Serializer,
	deserializer serialization.EventDeserializer,
) command.UnitOfWork {
	return &unitOfWork{client, dispatcher, serializer, deserializer}
}
