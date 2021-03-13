package messaging

import "github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"

type StoredEventSerializer interface {
	Serialize(event *storedevent.StoredEvent) ([]byte, error)
}
