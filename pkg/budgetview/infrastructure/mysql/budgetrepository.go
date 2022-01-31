package mysql

import (
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	commoninfrastructureuuid "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/uuid"

	"time"
)

type budgetRepository struct {
	client mysql.Client
}

func (repo *budgetRepository) Create(budget *model.Budget) error {
	const query = "INSERT INTO budget (id, title, currency, created_at, updated_at) VALUES (?, ?, ?, ?, ?)" +
		" ON DUPLICATE KEY UPDATE id = id"

	now := time.Now()
	binaryID := commoninfrastructureuuid.BinaryUUID(budget.ID)
	result, err := repo.client.Exec(query, binaryID, budget.Title, budget.Currency, now, now)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrBudgetAlreadyExists
	}
	return err
}

func NewBudgetRepository(client mysql.Client) model.BudgetRepository {
	return &budgetRepository{client}
}
