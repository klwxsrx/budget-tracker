package serialization

import "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"

type EventSerializer interface {
	Serialize(event event.Event) (string, error)
}
