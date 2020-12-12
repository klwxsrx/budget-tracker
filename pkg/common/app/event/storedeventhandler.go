package event

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"sync/atomic"
	"time"
)

const dispatchPeriod = time.Second

type StoredEventHandler interface {
	HandleStoredEvents()
	Stop()
}

type storedEventHandler struct {
	busHandler   *BusHandler
	logger       logger.Logger
	needDispatch int32
	stopChan     chan struct{}
}

func (d *storedEventHandler) HandleStoredEvents() {
	atomic.StoreInt32(&d.needDispatch, 1)
}

func (d *storedEventHandler) Stop() {
	d.stopChan <- struct{}{}
}

func (d *storedEventHandler) start() {
	errorsChan := make(chan error)
	go func() {
		for {
			select {
			case err := <-errorsChan:
				d.logger.WithError(err).Error("failed to handle event events")
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
						errorsChan <- err // TODO: test case without connection
					}
				}
			case <-d.stopChan:
				return
			}
		}
	}()
}

func NewStoredEventHandler(unsentEventProvider UnsentEventProvider, eventBus Bus, logger logger.Logger) StoredEventHandler {
	busHandler := &BusHandler{unsentEventProvider, eventBus}
	dispatcher := &storedEventHandler{busHandler, logger, 1, make(chan struct{})}
	go dispatcher.start()
	return dispatcher
}
