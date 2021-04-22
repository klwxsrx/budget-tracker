package storedevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	domainEvent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type Deserializer interface {
	Deserialize(event *storedevent.StoredEvent) (domainEvent.Event, error)
}

type deserializer struct{}

func (d *deserializer) Deserialize(event *storedevent.StoredEvent) (domainEvent.Event, error) {
	switch event.Type {
	case domain.EventTypeAccountCreated:
		return d.deserializeAccountCreatedEvent(event.EventData)
	case domain.EventTypeAccountTitleChanged:
		return d.deserializeAccountTitleChangedEvent(event.EventData)
	case domain.EventTypeAccountDeleted:
		return d.deserializeAccountDeletedEvent(event.EventData)
	default:
		return nil, errors.New(fmt.Sprintf("unknown event, %v", event.Type))
	}
}

func (d *deserializer) deserializeAccountCreatedEvent(eventJson []byte) (domainEvent.Event, error) {
	var event domain.AccountCreatedEvent
	err := json.Unmarshal(eventJson, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (d *deserializer) deserializeAccountTitleChangedEvent(eventJson []byte) (domainEvent.Event, error) {
	var event domain.AccountTitleChangedEvent
	err := json.Unmarshal(eventJson, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (d *deserializer) deserializeAccountDeletedEvent(eventJson []byte) (domainEvent.Event, error) {
	var event domain.AccountDeletedEvent
	err := json.Unmarshal(eventJson, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func NewDeserializer() Deserializer {
	return &deserializer{}
}
