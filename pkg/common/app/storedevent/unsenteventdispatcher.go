package storedevent

import (
	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"

	"sync"
)

type UnsentEventDispatcher interface {
	Dispatch()
	Start()
	Stop()
}

type unsentEventDispatcher struct {
	handler             *UnsentEventHandler
	logger              log.Logger
	dispatchRequestChan chan struct{}
	stopChan            chan struct{}
	sync.Mutex
	isStarted bool
}

func (d *unsentEventDispatcher) Dispatch() {
	select {
	case d.dispatchRequestChan <- struct{}{}:
	default:
	}
}

func (d *unsentEventDispatcher) Start() {
	d.Lock()
	defer d.Unlock()

	if d.isStarted {
		return
	}

	go d.run()
	d.Dispatch()
	d.isStarted = true
}

func (d *unsentEventDispatcher) Stop() {
	d.Lock()
	defer d.Unlock()

	if !d.isStarted {
		return
	}

	d.stopChan <- struct{}{}
	d.isStarted = false
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
	logger log.Logger,
) UnsentEventDispatcher {
	dispatcher := &unsentEventDispatcher{
		unsentEventHandler,
		logger,
		make(chan struct{}, 1),
		make(chan struct{}),
		sync.Mutex{},
		false,
	}
	return dispatcher
}
