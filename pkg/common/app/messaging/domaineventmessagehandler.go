package messaging

import (
	commonappevent "github.com/klwxsrx/budget-tracker/pkg/common/app/event"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

type domainEventMessageHandler struct {
	eventHandler commonappevent.DomainEventHandler
	deserializer storedevent.Deserializer
}

func (handler *domainEventMessageHandler) Handle(msg Message) error {
	event, err := handler.deserializer.Deserialize(msg.Type, msg.Data)
	if err != nil {
		return err
	}
	return handler.eventHandler.Handle(event)
}

func NewDomainEventMessageHandler(handler commonappevent.DomainEventHandler, deserializer storedevent.Deserializer) MessageHandler {
	return &domainEventMessageHandler{handler, deserializer}
}
