package realtime

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/realtime"
)

type Service interface {
	BudgetCreated(id uuid.UUID, title, currency string)
}

const (
	realtimeBudgetsChannelName = "budgets"
)

type service struct {
	client realtime.Client
	logger log.Logger
}

func (s *service) BudgetCreated(id uuid.UUID, title, currency string) {
	message := struct {
		baseJSON
		Budget budgetJSON `json:"budget"`
	}{baseJSON{"budget_created"}, budgetJSON{id, title, currency}}
	s.tryPublishMessage(realtimeBudgetsChannelName, message)
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
