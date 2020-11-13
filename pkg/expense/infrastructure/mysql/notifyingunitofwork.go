package mysql

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
)

type notifyingUnitOfWork struct {
	unitOfWork command.UnitOfWork
	dispatcher event.StoredEventNotificationDispatcher
}

func (uw *notifyingUnitOfWork) Execute(f func(r command.DomainRegistry) error) error {
	err := uw.unitOfWork.Execute(f)
	if err != nil {
		uw.dispatcher.Dispatch()
	}
	return err
}

func (uw *notifyingUnitOfWork) Critical(lock string, f func(r command.DomainRegistry) error) error {
	err := uw.unitOfWork.Critical(lock, f)
	if err != nil {
		uw.dispatcher.Dispatch()
	}
	return err
}

func NewNotifyingUnitOfWork(unitOfWork command.UnitOfWork, dispatcher event.StoredEventNotificationDispatcher) command.UnitOfWork {
	return &notifyingUnitOfWork{unitOfWork, dispatcher}
}
