package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
)

type accountRenamedJSON struct {
	AccountID uuid.UUID `json:"acc_id"`
	Title     string    `json:"title"`
}

type accountRenamedMessageHandler struct {
	service *service.AccountService
}

func (h *accountRenamedMessageHandler) MessageType() string {
	return AccountRenamedMessageType
}

func (h *accountRenamedMessageHandler) Handle(msg messaging.Message) error {
	if msg.Type != h.MessageType() {
		return fmt.Errorf("%w type %s for %T", messaging.ErrUnsupportedMessage, msg.Type, h)
	}

	var event accountRenamedJSON
	err := json.Unmarshal(msg.Data, &event)
	if err != nil {
		return fmt.Errorf("%w structure for %T: %v", messaging.ErrUnsupportedMessage, h, err)
	}

	return h.service.HandleAccountRenamed(event.AccountID, event.Title)
}

func NewAccountRenamedMessageHandler(srv *service.AccountService) messaging.TypedMessageHandler {
	return &accountRenamedMessageHandler{srv}
}
