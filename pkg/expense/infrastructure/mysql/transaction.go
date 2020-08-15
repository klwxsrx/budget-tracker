package mysql

import (
	"fmt"
	eventApp "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
)

type transaction struct {
	client       mysql.TransactionalClient
	dispatcher   eventApp.Dispatcher
	serializer   eventApp.Serializer
	deserializer serialization.EventDeserializer
}

func (t *transaction) Execute(f func(r command.DomainRegistry) error) error {
	tx, err := t.client.Begin()
	if err != nil {
		return fmt.Errorf("can't begin new transaction: %v", err)
	}
	registry := newDomainRegistry(tx, t.dispatcher, t.serializer, t.deserializer)
	err = f(registry)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (t *transaction) Critical(lock string, f func(r command.DomainRegistry) error) error {
	dbLock := mysql.NewLock(t.client, lock)
	err := dbLock.Get()
	if err != nil {
		return fmt.Errorf("can't create lock %v", err)
	}
	defer dbLock.Release()
	return t.Execute(f)
}

func NewTransaction(
	client mysql.TransactionalClient,
	dispatcher eventApp.Dispatcher,
	serializer eventApp.Serializer,
	deserializer serialization.EventDeserializer,
) command.Transaction {
	return &transaction{client, dispatcher, serializer, deserializer}
}
