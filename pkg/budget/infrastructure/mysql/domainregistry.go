package mysql

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/serialization"
	eventApp "github.com/klwxsrx/budget-tracker/pkg/common/app/event"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
)

type domainRegistry struct {
	accountService domain.AccountService
}

func (dr *domainRegistry) AccountService() domain.AccountService {
	return dr.accountService
}

func registerEventHandlers(dispatcher eventApp.Dispatcher, registry command.DomainRegistry, store storedevent.Store) eventApp.Dispatcher {
	dispatcher.Subscribe(storedevent.NewStoreEventHandler(store))
	// TODO: add event handlers
	return dispatcher
}

func newDomainRegistry(
	client mysql.Client,
	serializer eventApp.Serializer,
	deserializer serialization.EventDeserializer,
) command.DomainRegistry {
	dispatcher := eventApp.NewDispatcher()
	store := mysql.NewStore(client, serializer)
	accountRepo := NewAccountRepository(dispatcher, store, deserializer)
	registry := &domainRegistry{domain.NewAccountService(accountRepo)}
	registerEventHandlers(dispatcher, registry, store)
	return registry
}
