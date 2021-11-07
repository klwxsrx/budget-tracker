package event

import (
	"errors"

	"github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

var ErrUnexpectedEventType = errors.New("unexpected error type")

type DomainEventHandler interface {
	Handle(e event.Event) error
}
