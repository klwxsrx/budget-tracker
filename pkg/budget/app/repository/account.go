package repository

import (
	"errors"

	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	commondomain "github.com/klwxsrx/budget-tracker/pkg/common/domain"
)

type accountListRepository struct {
	aggregateRepo *AggregateRepository
}

func (r *accountListRepository) Update(list *domain.AccountList) error {
	return r.aggregateRepo.storeChanges(list)
}

func (r *accountListRepository) FindByID(id domain.BudgetID) (*domain.AccountList, error) {
	state := &domain.AccountListState{}
	err := r.aggregateRepo.loadChanges(commondomain.AggregateID(id), state)
	if errors.Is(err, errAggregateNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return domain.LoadAccountList(state), nil
}

func NewAccountRepository(
	aggregateRepo *AggregateRepository,
) domain.AccountListRepository {
	return &accountListRepository{aggregateRepo}
}
