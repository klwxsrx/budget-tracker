package event

import (
	"errors"

	"github.com/klwxsrx/budget-tracker/pkg/common/domain"
)

var ErrUnexpectedEventType = errors.New("unexpected error type")

type DomainEventHandler interface {
	Handle(event domain.Event) error
}
