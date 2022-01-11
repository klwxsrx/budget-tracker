package storedevent

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commondomain "github.com/klwxsrx/budget-tracker/pkg/common/domain"
)

var errUnknownAccountEventType = errors.New("unknown account event")

type baseEventJSON struct {
	AggregateID   uuid.UUID `json:"id"`
	AggregateName string    `json:"name"`
	Type          string    `json:"type"`
}

type budgetCreatedJSON struct {
	baseEventJSON
	Title    string `json:"title"`
	Currency string `json:"currency"`
}

type accountListCreatedJSON struct {
	baseEventJSON
}

type accountCreatedJSON struct {
	baseEventJSON
	AccountID      uuid.UUID `json:"acc_id"`
	Title          string    `json:"title"`
	InitialBalance int       `json:"balance"`
}

type accountReorderedJSON struct {
	baseEventJSON
	AccountID uuid.UUID `json:"acc_id"`
	Position  int       `json:"pos"`
}

type accountRenamedJSON struct {
	baseEventJSON
	AccountID uuid.UUID `json:"acc_id"`
	Title     string    `json:"title"`
}

type accountActivatedJSON struct {
	baseEventJSON
	AccountID uuid.UUID `json:"acc_id"`
}

type accountCancelledJSON struct {
	baseEventJSON
	AccountID uuid.UUID `json:"acc_id"`
}

type accountDeletedJSON struct {
	baseEventJSON
	AccountID uuid.UUID `json:"acc_id"`
}

type eventSerializer struct{}

func (s *eventSerializer) Serialize(event commondomain.Event) (string, error) {
	var err error
	var eventJSON interface{}
	switch event.Type() {
	case domain.EventTypeBudgetCreated:
		eventJSON, err = s.createBudgetCreatedJSON(event)
	case domain.EventTypeAccountListCreated:
		eventJSON, err = s.createAccountListCreatedJSON(event)
	case domain.EventTypeAccountCreated:
		eventJSON, err = s.createAccountCreatedJSON(event)
	case domain.EventTypeAccountReordered:
		eventJSON, err = s.createAccountReorderedJSON(event)
	case domain.EventTypeAccountRenamed:
		eventJSON, err = s.createAccountRenamedJSON(event)
	case domain.EventTypeAccountActivated:
		eventJSON, err = s.createAccountActivatedJSON(event)
	case domain.EventTypeAccountCancelled:
		eventJSON, err = s.createAccountCancelledJSON(event)
	case domain.EventTypeAccountDeleted:
		eventJSON, err = s.createAccountDeletedJSON(event)
	default:
		return "", fmt.Errorf("unknown event type %s", event.Type())
	}
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(eventJSON)
	if err != nil {
		return "", fmt.Errorf("can't serialize event - %s: %w", event, err)
	}
	return string(result), nil
}

func (s *eventSerializer) createBudgetCreatedJSON(event commondomain.Event) (*budgetCreatedJSON, error) {
	createdEvent, ok := event.(*domain.BudgetCreatedEvent)
	if !ok {
		return nil, errUnknownAccountEventType
	}
	return &budgetCreatedJSON{
		baseEventJSON: baseEventJSON{
			AggregateID:   event.AggregateID().UUID,
			AggregateName: event.AggregateName(),
			Type:          event.Type(),
		},
		Title:    createdEvent.Title,
		Currency: createdEvent.Currency,
	}, nil
}

func (s *eventSerializer) createAccountListCreatedJSON(event commondomain.Event) (*accountListCreatedJSON, error) {
	_, ok := event.(*domain.AccountListCreatedEvent)
	if !ok {
		return nil, errUnknownAccountEventType
	}
	return &accountListCreatedJSON{baseEventJSON{
		AggregateID:   event.AggregateID().UUID,
		AggregateName: event.AggregateName(),
		Type:          event.Type(),
	}}, nil
}

func (s *eventSerializer) createAccountCreatedJSON(event commondomain.Event) (*accountCreatedJSON, error) {
	createdEvent, ok := event.(*domain.AccountCreatedEvent)
	if !ok {
		return nil, errUnknownAccountEventType
	}
	return &accountCreatedJSON{
		baseEventJSON: baseEventJSON{
			AggregateID:   createdEvent.EventAggregateID,
			AggregateName: createdEvent.EventAggregateName,
			Type:          createdEvent.EventType,
		},
		AccountID:      createdEvent.AccountID.UUID,
		Title:          createdEvent.Title,
		InitialBalance: createdEvent.InitialBalance,
	}, nil
}

func (s *eventSerializer) createAccountReorderedJSON(event commondomain.Event) (*accountReorderedJSON, error) {
	reorderedEvent, ok := event.(*domain.AccountReorderedEvent)
	if !ok {
		return nil, errUnknownAccountEventType
	}
	return &accountReorderedJSON{
		baseEventJSON: baseEventJSON{
			AggregateID:   reorderedEvent.EventAggregateID,
			AggregateName: reorderedEvent.EventAggregateName,
			Type:          reorderedEvent.EventType,
		},
		AccountID: reorderedEvent.AccountID.UUID,
		Position:  reorderedEvent.Position,
	}, nil
}

func (s *eventSerializer) createAccountRenamedJSON(event commondomain.Event) (*accountRenamedJSON, error) {
	renamedEvent, ok := event.(*domain.AccountRenamedEvent)
	if !ok {
		return nil, errUnknownAccountEventType
	}
	return &accountRenamedJSON{
		baseEventJSON: baseEventJSON{
			AggregateID:   renamedEvent.EventAggregateID,
			AggregateName: renamedEvent.EventAggregateName,
			Type:          renamedEvent.EventType,
		},
		AccountID: renamedEvent.AccountID.UUID,
		Title:     renamedEvent.Title,
	}, nil
}

func (s *eventSerializer) createAccountActivatedJSON(event commondomain.Event) (*accountActivatedJSON, error) {
	activatedEvent, ok := event.(*domain.AccountActivatedEvent)
	if !ok {
		return nil, errUnknownAccountEventType
	}
	return &accountActivatedJSON{
		baseEventJSON: baseEventJSON{
			AggregateID:   activatedEvent.EventAggregateID,
			AggregateName: activatedEvent.EventAggregateName,
			Type:          activatedEvent.EventType,
		},
		AccountID: activatedEvent.AccountID.UUID,
	}, nil
}

func (s *eventSerializer) createAccountCancelledJSON(event commondomain.Event) (*accountCancelledJSON, error) {
	cancelledEvent, ok := event.(*domain.AccountCancelledEvent)
	if !ok {
		return nil, errUnknownAccountEventType
	}
	return &accountCancelledJSON{
		baseEventJSON: baseEventJSON{
			AggregateID:   cancelledEvent.EventAggregateID,
			AggregateName: cancelledEvent.EventAggregateName,
			Type:          cancelledEvent.EventType,
		},
		AccountID: cancelledEvent.AccountID.UUID,
	}, nil
}

func (s *eventSerializer) createAccountDeletedJSON(event commondomain.Event) (*accountDeletedJSON, error) {
	deletedEvent, ok := event.(*domain.AccountDeletedEvent)
	if !ok {
		return nil, errUnknownAccountEventType
	}
	return &accountDeletedJSON{
		baseEventJSON: baseEventJSON{
			AggregateID:   deletedEvent.EventAggregateID,
			AggregateName: deletedEvent.EventAggregateName,
			Type:          deletedEvent.EventType,
		},
		AccountID: deletedEvent.AccountID.UUID,
	}, nil
}

func NewSerializer() storedevent.Serializer {
	return &eventSerializer{}
}
