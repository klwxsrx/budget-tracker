package app

import (
	domain "github.com/klwxsrx/expense-tracker/pkg/expense/domain/account"
)

type DomainRegistry interface {
	AccountService() domain.Service
}

type Transaction interface {
	Execute(f func(r DomainRegistry) error) error
	Critical(lock string, f func(r DomainRegistry) error) error
}
