package mysql

import (
	"fmt"

	"github.com/klwxsrx/budget-tracker/pkg/budget/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
)

type unitOfWork struct {
	client       mysql.TransactionalClient
	serializer   messaging.DomainEventSerializer
	deserializer messaging.DomainEventDeserializer
}

func (uw *unitOfWork) Execute(f func(r service.DomainRegistry) error) error {
	tx, err := uw.client.Begin()
	if err != nil {
		return fmt.Errorf("can't begin new transaction: %w", err)
	}

	registry := service.NewDomainRegistry(uw.newStore(tx), uw.deserializer)
	err = f(registry)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (uw *unitOfWork) Critical(lock string, f func(r service.DomainRegistry) error) error {
	dbLock := mysql.NewLock(uw.client, lock)
	err := dbLock.Get()
	if err != nil {
		return fmt.Errorf("can't create lock %v: %w", lock, err)
	}
	defer dbLock.Release()
	return uw.Execute(f)
}

func (uw *unitOfWork) newStore(client mysql.Client) storedevent.Store {
	store := mysql.NewEventStore(client, uw.serializer)
	return mysql.NewUnsentEventStoreDecorator(client, store)
}

func NewUnitOfWork(
	client mysql.TransactionalClient,
	serializer messaging.DomainEventSerializer,
	deserializer messaging.DomainEventDeserializer,
) service.UnitOfWork {
	return &unitOfWork{client, serializer, deserializer}
}
