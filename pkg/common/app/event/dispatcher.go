package event

import "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"

type dispatcher struct {
	handlers []event.Handler
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

func (d *dispatcher) Subscribe(h event.Handler) {
	d.handlers = append(d.handlers, h)
}

func NewDispatcher() event.Dispatcher {
	return &dispatcher{make([]event.Handler, 0)}
}
