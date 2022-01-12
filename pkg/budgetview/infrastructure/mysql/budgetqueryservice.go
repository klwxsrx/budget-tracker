package mysql

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/query"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	commoninfrastructureuuid "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/uuid"
)

type budgetQueryService struct {
	client mysql.Client
}

func (s *budgetQueryService) ListBudgets() ([]query.Budget, error) {
	var budgets []sqlxBudget
	err := s.client.Select(&budgets, "SELECT * FROM budget")
	if err != nil {
		return nil, err
	}

	result := make([]query.Budget, 0, len(budgets))
	for _, budget := range budgets {
		result = append(result, query.Budget{
			ID:       uuid.UUID(budget.ID),
			Title:    budget.Title,
			Currency: budget.Currency,
		})
	}
	return result, nil
}

func (s *budgetQueryService) ExistByIDs(ids []uuid.UUID) (bool, error) {
	if ids == nil {
		return false, nil
	}

	binaryIDs := make([]commoninfrastructureuuid.BinaryUUID, 0, len(ids))
	for _, id := range ids {
		binaryIDs = append(binaryIDs, commoninfrastructureuuid.BinaryUUID(id))
	}

	var count int
	sql, params, err := sqlx.In("SELECT COUNT(*) FROM budget WHERE id IN (?)", binaryIDs)
	if err != nil {
		return false, err
	}

	err = s.client.Get(&count, sql, params...)
	if err != nil {
		return false, err
	}

	return count == len(ids), nil
}

func NewBudgetQueryService(client mysql.Client) query.BudgetQueryService {
	return &budgetQueryService{client}
}

type sqlxBudget struct {
	ID       commoninfrastructureuuid.BinaryUUID `db:"id"`
	Title    string                              `db:"title"`
	Currency string                              `db:"currency"`
}
