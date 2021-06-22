package storedevent

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commonEvent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

var errorUnknownAccountEventType = errors.New("unknown account event")

type baseEventJson struct {
	AggregateID   uuid.UUID `json:"id"`
	AggregateName string    `json:"name"`
	Type          string    `json:"type"`
}

type accountListCreatedJson struct {
	baseEventJson
}

type accountCreatedJson struct {
	baseEventJson
	AccountID      uuid.UUID `json:"acc_id"`
	Title          string    `json:"title"`
	Currency       string    `json:"curr"`
	InitialBalance int       `json:"balance"`
}

type accountReorderedJson struct {
	baseEventJson
	AccountID uuid.UUID `json:"acc_id"`
	Position  int       `json:"pos"`
}

type accountRenamedJson struct {
	baseEventJson
	AccountID uuid.UUID `json:"acc_id"`
	Title     string    `json:"title"`
}

type accountActivatedJson struct {
	baseEventJson
	AccountID uuid.UUID `json:"acc_id"`
}

type accountCancelledJson struct {
	baseEventJson
	AccountID uuid.UUID `json:"acc_id"`
}

type accountDeletedJson struct {
	baseEventJson
	AccountID uuid.UUID `json:"acc_id"`
}

type eventSerializer struct{}

func (s *eventSerializer) Serialize(event commonEvent.Event) (string, error) {
	var err error
	var eventJson interface{}
	switch event.Type() {
	case domain.EventTypeAccountListCreated:
		eventJson, err = s.createAccountListCreatedJson(event)
	case domain.EventTypeAccountCreated:
		eventJson, err = s.createAccountCreatedJson(event)
	case domain.EventTypeAccountReordered:
		eventJson, err = s.createAccountReorderedJson(event)
	case domain.EventTypeAccountRenamed:
		eventJson, err = s.createAccountRenamedJson(event)
	case domain.EventTypeAccountActivated:
		eventJson, err = s.createAccountActivatedJson(event)
	case domain.EventTypeAccountCancelled:
		eventJson, err = s.createAccountCancelledJson(event)
	case domain.EventTypeAccountDeleted:
		eventJson, err = s.createAccountDeletedJson(event)
	default:
		return "", errors.New(fmt.Sprintf("unkown event type %v", event.Type()))
	}

	result, err := json.Marshal(eventJson)
	if err != nil {
		return "", fmt.Errorf("can't serialize event - %s: %v", event, err)
	}
	return string(result), nil
}

func (s *eventSerializer) createAccountListCreatedJson(event commonEvent.Event) (*accountListCreatedJson, error) {
	_, ok := event.(*domain.AccountListCreatedEvent)
	if !ok {
		return nil, errorUnknownAccountEventType
	}
	return &accountListCreatedJson{baseEventJson{
		AggregateID:   event.AggregateID().UUID,
		AggregateName: event.AggregateName(),
		Type:          event.Type(),
	}}, nil
}

func (s *eventSerializer) createAccountCreatedJson(event commonEvent.Event) (*accountCreatedJson, error) {
	createdEvent, ok := event.(*domain.AccountCreatedEvent)
	if !ok {
		return nil, errorUnknownAccountEventType
	}
	return &accountCreatedJson{
		baseEventJson: baseEventJson{
			AggregateID:   createdEvent.AggregateId,
			AggregateName: createdEvent.Name,
			Type:          createdEvent.EventType,
		},
		AccountID:      createdEvent.AccountID.UUID,
		Title:          createdEvent.Title,
		Currency:       string(createdEvent.Currency),
		InitialBalance: createdEvent.InitialBalance,
	}, nil
}

func (s *eventSerializer) createAccountReorderedJson(event commonEvent.Event) (*accountReorderedJson, error) {
	reorderedEvent, ok := event.(*domain.AccountReorderedEvent)
	if !ok {
		return nil, errorUnknownAccountEventType
	}
	return &accountReorderedJson{
		baseEventJson: baseEventJson{
			AggregateID:   reorderedEvent.AggregateId,
			AggregateName: reorderedEvent.Name,
			Type:          reorderedEvent.EventType,
		},
		AccountID: reorderedEvent.AccountID.UUID,
		Position:  reorderedEvent.Position,
	}, nil
}

func (s *eventSerializer) createAccountRenamedJson(event commonEvent.Event) (*accountRenamedJson, error) {
	renamedEvent, ok := event.(*domain.AccountRenamedEvent)
	if !ok {
		return nil, errorUnknownAccountEventType
	}
	return &accountRenamedJson{
		baseEventJson: baseEventJson{
			AggregateID:   renamedEvent.AggregateId,
			AggregateName: renamedEvent.Name,
			Type:          renamedEvent.EventType,
		},
		AccountID: renamedEvent.AccountID.UUID,
		Title:     renamedEvent.Title,
	}, nil
}

func (s *eventSerializer) createAccountActivatedJson(event commonEvent.Event) (*accountActivatedJson, error) {
	activatedEvent, ok := event.(*domain.AccountActivatedEvent)
	if !ok {
		return nil, errorUnknownAccountEventType
	}
	return &accountActivatedJson{
		baseEventJson: baseEventJson{
			AggregateID:   activatedEvent.AggregateId,
			AggregateName: activatedEvent.Name,
			Type:          activatedEvent.EventType,
		},
		AccountID: activatedEvent.AccountID.UUID,
	}, nil
}

func (s *eventSerializer) createAccountCancelledJson(event commonEvent.Event) (*accountCancelledJson, error) {
	cancelledEvent, ok := event.(*domain.AccountCancelledEvent)
	if !ok {
		return nil, errorUnknownAccountEventType
	}
	return &accountCancelledJson{
		baseEventJson: baseEventJson{
			AggregateID:   cancelledEvent.AggregateId,
			AggregateName: cancelledEvent.Name,
			Type:          cancelledEvent.EventType,
		},
		AccountID: cancelledEvent.AccountID.UUID,
	}, nil
}

func (s *eventSerializer) createAccountDeletedJson(event commonEvent.Event) (*accountDeletedJson, error) {
	deletedEvent, ok := event.(*domain.AccountDeletedEvent)
	if !ok {
		return nil, errorUnknownAccountEventType
	}
	return &accountDeletedJson{
		baseEventJson: baseEventJson{
			AggregateID:   deletedEvent.AggregateId,
			AggregateName: deletedEvent.Name,
			Type:          deletedEvent.EventType,
		},
		AccountID: deletedEvent.AccountID.UUID,
	}, nil
}

func NewSerializer() storedevent.Serializer {
	return &eventSerializer{}
}
