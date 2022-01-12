package command

import (
	"errors"
	"fmt"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
	"github.com/klwxsrx/budget-tracker/pkg/common/domain"
)

type Type string

type Command interface {
	Type() Type
}

type Base struct {
	CommandType string
}

func (b *Base) Type() Type {
	return Type(b.CommandType)
}

type Bus interface {
	Publish(c Command) error
}

type BusRegistry interface {
	Bus
	Register(h Handler) error
}

type Handler interface {
	Execute(c Command) error
	Type() Type
}

type bus struct {
	registry map[Type]Handler
	logger   log.Logger
}

func (b *bus) Publish(c Command) error {
	handler, ok := b.registry[c.Type()]
	if !ok {
		err := fmt.Errorf("cannot find handler for %v", c.Type())
		b.logger.Error(err.Error())
		return err
	}

	err := handler.Execute(c)
	loggerWithFields := b.logger.WithError(err).With(log.Fields{
		"command": c.Type(),
	})

	if err == nil || errors.Is(err, domain.Error) {
		loggerWithFields.Info("command handled")
	} else {
		loggerWithFields.With(log.Fields{
			"data": c,
		}).Error("command handled with error")
	}
	return err
}

func (b *bus) Register(h Handler) error {
	if _, exists := b.registry[h.Type()]; exists {
		return fmt.Errorf("handler is already set for %v", h.Type())
	}
	b.registry[h.Type()] = h
	return nil
}

func NewBusRegistry(logger log.Logger) BusRegistry {
	return &bus{make(map[Type]Handler), logger}
}
