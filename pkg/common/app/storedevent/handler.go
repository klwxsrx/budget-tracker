package storedevent

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/persistence"
)

const dispatchPeriod = time.Second

type Handler interface {
	HandleUnsentStoredEvents()
}

type handler struct {
	busHandler   *UnsentEventBusHandler
	logger       logger.Logger
	needDispatch int32
}

func (d *handler) HandleUnsentStoredEvents() {
	atomic.StoreInt32(&d.needDispatch, 1)
}

func (d *handler) start(ctx context.Context) {
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

func NewHandler(
	ctx context.Context,
	unsentEventProvider UnsentEventProvider,
	eventBus Bus,
	sync persistence.Synchronization,
	loggerImpl logger.Logger,
) Handler {
	busHandler := &UnsentEventBusHandler{unsentEventProvider, eventBus, sync}
	dispatcher := &handler{busHandler, loggerImpl, 1}
	dispatcher.start(ctx)
	return dispatcher
}
