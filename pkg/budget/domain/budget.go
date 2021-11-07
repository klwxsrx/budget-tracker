package domain

import "github.com/google/uuid"

const budgetAggregateName = "budget"

type BudgetID struct {
	uuid.UUID
}
