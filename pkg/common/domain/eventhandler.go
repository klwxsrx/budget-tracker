package domain

import (
	"errors"
)

var ErrUnexpectedEventType = errors.New("unexpected error type")

type EventHandler interface {
	Handle(event Event) error
}
