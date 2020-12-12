package event

import domain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"

type Serializer interface {
	Serialize(event domain.Event) (string, error)
}
