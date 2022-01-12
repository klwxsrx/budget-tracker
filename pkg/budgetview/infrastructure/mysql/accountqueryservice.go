package mysql

import (
	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/query"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	commoninfrastructureuuid "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/uuid"
)

type accountQueryService struct {
	client mysql.Client
}

func (s *accountQueryService) ListAccounts(budgetID uuid.UUID) ([]query.Account, error) {
	var accounts []sqlxAccount
	err := s.client.Select(&accounts, "SELECT * FROM account WHERE budget_id = ?", commoninfrastructureuuid.BinaryUUID(budgetID))
	if err != nil {
		return nil, err
	}

	result := make([]query.Account, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, query.Account{
			ID:             uuid.UUID(account.ID),
			BudgetID:       uuid.UUID(account.BudgetID),
			Title:          account.Title,
			Status:         account.Status,
			InitialBalance: account.InitialBalance,
			CurrentBalance: account.CurrentBalance,
		})
	}
	return result, nil
}

func NewAccountQueryService(client mysql.Client) query.AccountQueryService {
	return &accountQueryService{client}
}

type sqlxAccount struct {
	ID             commoninfrastructureuuid.BinaryUUID `db:"id"`
	BudgetID       commoninfrastructureuuid.BinaryUUID `db:"budget_id"`
	Title          string                              `db:"title"`
	Status         int                                 `db:"status"`
	InitialBalance int                                 `db:"initial_balance"`
	CurrentBalance int                                 `db:"current_balance"`
}
