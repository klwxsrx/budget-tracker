package storedevent

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

type handingUnitOfWork struct {
	unitOfWork service.UnitOfWork
	handler    storedevent.Handler
}

func (uw *handingUnitOfWork) Execute(f func(r service.DomainRegistry) error) error {
	err := uw.unitOfWork.Execute(f)
	if err == nil {
		uw.handler.HandleUnsentStoredEvents()
	}
	return err
}

func (uw *handingUnitOfWork) Critical(lock string, f func(r service.DomainRegistry) error) error {
	err := uw.unitOfWork.Critical(lock, f)
	if err == nil {
		uw.handler.HandleUnsentStoredEvents()
	}
	return err
}

func NewHandlingUnitOfWork(unitOfWork service.UnitOfWork, dispatcher storedevent.Handler) service.UnitOfWork {
	return &handingUnitOfWork{unitOfWork, dispatcher}
}
