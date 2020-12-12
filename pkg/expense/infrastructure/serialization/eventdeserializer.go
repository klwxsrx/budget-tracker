package serialization

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
	domainEvent "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

type EventDeserializer interface {
	Deserialize(event *storedevent.StoredEvent) (domainEvent.Event, error)
}

type eventDeserializer struct{}

func (ed *eventDeserializer) Deserialize(event *storedevent.StoredEvent) (domainEvent.Event, error) {
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

func (ed *eventDeserializer) deserializeAccountCreatedEvent(eventJson []byte) (domainEvent.Event, error) {
	var event domain.AccountCreatedEvent
	err := json.Unmarshal(eventJson, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (ed *eventDeserializer) deserializeAccountTitleChangedEvent(eventJson []byte) (domainEvent.Event, error) {
	var event domain.AccountTitleChangedEvent
	err := json.Unmarshal(eventJson, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (ed *eventDeserializer) deserializeAccountDeletedEvent(eventJson []byte) (domainEvent.Event, error) {
	var event domain.AccountDeletedEvent
	err := json.Unmarshal(eventJson, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func NewEventDeserializer() EventDeserializer {
	return &eventDeserializer{}
}
