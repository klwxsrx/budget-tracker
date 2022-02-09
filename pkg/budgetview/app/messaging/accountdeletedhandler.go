package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
)

type accountDeletedJSON struct {
	BudgetID  uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"acc_id"`
}

type accountDeletedMessageHandler struct {
	service *service.AccountService
}

func (h *accountDeletedMessageHandler) MessageType() string {
	return AccountDeletedMessageType
}

func (h *accountDeletedMessageHandler) Handle(msg messaging.Message) error {
	if msg.Type != h.MessageType() {
		return fmt.Errorf("%w type %s for %T", messaging.ErrUnsupportedMessage, msg.Type, h)
	}

	var event accountDeletedJSON
	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		return fmt.Errorf("%w structure for %T: %v", messaging.ErrUnsupportedMessage, h, err)
	}

	return h.service.HandleAccountDeleted(event.BudgetID, event.AccountID)
}

func NewAccountDeletedMessageHandler(srv *service.AccountService) messaging.TypedMessageHandler {
	return &accountDeletedMessageHandler{srv}
}
