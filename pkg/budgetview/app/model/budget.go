package model

import (
	"errors"

	"github.com/google/uuid"
)

var ErrBudgetAlreadyExists = errors.New("budget already exists")

type BudgetRepository interface {
	Create(budget *Budget) error
}

type Budget struct {
	ID       uuid.UUID
	Title    string
	Currency string
}
