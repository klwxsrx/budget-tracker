package app

import "github.com/klwxsrx/expense-tracker/pkg/expense/domain"

type DomainRegistry interface {
	AccountService() domain.AccountService
}

type Transaction interface {
	Execute(f func(r DomainRegistry) error) error
	Critical(lock string, f func(r DomainRegistry) error) error
}
