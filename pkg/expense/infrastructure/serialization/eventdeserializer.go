package serialization

import (
	"encoding/json"
	"errors"
	"fmt"
	eventApp "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	eventDomain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

type EventDeserializer interface {
	Deserialize(event *eventApp.StoredEvent) (eventDomain.Event, error)
}

type eventDeserializer struct{}

func (ed *eventDeserializer) Deserialize(event *eventApp.StoredEvent) (eventDomain.Event, error) {
	switch event.Type {
	case domain.AccountCreatedEventType:
		return ed.deserializeAccountCreatedEvent(event.EventData)
	case domain.AccountTitleChangedEventType:
		return ed.deserializeAccountTitleChangedEvent(event.EventData)
	case domain.AccountDeletedEventType:
		return ed.deserializeAccountDeletedEvent(event.EventData)
	default:
		return nil, errors.New(fmt.Sprintf("unknown event, %v", event.Type))
	}
}

func (ed *eventDeserializer) deserializeAccountCreatedEvent(eventJson []byte) (eventDomain.Event, error) {
	var domainEvent domain.AccountCreatedEvent
	err := json.Unmarshal(eventJson, &domainEvent)
	if err != nil {
		return nil, err
	}
	return &domainEvent, nil
}

func (ed *eventDeserializer) deserializeAccountTitleChangedEvent(eventJson []byte) (eventDomain.Event, error) {
	var domainEvent domain.AccountTitleChangedEvent
	err := json.Unmarshal(eventJson, &domainEvent)
	if err != nil {
		return nil, err
	}
	return &domainEvent, nil
}

func (ed *eventDeserializer) deserializeAccountDeletedEvent(eventJson []byte) (eventDomain.Event, error) {
	var domainEvent domain.AccountDeletedEvent
	err := json.Unmarshal(eventJson, &domainEvent)
	if err != nil {
		return nil, err
	}
	return &domainEvent, nil
}

func NewEventDeserializer() EventDeserializer {
	return &eventDeserializer{}
}
