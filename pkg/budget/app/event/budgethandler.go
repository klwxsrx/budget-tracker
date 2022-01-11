package event

import (
	"errors"

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
	err := handler.bus.Publish(cmd)
	if errors.Is(err, domain.ErrAccountListAlreadyExists) {
		return nil
	}
	return err
}

func NewBudgetCreatedEventHandler(bus command.Bus) commonappevent.DomainEventHandler {
	return &budgetCreatedEventHandler{bus}
}
