package storedevent

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
)

type handingUnitOfWork struct {
	unitOfWork command.UnitOfWork
	handler    storedevent.Handler
}

func (uw *handingUnitOfWork) Execute(f func(r command.DomainRegistry) error) error {
	err := uw.unitOfWork.Execute(f)
	if err != nil {
		uw.handler.HandleUnsentStoredEvents()
	}
	return err
}

func (uw *handingUnitOfWork) Critical(lock string, f func(r command.DomainRegistry) error) error {
	err := uw.unitOfWork.Critical(lock, f)
	if err != nil {
		uw.handler.HandleUnsentStoredEvents()
	}
	return err
}

func NewHandlingUnitOfWork(unitOfWork command.UnitOfWork, dispatcher storedevent.Handler) command.UnitOfWork {
	return &handingUnitOfWork{unitOfWork, dispatcher}
}
