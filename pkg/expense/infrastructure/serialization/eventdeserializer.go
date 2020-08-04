package serialization

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	commonDomain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

type EventDeserializer interface {
	Deserialize(event *event.StoredEvent) (commonDomain.Event, error)
}

type eventDeserializer struct{}

func (ed *eventDeserializer) Deserialize(event *event.StoredEvent) (commonDomain.Event, error) {
	switch event.Type {
	case domain.AccountCreatedEvent:
		return ed.deserializeAccountCreatedEvent(event.EventData)
	case domain.AccountTitleChangedEvent:
		return ed.deserializeAccountTitleChangedEvent(event.EventData)
	case domain.AccountDeletedEvent:
		return ed.deserializeAccountDeletedEvent(event.EventData)
	default:
		return nil, errors.New(fmt.Sprintf("unknown event, %v", event.Type))
	}
}

func (ed *eventDeserializer) deserializeAccountCreatedEvent(eventJson []byte) (commonDomain.Event, error) {
	var domainEvent domain.AccountCreated
	err := json.Unmarshal(eventJson, &domainEvent)
	if err != nil {
		return nil, err
	}
	return &domainEvent, nil
}

func (ed *eventDeserializer) deserializeAccountTitleChangedEvent(eventJson []byte) (commonDomain.Event, error) {
	var domainEvent domain.AccountTitleChanged
	err := json.Unmarshal(eventJson, &domainEvent)
	if err != nil {
		return nil, err
	}
	return &domainEvent, nil
}

func (ed *eventDeserializer) deserializeAccountDeletedEvent(eventJson []byte) (commonDomain.Event, error) {
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
