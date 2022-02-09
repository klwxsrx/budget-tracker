package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
)

type accountReorderedJSON struct {
	BudgetID  uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"acc_id"`
	Position  int       `json:"pos"`
}

type accountReorderedMessageHandler struct {
	service *service.AccountService
}

func (h *accountReorderedMessageHandler) MessageType() string {
	return AccountReorderedMessageType
}

func (h *accountReorderedMessageHandler) Handle(msg messaging.Message) error {
	if msg.Type != h.MessageType() {
		return fmt.Errorf("%w type %s for %T", messaging.ErrUnsupportedMessage, msg.Type, h)
	}

	var event accountReorderedJSON
	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		return fmt.Errorf("%w structure for %T: %v", messaging.ErrUnsupportedMessage, h, err)
	}

	return h.service.HandleAccountReordered(event.BudgetID, event.AccountID, event.Position)
}

func NewAccountReorderedMessageHandler(srv *service.AccountService) messaging.TypedMessageHandler {
	return &accountReorderedMessageHandler{srv}
}
