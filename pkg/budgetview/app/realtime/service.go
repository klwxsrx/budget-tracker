package realtime

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/realtime"
)

type Service interface {
	BudgetCreated(budget *model.Budget)
	AccountCreated(budgetID uuid.UUID, account *model.Account)
	AccountRenamed(budgetID, accountID uuid.UUID, title string)
	AccountStatusChanged(budgetID, accountID uuid.UUID, status int)
	AccountDeleted(budgetID, accountID uuid.UUID)
	AccountsReordered(budgetID uuid.UUID, accountIDs []uuid.UUID)
}

const (
	realtimeBudgetsChannelName        = "budgets"
	realtimeBudgetChannelNameTemplate = "budget_%s"

	accountStatusActive    = "active"
	accountStatusCancelled = "cancelled"
)

type service struct {
	client realtime.Client
	logger log.Logger
}

func (s *service) BudgetCreated(budget *model.Budget) {
	message := struct {
		baseJSON
		Budget budgetJSON `json:"budget"`
	}{baseJSON{"budget_created"}, budgetJSON{budget.ID, budget.Title, budget.Currency}}

	s.tryPublishMessage(realtimeBudgetsChannelName, message)
}

func (s *service) AccountCreated(budgetID uuid.UUID, account *model.Account) {
	message := struct {
		baseJSON
		Account accountJSON `json:"account"`
	}{
		baseJSON{"account_created"},
		accountJSON{
			account.AccountID,
			account.Title,
			s.getStringAccountStatus(account.Status),
			account.InitialBalance,
			account.CurrentBalance,
			account.Position,
		},
	}

	s.tryPublishMessage(s.budgetChannelName(budgetID), message)
}

func (s *service) AccountRenamed(budgetID, accountID uuid.UUID, title string) {
	message := struct {
		baseJSON
		AccountID uuid.UUID `json:"account_id"`
		Title     string    `json:"title"`
	}{baseJSON{"account_renamed"}, accountID, title}

	s.tryPublishMessage(s.budgetChannelName(budgetID), message)
}

func (s *service) AccountStatusChanged(budgetID, accountID uuid.UUID, status int) {
	message := struct {
		baseJSON
		AccountID uuid.UUID `json:"account_id"`
		Status    string    `json:"status"`
	}{baseJSON{"account_status_changed"}, accountID, s.getStringAccountStatus(status)}

	s.tryPublishMessage(s.budgetChannelName(budgetID), message)
}

func (s *service) AccountDeleted(budgetID, accountID uuid.UUID) {
	message := struct {
		baseJSON
		AccountID uuid.UUID `json:"account_id"`
	}{baseJSON{"account_deleted"}, accountID}

	s.tryPublishMessage(s.budgetChannelName(budgetID), message)
}

func (s *service) AccountsReordered(budgetID uuid.UUID, accountIDs []uuid.UUID) {
	message := struct {
		baseJSON
		OrderedIDs []uuid.UUID `json:"ordered_ids"`
	}{baseJSON{"accounts_reordered"}, accountIDs}

	s.tryPublishMessage(s.budgetChannelName(budgetID), message)
}

func (s *service) tryPublishMessage(channel string, jsonMessage interface{}) {
	data, err := json.Marshal(jsonMessage)
	if err != nil {
		s.logger.WithError(err).Error(fmt.Sprintf("failed to encode json of struct: %v", jsonMessage))
	}
	err = s.client.PublishMessage(channel, data)
	if err != nil {
		s.logger.WithError(err).Error("failed to publish message")
	} else {
		s.logger.Info(fmt.Sprintf("realtime msg has been sent to the channel \"%s\"", channel))
	}
}

func (s *service) getStringAccountStatus(status int) string {
	switch status {
	case model.AccountStatusActive:
		return accountStatusActive
	case model.AccountStatusCancelled:
	}
	return accountStatusCancelled
}

func (s *service) budgetChannelName(budgetID uuid.UUID) string {
	return fmt.Sprintf(realtimeBudgetChannelNameTemplate, budgetID)
}

func NewService(client realtime.Client, logger log.Logger) Service {
	return &service{client, logger}
}

type baseJSON struct {
	Type string `json:"type"`
}

type budgetJSON struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Currency string    `json:"currency"`
}

type accountJSON struct {
	AccountID      uuid.UUID `json:"account_id"`
	Title          string    `json:"title"`
	Status         string    `json:"status"`
	InitialBalance int       `json:"initial_balance"`
	CurrentBalance int       `json:"current_balance"`
	Position       int       `json:"position"`
}
