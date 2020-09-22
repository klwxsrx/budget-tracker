package mysql

import (
	eventApp "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event"
	eventMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
)

type domainRegistry struct {
	accountService domain.AccountService
}

func (dr *domainRegistry) AccountService() domain.AccountService {
	return dr.accountService
}

func registerEventHandlers(dispatcher eventApp.Dispatcher, registry command.DomainRegistry, store eventApp.Store) eventApp.Dispatcher {
	dispatcher.Subscribe(event.NewStoreEventHandler(store))
	// TODO: add event handlers
	return dispatcher
}

func newDomainRegistry(
	client mysql.Client,
	serializer eventApp.Serializer,
	deserializer serialization.EventDeserializer,
) command.DomainRegistry {
	dispatcher := eventApp.NewDispatcher()
	store := eventMysql.NewStore(client, serializer)
	accountRepo := NewAccountRepository(dispatcher, store, deserializer)
	registry := &domainRegistry{domain.NewAccountService(accountRepo)}
	registerEventHandlers(dispatcher, registry, store)
	return registry
}
