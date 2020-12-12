package command

import (
	"errors"
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
)

type Type string

type Command interface {
	GetType() Type
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
	GetType() Type
}

type bus struct {
	registry map[Type]Handler
	logger   logger.Logger
}

func (b *bus) Publish(c Command) error {
	handler, ok := b.registry[c.GetType()]
	if !ok {
		return errors.New(fmt.Sprintf("cannot find handler for %v", c.GetType()))
	}
	err := handler.Execute(c) // TODO: return success/failed result

	loggerWithFields := b.logger.With(logger.Fields{
		"command":     c.GetType(),
		"data":        c,
		"resultError": err,
	})
	if err != nil {
		loggerWithFields.Warn("command handled with error")
	} else {
		loggerWithFields.Info("command handled")
	}
	return nil
}

func (b *bus) Register(h Handler) error {
	if _, exists := b.registry[h.GetType()]; exists {
		return errors.New(fmt.Sprintf("handler is already set for %v", h.GetType()))
	}
	b.registry[h.GetType()] = h
	return nil
}

func NewBusRegistry(logger logger.Logger) BusRegistry {
	return &bus{make(map[Type]Handler), logger}
}
