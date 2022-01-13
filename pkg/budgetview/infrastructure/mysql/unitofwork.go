package mysql

import (
	"fmt"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
)

type repoProvider struct{}

type unitOfWork struct {
	client mysql.TransactionalClient
}

func (uw *unitOfWork) Execute(f func(r service.RepositoryProvider) error) error {
	tx, err := uw.client.Begin()
	if err != nil {
		return fmt.Errorf("can't begin new transaction: %w", err)
	}

	registry := &repoProvider{}
	err = f(registry)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func NewUnitOfWork(client mysql.TransactionalClient) service.UnitOfWork {
	return &unitOfWork{client}
}
