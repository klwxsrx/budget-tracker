package service

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/repository"
	commonStoredEvent "github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent" // TODO: install golangci & fix wrong aliases
	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type UnitOfWork interface {
	Execute(f func(r DomainRegistry) error) error
	Critical(lock string, f func(r DomainRegistry) error) error
}

type DomainRegistry interface {
	AccountListService() domain.AccountListService
}

type domainRegistry struct {
	accountListService domain.AccountListService
}

func (dr *domainRegistry) AccountListService() domain.AccountListService {
	return dr.accountListService
}

func registerEventHandlers(dispatcher event.Dispatcher, registry DomainRegistry, store commonStoredEvent.Store) event.Dispatcher {
	dispatcher.Subscribe(commonStoredEvent.NewStoreEventHandler(store))
	// TODO: add event handlers
	return dispatcher
}

func NewDomainRegistry(
	store commonStoredEvent.Store,
	deserializer commonStoredEvent.Deserializer,
) DomainRegistry {
	dispatcher := event.NewDispatcher()
	accountRepo := repository.NewAccountRepository(dispatcher, store, deserializer)
	registry := &domainRegistry{domain.NewAccountListService(accountRepo)}
	registerEventHandlers(dispatcher, registry, store)
	return registry
}
