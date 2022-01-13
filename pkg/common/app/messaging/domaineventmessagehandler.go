package messaging

import "github.com/klwxsrx/budget-tracker/pkg/common/domain"

type domainEventMessageHandler struct {
	eventHandler domain.EventHandler
	deserializer DomainEventDeserializer
}

func (handler *domainEventMessageHandler) Handle(msg Message) error {
	event, err := handler.deserializer.Deserialize(msg.Type, msg.Data)
	if err != nil {
		return err
	}
	return handler.eventHandler.Handle(event)
}

func NewDomainEventMessageHandler(handler domain.EventHandler, deserializer DomainEventDeserializer) MessageHandler {
	return &domainEventMessageHandler{handler, deserializer}
}
