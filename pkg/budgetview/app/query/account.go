package query

import (
	"github.com/google/uuid"
)

const (
	AccountStatusActive    = "active"
	AccountStatusCancelled = "cancelled"
)

type Account struct {
	ID             uuid.UUID
	BudgetID       uuid.UUID
	Title          string
	Status         int
	InitialBalance int
	CurrentBalance int
}

type AccountQueryService interface {
	ListAccounts(budgetID uuid.UUID) ([]Account, error)
}
