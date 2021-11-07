package storedevent

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"

	"sync"
)

type UnsentEventDispatcher interface {
	Dispatch()
	Start()
	Stop()
}

type unsentEventDispatcher struct {
	handler             *UnsentEventHandler
	logger              logger.Logger
	dispatchRequestChan chan struct{}
	starter             sync.Once
	stopChan            chan struct{}
}

func (d *unsentEventDispatcher) Dispatch() {
	select {
	case d.dispatchRequestChan <- struct{}{}:
	default:
	}
}

func (d *unsentEventDispatcher) Start() {
	d.starter.Do(func() {
		go d.run()
		d.Dispatch()
	})
}

func (d *unsentEventDispatcher) Stop() { // TODO: stop before start
	d.stopChan <- struct{}{}
}

func (d *unsentEventDispatcher) run() {
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
}

func NewUnsentEventDispatcher(
	unsentEventHandler *UnsentEventHandler,
	loggerImpl logger.Logger,
) UnsentEventDispatcher {
	dispatcher := &unsentEventDispatcher{
		unsentEventHandler,
		loggerImpl,
		make(chan struct{}, 1),
		sync.Once{},
		make(chan struct{}),
	}
	return dispatcher
}
