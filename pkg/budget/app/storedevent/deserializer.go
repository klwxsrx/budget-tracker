package storedevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	domainEvent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type deserializer struct{}

func (d *deserializer) Deserialize(event *storedevent.StoredEvent) (domainEvent.Event, error) {
	switch event.Type {
	case domain.EventTypeAccountListCreated:
		return d.deserializeAccountListCreatedEvent(event.EventData)
	case domain.EventTypeAccountCreated:
		return d.deserializeAccountCreatedEvent(event.EventData)
	case domain.EventTypeAccountReordered:
		return d.deserializeAccountReorderedEvent(event.EventData)
	case domain.EventTypeAccountRenamed:
		return d.deserializeAccountRenamedEvent(event.EventData)
	case domain.EventTypeAccountActivated:
		return d.deserializeAccountActivatedEvent(event.EventData)
	case domain.EventTypeAccountCancelled:
		return d.deserializeAccountCancelledEvent(event.EventData)
	case domain.EventTypeAccountDeleted:
		return d.deserializeAccountDeletedEvent(event.EventData)
	default:
		return nil, errors.New(fmt.Sprintf("unknown event, %v", event.Type))
	}
}

func (d *deserializer) deserializeAccountListCreatedEvent(eventPayload []byte) (domainEvent.Event, error) {
	var event accountListCreatedJson
	err := json.Unmarshal(eventPayload, &event)
	if err != nil {
		return nil, err
	}
	return domain.NewEventAccountListCreated(domain.BudgetID{UUID: event.AggregateID}), nil
}

func (d *deserializer) deserializeAccountCreatedEvent(eventPayload []byte) (domainEvent.Event, error) {
	var event accountCreatedJson
	err := json.Unmarshal(eventPayload, &event)
	if err != nil {
		return nil, err
	}
	return domain.NewEventAccountCreated(
		domain.BudgetID{UUID: event.AggregateID},
		domain.AccountID{UUID: event.AccountID},
		event.Title,
		domain.Currency(event.Currency),
		event.InitialBalance,
	), nil
}

func (d *deserializer) deserializeAccountReorderedEvent(eventPayload []byte) (domainEvent.Event, error) {
	var event accountReorderedJson
	err := json.Unmarshal(eventPayload, &event)
	if err != nil {
		return nil, err
	}
	return domain.NewEventAccountReordered(
		domain.BudgetID{UUID: event.AggregateID},
		domain.AccountID{UUID: event.AccountID},
		event.Position,
	), nil
}

func (d *deserializer) deserializeAccountRenamedEvent(eventPayload []byte) (domainEvent.Event, error) {
	var event accountRenamedJson
	err := json.Unmarshal(eventPayload, &event)
	if err != nil {
		return nil, err
	}
	return domain.NewEventAccountRenamed(
		domain.BudgetID{UUID: event.AggregateID},
		domain.AccountID{UUID: event.AccountID},
		event.Title,
	), nil
}

func (d *deserializer) deserializeAccountActivatedEvent(eventPayload []byte) (domainEvent.Event, error) {
	var event accountActivatedJson
	err := json.Unmarshal(eventPayload, &event)
	if err != nil {
		return nil, err
	}
	return domain.NewEventAccountActivated(
		domain.BudgetID{UUID: event.AggregateID},
		domain.AccountID{UUID: event.AccountID},
	), nil
}

func (d *deserializer) deserializeAccountCancelledEvent(eventPayload []byte) (domainEvent.Event, error) {
	var event accountCancelledJson
	err := json.Unmarshal(eventPayload, &event)
	if err != nil {
		return nil, err
	}
	return domain.NewEventAccountCancelled(
		domain.BudgetID{UUID: event.AggregateID},
		domain.AccountID{UUID: event.AccountID},
	), nil
}

func (d *deserializer) deserializeAccountDeletedEvent(eventPayload []byte) (domainEvent.Event, error) {
	var event accountDeletedJson
	err := json.Unmarshal(eventPayload, &event)
	if err != nil {
		return nil, err
	}
	return domain.NewEventAccountDeleted(
		domain.BudgetID{UUID: event.AggregateID},
		domain.AccountID{UUID: event.AccountID},
	), nil
}

func NewDeserializer() storedevent.Deserializer {
	return &deserializer{}
}
