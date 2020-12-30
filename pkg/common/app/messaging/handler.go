package messaging

type MessageID []byte
type Message []byte

type MessageHandler interface {
	Handle(e Message) error
	SetOffset(id MessageID) error
	GetOffset() *MessageID
	GetName() string
}
