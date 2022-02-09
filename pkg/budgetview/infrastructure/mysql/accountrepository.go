package mysql

import (
	sql2 "database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	commoninfrastructureuuid "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/uuid"
)

type accountRepository struct {
	client mysql.Client
}

func (r *accountRepository) FindByBudgetID(id uuid.UUID) ([]*model.Account, error) {
	const sql = "SELECT * FROM account WHERE budget_id = ? ORDER BY position ASC"

	var accounts []sqlxAccount
	err := r.client.Select(&accounts, sql, commoninfrastructureuuid.BinaryUUID(id))
	if err != nil {
		return nil, err
	}

	result := make([]*model.Account, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, &model.Account{
			BudgetID:       uuid.UUID(account.BudgetID),
			AccountID:      uuid.UUID(account.ID),
			Title:          account.Title,
			Status:         account.Status,
			InitialBalance: account.InitialBalance,
			CurrentBalance: account.CurrentBalance,
			Position:       account.Position,
		})
	}
	return result, nil
}

func (r *accountRepository) Create(account *model.Account) error {
	const query = "INSERT INTO account" +
		"(id, budget_id, status, title, initial_balance, current_balance, position, created_at, updated_at)" +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE id = id"

	now := time.Now()
	binaryID := commoninfrastructureuuid.BinaryUUID(account.AccountID)
	binaryBudgetID := commoninfrastructureuuid.BinaryUUID(account.BudgetID)

	result, err := r.client.Exec(query,
		binaryID,
		binaryBudgetID,
		account.Status,
		account.Title,
		account.InitialBalance,
		account.CurrentBalance,
		account.Position,
		now,
		now,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrAccountAlreadyExists
	}
	return err
}

func (r *accountRepository) Update(account *model.Account) error {
	const query = "UPDATE account SET " +
		"budget_id = ?, status = ?, title = ?, initial_balance = ?, current_balance = ?, position = ?, updated_at = ? " +
		"WHERE id = ?"

	now := time.Now()
	binaryID := commoninfrastructureuuid.BinaryUUID(account.AccountID)
	binaryBudgetID := commoninfrastructureuuid.BinaryUUID(account.BudgetID)

	result, err := r.client.Exec(query,
		binaryBudgetID,
		account.Status,
		account.Title,
		account.InitialBalance,
		account.CurrentBalance,
		account.Position,
		now,
		binaryID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrAccountDoesNotExist
	}
	return err
}

func (r *accountRepository) FindByID(id uuid.UUID) (*model.Account, error) {
	const sql = "SELECT * FROM account WHERE id = ?"

	var account sqlxAccount
	err := r.client.Get(&account, sql, commoninfrastructureuuid.BinaryUUID(id))
	if errors.Is(err, sql2.ErrNoRows) {
		return nil, model.ErrAccountDoesNotExist
	} else if err != nil {
		return nil, err
	}

	return &model.Account{
		BudgetID:       uuid.UUID(account.BudgetID),
		AccountID:      uuid.UUID(account.ID),
		Title:          account.Title,
		Status:         account.Status,
		InitialBalance: account.InitialBalance,
		CurrentBalance: account.CurrentBalance,
		Position:       account.Position,
	}, nil
}

func (r *accountRepository) Delete(id uuid.UUID) error {
	const sql = "DELETE FROM account WHERE id = ?"

	binaryID := commoninfrastructureuuid.BinaryUUID(id)
	result, err := r.client.Exec(sql, binaryID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrAccountDoesNotExist
	}
	return nil
}

func NewAccountRepository(client mysql.Client) model.AccountRepository {
	return &accountRepository{client}
}
