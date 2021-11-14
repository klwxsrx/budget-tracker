package event

import (
	budgetappcommand "github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	commonappevent "github.com/klwxsrx/budget-tracker/pkg/common/app/event"
	commondomain "github.com/klwxsrx/budget-tracker/pkg/common/domain"
)

type budgetCreatedEventHandler struct {
	bus command.Bus
}

func (handler *budgetCreatedEventHandler) Handle(event commondomain.Event) error {
	budgetCreated, ok := event.(*domain.BudgetCreatedEvent)
	if !ok {
		return commonappevent.ErrUnexpectedEventType
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

func NewBudgetCreatedEventHandler(bus command.Bus) commonappevent.DomainEventHandler {
	return &budgetCreatedEventHandler{bus}
}
