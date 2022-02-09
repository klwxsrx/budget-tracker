package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
)

type accountCreatedJSON struct {
	BudgetID       uuid.UUID `json:"id"`
	AccountID      uuid.UUID `json:"acc_id"`
	Title          string    `json:"title"`
	InitialBalance int       `json:"balance"`
}

type accountCreatedMessageHandler struct {
	service *service.AccountService
}

func (h *accountCreatedMessageHandler) MessageType() string {
	return AccountCreatedMessageType
}

func (h *accountCreatedMessageHandler) Handle(msg messaging.Message) error {
	if msg.Type != h.MessageType() {
		return fmt.Errorf("%w type %s for %T", messaging.ErrUnsupportedMessage, msg.Type, h)
	}

	var event accountCreatedJSON
	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		return fmt.Errorf("%w structure for %T: %v", messaging.ErrUnsupportedMessage, h, err)
	}

	return h.service.HandleAccountCreated(event.BudgetID, event.AccountID, event.Title, event.InitialBalance)
}

func NewAccountCreatedMessageHandler(srv *service.AccountService) messaging.TypedMessageHandler {
	return &accountCreatedMessageHandler{srv}
}
