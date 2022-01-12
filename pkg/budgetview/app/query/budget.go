package query

import "github.com/google/uuid"

type Budget struct {
	ID       uuid.UUID
	Title    string
	Currency string
}

type BudgetQueryService interface {
	ListBudgets() ([]Budget, error)
	ExistByIDs(ids []uuid.UUID) (bool, error)
}
