package event

import (
	budgetappcommand "github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/event"
	commondomainevent "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type budgetCreatedEventHandler struct {
	bus command.Bus
}

func (handler *budgetCreatedEventHandler) Handle(evt commondomainevent.Event) error {
	budgetCreated, ok := evt.(*domain.BudgetCreatedEvent)
	if !ok {
		return event.ErrUnexpectedEventType
	}

	cmd := budgetappcommand.NewAccountCreateList(budgetCreated.EventAggregateID)
	result, err := handler.bus.Publish(cmd)
	switch result {
	case command.ResultSuccess, command.ResultDuplicateConflict:
		return nil
	case command.ResultInvalidArgument, command.ResultNotFound, command.ResultUnknownError:
	}
	return err
}

func NewBudgetCreatedEventHandler(bus command.Bus) event.DomainEventHandler {
	return &budgetCreatedEventHandler{bus}
}
