package messaging

type StoredEventNotifier interface {
	NotifyOfCreatedEvents() error
}
