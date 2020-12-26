package messaging

type ID []byte
type Message []byte

type Subscriber interface {
	Handle(e Message) error
	SetOffset(id ID) error
	GetOffset() *ID
	GetName() string
}
