package messaging

import (
	"errors"
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

var ErrUnsupportedMessage = errors.New("unsupported message")

type TypedMessageHandler interface {
	MessageHandler
	MessageType() string
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

func (h *CompositeTypedMessageHandler) SubscribeTyped(handler TypedMessageHandler) {
	h.handlers[handler.MessageType()] = handler
}

func NewCompositeTypedMessageHandler() *CompositeTypedMessageHandler {
	return &CompositeTypedMessageHandler{
		make(map[string]MessageHandler),
	}
}
