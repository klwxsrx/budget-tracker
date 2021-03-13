package command

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
)

type DomainRegistry interface {
	AccountService() domain.AccountService
}

type UnitOfWork interface {
	Execute(f func(r DomainRegistry) error) error
	Critical(lock string, f func(r DomainRegistry) error) error
}
