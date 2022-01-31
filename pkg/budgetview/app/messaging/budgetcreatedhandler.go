package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
)

type budgetCreatedJSON struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Currency string    `json:"currency"`
}

type budgetCreatedMessageHandler struct {
	service *service.BudgetService
}

func (h *budgetCreatedMessageHandler) MessageType() string {
	return BudgetCreatedMessageType
}

func (h *budgetCreatedMessageHandler) Handle(msg messaging.Message) error {
	if msg.Type != h.MessageType() {
		return fmt.Errorf("%w type %s for %T", messaging.ErrUnsupportedMessage, msg.Type, h)
	}

	var event budgetCreatedJSON
	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		return fmt.Errorf("%w structure for %T: %v", messaging.ErrUnsupportedMessage, h, err)
	}

	return h.service.HandleBudgetCreated(event.ID, event.Title, event.Currency)
}

func NewBudgetCreatedMessageHandler(srv *service.BudgetService) messaging.TypedMessageHandler {
	return &budgetCreatedMessageHandler{srv}
}
