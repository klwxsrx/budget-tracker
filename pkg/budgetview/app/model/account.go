package model

import (
	"errors"

	"github.com/google/uuid"
)

const (
	AccountStatusActive int = iota
	AccountStatusCancelled
)

var ErrAccountAlreadyExists = errors.New("account already exists")
var ErrAccountDoesNotExist = errors.New("account does not exist")

type AccountRepository interface {
	FindByID(id uuid.UUID) (*Account, error)
	FindByBudgetID(id uuid.UUID) ([]*Account, error)
	Create(account *Account) error
	Update(account *Account) error
	Delete(id uuid.UUID) error
}

type Account struct {
	BudgetID       uuid.UUID
	AccountID      uuid.UUID
	Title          string
	Status         int
	InitialBalance int
	CurrentBalance int
	Position       int
}
