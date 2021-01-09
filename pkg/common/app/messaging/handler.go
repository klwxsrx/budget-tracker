package messaging

import "time"

type MessageID []byte

type Message struct {
	ID        MessageID
	Type      string
	Data      []byte
	EventTime time.Time
}

type MessageHandler interface {
	Handle(e *Message) error
}

type NamedMessageHandler interface {
	MessageHandler
	Name() string
	LatestMessageID() (*MessageID, error)
}
