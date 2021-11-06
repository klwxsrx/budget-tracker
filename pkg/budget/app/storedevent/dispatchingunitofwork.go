package storedevent

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/service"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

type dispatchingUnitOfWork struct {
	unitOfWork service.UnitOfWork
	dispatcher storedevent.UnsentEventDispatcher
}

func (uw *dispatchingUnitOfWork) Execute(f func(r service.DomainRegistry) error) error {
	err := uw.unitOfWork.Execute(f)
	if err == nil {
		uw.dispatcher.Dispatch()
	}
	return err
}

func (uw *dispatchingUnitOfWork) Critical(lock string, f func(r service.DomainRegistry) error) error {
	err := uw.unitOfWork.Critical(lock, f)
	if err == nil {
		uw.dispatcher.Dispatch()
	}
	return err
}

func NewDispatchingUnitOfWork(unitOfWork service.UnitOfWork, dispatcher storedevent.UnsentEventDispatcher) service.UnitOfWork {
	return &dispatchingUnitOfWork{unitOfWork, dispatcher}
}
