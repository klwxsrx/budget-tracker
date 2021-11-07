package storedevent

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
)

type UnsentEventDispatcher interface {
	Dispatch()
	Stop()
}

type unsentEventDispatcher struct {
	handler             *UnsentEventHandler
	logger              logger.Logger
	dispatchRequestChan chan struct{}
	stopChan            chan struct{}
}

func (d *unsentEventDispatcher) Dispatch() {
	select {
	case d.dispatchRequestChan <- struct{}{}:
	default:
	}
}

func (d *unsentEventDispatcher) Stop() {
	d.stopChan <- struct{}{}
}

func (d *unsentEventDispatcher) run() {
	go func() {
		for {
			select {
			case <-d.dispatchRequestChan:
				err := d.handler.ProcessUnsentEvents()
				if err != nil {
					d.logger.WithError(err).Error("failed to process unsent events")
					d.Dispatch()
				}
			case <-d.stopChan:
				return
			}
		}
	}()
}

func NewUnsentEventDispatcher(
	unsentEventHandler *UnsentEventHandler,
	loggerImpl logger.Logger,
) UnsentEventDispatcher {
	dispatcher := &unsentEventDispatcher{
		unsentEventHandler,
		loggerImpl,
		make(chan struct{}, 1),
		make(chan struct{}),
	}
	dispatcher.run()
	dispatcher.Dispatch()
	return dispatcher
}
