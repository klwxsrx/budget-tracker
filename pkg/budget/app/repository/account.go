package repository

import (
	"errors"

	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commondomain "github.com/klwxsrx/budget-tracker/pkg/common/domain"
)

type accountListRepository struct {
	aggregateRepository
}

func (r *accountListRepository) Update(list *domain.AccountList) error {
	return r.storeChanges(list)
}

func (r *accountListRepository) FindByID(id domain.BudgetID) (*domain.AccountList, error) {
	state := &domain.AccountListState{}
	err := r.loadChanges(commondomain.AggregateID(id), state)
	if errors.Is(err, errAggregateNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return domain.LoadAccountList(state), nil
}

func NewAccountRepository(
	store storedevent.Store,
	deserializer storedevent.Deserializer,
) domain.AccountListRepository {
	return &accountListRepository{aggregateRepository{
		store:        store,
		deserializer: deserializer,
	}}
}
