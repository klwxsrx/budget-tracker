package storedevent

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
)

const dispatchPeriod = time.Second

type UnsentEventDispatcher interface {
	Dispatch()
}

type unsentEventDispatcher struct {
	handler      *UnsentEventHandler
	logger       logger.Logger
	needDispatch int32
}

func (d *unsentEventDispatcher) Dispatch() {
	atomic.StoreInt32(&d.needDispatch, 1)
}

func (d *unsentEventDispatcher) run(ctx context.Context) {
	ticker := time.NewTicker(dispatchPeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				needDispatch := atomic.SwapInt32(&d.needDispatch, 0)
				if needDispatch == 1 {
					err := d.handler.ProcessUnsentEvents()
					if err != nil {
						atomic.StoreInt32(&d.needDispatch, 1)
						d.logger.WithError(err).Error("failed to handle unsent events")
					}
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func NewUnsentEventDispatcher(
	ctx context.Context,
	unsentEventHandler *UnsentEventHandler,
	loggerImpl logger.Logger,
) UnsentEventDispatcher {
	dispatcher := &unsentEventDispatcher{unsentEventHandler, loggerImpl, 1}
	dispatcher.run(ctx)
	return dispatcher
}
