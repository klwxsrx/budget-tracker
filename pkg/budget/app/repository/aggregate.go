package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
	commondomain "github.com/klwxsrx/budget-tracker/pkg/common/domain"
)

var errAggregateNotFound = errors.New("aggregate not found")

type AggregateRepository struct {
	store        storedevent.Store
	deserializer messaging.DomainEventDeserializer
}

func (repo *AggregateRepository) storeChanges(aggregate domain.Aggregate) error {
	for _, event := range aggregate.GetChanges() {
		_, err := repo.store.Append(event)
		if err != nil {
			return fmt.Errorf("can't update aggregate: %w", err)
		}
	}
	return nil
}

func (repo *AggregateRepository) loadChanges(id commondomain.AggregateID, state domain.AggregateState) error {
	storedEvents, err := repo.store.GetByAggregate(id, state.AggregateName(), storedevent.ID{UUID: uuid.Nil})
	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}
	if len(storedEvents) == 0 {
		return fmt.Errorf("%w, id: %v", errAggregateNotFound, id)
	}
	for _, storedEvent := range storedEvents {
		event, err := repo.deserializer.Deserialize(storedEvent.Type, storedEvent.EventData)
		if err != nil {
			return fmt.Errorf("failed to deserialize events: %w", err)
		}
		err = state.Apply(event)
		if err != nil {
			return fmt.Errorf("failed to apply event: %w", err)
		}
	}
	return nil
}

func NewAggregateRepository(
	store storedevent.Store,
	deserializer messaging.DomainEventDeserializer,
) *AggregateRepository {
	return &AggregateRepository{store, deserializer}
}
