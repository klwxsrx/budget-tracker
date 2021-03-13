package event

import domain "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"

type Handler interface {
	Handle(e domain.Event) error
}

type Dispatcher interface {
	Dispatch(events []domain.Event) error
	Subscribe(h Handler)
}

type dispatcher struct {
	handlers []Handler
}

func (d *dispatcher) Dispatch(events []domain.Event) error {
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

func NewDispatcher() Dispatcher {
	return &dispatcher{make([]Handler, 0)}
}
