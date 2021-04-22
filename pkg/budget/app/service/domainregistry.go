package service

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/repository"
	commonStoredEvent "github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type domainRegistry struct {
	accountService domain.AccountService
}

func (dr *domainRegistry) AccountService() domain.AccountService {
	return dr.accountService
}

func registerEventHandlers(dispatcher event.Dispatcher, registry command.DomainRegistry, store commonStoredEvent.Store) event.Dispatcher {
	dispatcher.Subscribe(commonStoredEvent.NewStoreEventHandler(store))
	// TODO: add event handlers
	return dispatcher
}

func NewDomainRegistry(
	store commonStoredEvent.Store,
	deserializer storedevent.Deserializer,
) command.DomainRegistry {
	dispatcher := event.NewDispatcher()
	accountRepo := repository.NewAccountRepository(dispatcher, store, deserializer)
	registry := &domainRegistry{domain.NewAccountService(accountRepo)}
	registerEventHandlers(dispatcher, registry, store)
	return registry
}
