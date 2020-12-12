package event

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
)

type storedEventHandingUnitOfWork struct {
	unitOfWork         command.UnitOfWork
	storedEventHandler event.StoredEventHandler
}

func (uw *storedEventHandingUnitOfWork) Execute(f func(r command.DomainRegistry) error) error {
	err := uw.unitOfWork.Execute(f)
	if err != nil {
		uw.storedEventHandler.HandleStoredEvents()
	}
	return err
}

func (uw *storedEventHandingUnitOfWork) Critical(lock string, f func(r command.DomainRegistry) error) error {
	err := uw.unitOfWork.Critical(lock, f)
	if err != nil {
		uw.storedEventHandler.HandleStoredEvents()
	}
	return err
}

func NewStoredEventHandlingUnitOfWork(unitOfWork command.UnitOfWork, dispatcher event.StoredEventHandler) command.UnitOfWork {
	return &storedEventHandingUnitOfWork{unitOfWork, dispatcher}
}
