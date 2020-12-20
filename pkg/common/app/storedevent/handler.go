package storedevent

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/persistence"
	"sync/atomic"
	"time"
)

const dispatchPeriod = time.Second

type Handler interface {
	HandleUnsentStoredEvents()
	Stop()
}

type handler struct {
	busHandler     *UnsentEventBusHandler
	logger         logger.Logger
	needDispatch   int32
	stopChan       chan struct{}
	stopErrorsChan chan struct{}
}

func (d *handler) HandleUnsentStoredEvents() {
	atomic.StoreInt32(&d.needDispatch, 1)
}

func (d *handler) Stop() {
	d.stopChan <- struct{}{}
}

func (d *handler) start() {
	errorsChan := make(chan error)
	go func() {
		for {
			select {
			case err := <-errorsChan:
				d.logger.WithError(err).Error("failed to handle unsent events")
			case <-d.stopErrorsChan:
				return
			}
		}
	}()

	ticker := time.NewTicker(dispatchPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				needDispatch := atomic.SwapInt32(&d.needDispatch, 0)
				if needDispatch == 1 {
					err := d.busHandler.ProcessUnsentEvents()
					if err != nil {
						atomic.StoreInt32(&d.needDispatch, 1)
						errorsChan <- err
					}
				}
			case <-d.stopChan:
				ticker.Stop()
				d.stopErrorsChan <- struct{}{}
				return
			}
		}
	}()
}

func NewHandler(unsentEventProvider UnsentEventProvider, eventBus Bus, sync persistence.Synchronization, logger logger.Logger) Handler {
	busHandler := &UnsentEventBusHandler{unsentEventProvider, eventBus, sync}
	dispatcher := &handler{busHandler, logger, 1, make(chan struct{}), make(chan struct{})}
	go dispatcher.start()
	return dispatcher
}
