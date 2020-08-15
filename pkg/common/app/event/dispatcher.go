package event

import "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"

type Handler interface {
	Handle(e event.Event) error
}

type Dispatcher interface {
	Dispatch(events []event.Event) error
}

type Subscriber interface {
	Subscribe(h Handler)
}

type dispatcher struct {
	handlers []Handler
}

func (d *dispatcher) Dispatch(events []event.Event) error {
	for _, e := range events {
		for _, h := range d.handlers {
			err := h.Handle(e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *dispatcher) Subscribe(h Handler) {
	d.handlers = append(d.handlers, h)
}

var dispatcherImpl = &dispatcher{make([]Handler, 0)}

func NewDispatcher() Dispatcher {
	return dispatcherImpl
}

func NewSubscriber() Subscriber {
	return dispatcherImpl
}
