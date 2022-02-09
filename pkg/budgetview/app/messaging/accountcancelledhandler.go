package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
)

type accountCancelledJSON struct {
	AccountID uuid.UUID `json:"acc_id"`
}

type accountCancelledMessageHandler struct {
	service *service.AccountService
}

func (h *accountCancelledMessageHandler) MessageType() string {
	return AccountCancelledMessageType
}

func (h *accountCancelledMessageHandler) Handle(msg messaging.Message) error {
	if msg.Type != h.MessageType() {
		return fmt.Errorf("%w type %s for %T", messaging.ErrUnsupportedMessage, msg.Type, h)
	}

	var event accountCancelledJSON
	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		return fmt.Errorf("%w structure for %T: %v", messaging.ErrUnsupportedMessage, h, err)
	}

	return h.service.HandleAccountCancelled(event.AccountID)
}

func NewAccountCancelledMessageHandler(srv *service.AccountService) messaging.TypedMessageHandler {
	return &accountCancelledMessageHandler{srv}
}
