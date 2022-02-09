package mysql

import (
	"fmt"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
)

type repoProvider struct {
	client mysql.Client
}

func (p *repoProvider) BudgetRepository() model.BudgetRepository {
	return NewBudgetRepository(p.client)
}

func (p *repoProvider) AccountRepository() model.AccountRepository {
	return NewAccountRepository(p.client)
}

func newRepoProvider(client mysql.Client) service.RepositoryProvider {
	return &repoProvider{client}
}

type unitOfWork struct {
	client mysql.TransactionalClient
}

func (uw *unitOfWork) Execute(f func(r service.RepositoryProvider) error) error {
	tx, err := uw.client.Begin()
	if err != nil {
		return fmt.Errorf("can't begin new transaction: %w", err)
	}

	repos := newRepoProvider(tx)
	err = f(repos)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func NewUnitOfWork(client mysql.TransactionalClient) service.UnitOfWork {
	return &unitOfWork{client}
}
