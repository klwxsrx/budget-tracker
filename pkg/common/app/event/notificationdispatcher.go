package event

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/app/event/messaging"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"sync/atomic"
	"time"
)

const dispatchPeriod = time.Second

type StoredEventNotificationDispatcher interface {
	Dispatch()
	Start()
	Stop()
}

type notificationDispatcher struct {
	notifier     messaging.StoredEventNotifier
	logger       logger.Logger
	needDispatch int32
	stopChan     chan struct{}
}

func (d *notificationDispatcher) Dispatch() {
	atomic.StoreInt32(&d.needDispatch, 1)
}

func (d *notificationDispatcher) Start() {
	errorsChan := make(chan error)
	go func() {
		for {
			select {
			case err := <-errorsChan:
				d.logger.With(logger.Fields{"error": err}).Error("failed to dispatch event notification")
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
					err := d.notifier.NotifyOfCreatedEvents()
					if err != nil {
						atomic.StoreInt32(&d.needDispatch, 1)
						errorsChan <- err
					}
				}
			case <-d.stopChan:
				return
			}
		}
	}()
}

func (d *notificationDispatcher) Stop() {
	d.stopChan <- struct{}{}
}

func NewStoredEventNotificationDispatcher(notifier messaging.StoredEventNotifier, logger logger.Logger) StoredEventNotificationDispatcher {
	return &notificationDispatcher{notifier, logger, 1, make(chan struct{})}
}