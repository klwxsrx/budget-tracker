package messaging

import (
	"time"
)

type MessageID []byte

type Message struct {
	ID        MessageID
	Type      string
	Data      []byte
	EventTime time.Time
}

type MessageHandler interface {
	Handle(msg Message) error
}

type CompositeTypedMessageHandler struct {
	handlers map[string]MessageHandler
}

func (h *CompositeTypedMessageHandler) Handle(msg Message) error {
	handler, ok := h.handlers[msg.Type]
	if !ok {
		return nil
	}
	return handler.Handle(msg)
}

func (h *CompositeTypedMessageHandler) Subscribe(messageType string, handler MessageHandler) {
	h.handlers[messageType] = handler
}

func NewCompositeTypedMessageHandler() *CompositeTypedMessageHandler {
	return &CompositeTypedMessageHandler{
		make(map[string]MessageHandler),
	}
}
