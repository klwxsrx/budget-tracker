package service

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/repository"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
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

func NewDomainRegistry(
	store storedevent.Store,
	deserializer messaging.DomainEventDeserializer,
) DomainRegistry {
	aggregateRepo := repository.NewAggregateRepository(store, deserializer)
	accountRepo := repository.NewAccountRepository(aggregateRepo)
	return &domainRegistry{domain.NewAccountListService(accountRepo)}
}
