package persistence

import (
	"fmt"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/realtime"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
	commonapprealtime "github.com/klwxsrx/budget-tracker/pkg/common/app/realtime"
	commoninfrastructuremysql "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
)

type repoProvider struct {
	budgetRepo      model.BudgetRepository
	accountRepo     model.AccountRepository
	realtimeService realtime.Service
}

func (p *repoProvider) BudgetRepository() model.BudgetRepository {
	return p.budgetRepo
}

func (p *repoProvider) AccountRepository() model.AccountRepository {
	return p.accountRepo
}

func (p *repoProvider) RealtimeService() realtime.Service {
	return p.realtimeService
}

func newRepoProvider(
	budgetRepo model.BudgetRepository,
	accountRepo model.AccountRepository,
	realtimeService realtime.Service,
) service.RepositoryProvider {
	return &repoProvider{
		budgetRepo,
		accountRepo,
		realtimeService,
	}
}

type unitOfWork struct {
	mysqlClient    commoninfrastructuremysql.TransactionalClient
	realtimeClient commonapprealtime.Client
	logger         log.Logger
}

func (uw *unitOfWork) Execute(f func(r service.RepositoryProvider) error) error {
	tx, err := uw.mysqlClient.Begin()
	if err != nil {
		return fmt.Errorf("can't begin new transaction: %w", err)
	}

	scheduledRealtimeClient := NewScheduledRealtimeClient(uw.realtimeClient, uw.logger)
	repos := newRepoProvider(
		mysql.NewBudgetRepository(tx),
		mysql.NewAccountRepository(tx),
		realtime.NewService(scheduledRealtimeClient, uw.logger),
	)
	err = f(repos)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err == nil {
		scheduledRealtimeClient.Commit()
	}
	return err
}

func NewUnitOfWork(
	mysqlClient commoninfrastructuremysql.TransactionalClient,
	realtimeClient commonapprealtime.Client,
	logger log.Logger,
) service.UnitOfWork {
	return &unitOfWork{mysqlClient, realtimeClient, logger}
}
