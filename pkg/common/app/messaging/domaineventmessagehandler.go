package messaging

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/event"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

type domainEventMessageHandler struct {
	eventHandler event.DomainEventHandler
	deserializer storedevent.Deserializer
}

func (handler *domainEventMessageHandler) Handle(msg Message) error {
	evt, err := handler.deserializer.Deserialize(msg.Type, msg.Data)
	if err != nil {
		return err
	}
	return handler.eventHandler.Handle(evt)
}

func NewDomainEventMessageHandler(handler event.DomainEventHandler, deserializer storedevent.Deserializer) MessageHandler {
	return &domainEventMessageHandler{handler, deserializer}
}
