package messaging

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

type EventDeserializer interface {
	Deserialize(event *messaging.Event) (event.Event, error)
}

type eventDeserializer struct{}

func (ed *eventDeserializer) Deserialize(event *messaging.Event) (event.Event, error) {
	switch event.Type {
	case domain.AccountCreatedEvent:
		return ed.deserializeAccountCreatedEvent(event.Payload)
	case domain.AccountTitleChangedEvent:
		return ed.deserializeAccountTitleChangedEvent(event.Payload)
	case domain.AccountDeletedEvent:
		return ed.deserializeAccountDeletedEvent(event.Payload)
	default:
		return nil, errors.New(fmt.Sprintf("unknown event, %v", event.Type))
	}
}

func (ed *eventDeserializer) deserializeAccountCreatedEvent(eventJson []byte) (event.Event, error) {
	var domainEvent domain.AccountCreated
	err := json.Unmarshal(eventJson, &domainEvent)
	if err != nil {
		return nil, err
	}
	return &domainEvent, nil
}

func (ed *eventDeserializer) deserializeAccountTitleChangedEvent(eventJson []byte) (event.Event, error) {
	var domainEvent domain.AccountTitleChanged
	err := json.Unmarshal(eventJson, &domainEvent)
	if err != nil {
		return nil, err
	}
	return &domainEvent, nil
}

func (ed *eventDeserializer) deserializeAccountDeletedEvent(eventJson []byte) (event.Event, error) {
	var domainEvent domain.AccountDeleted
	err := json.Unmarshal(eventJson, &domainEvent)
	if err != nil {
		return nil, err
	}
	return &domainEvent, nil
}

func NewEventDeserializer() EventDeserializer {
	return &eventDeserializer{}
}
