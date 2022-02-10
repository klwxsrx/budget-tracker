package realtime

type Client interface {
	PublishMessage(channel string, data []byte) error
}
