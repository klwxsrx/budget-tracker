package serialization

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	domainEvent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type EventDeserializer interface {
	Deserialize(event *storedevent.StoredEvent) (domainEvent.Event, error)
}

type eventDeserializer struct{}

func (ed *eventDeserializer) Deserialize(event *storedevent.StoredEvent) (domainEvent.Event, error) {
	switch event.Type {
	case domain.EventTypeAccountCreated:
		return ed.deserializeAccountCreatedEvent(event.EventData)
	case domain.EventTypeAccountTitleChanged:
		return ed.deserializeAccountTitleChangedEvent(event.EventData)
	case domain.EventTypeAccountDeleted:
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
