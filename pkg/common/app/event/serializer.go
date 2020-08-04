package event

import "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"

type Serializer interface {
	Serialize(event event.Event) (string, error)
}
