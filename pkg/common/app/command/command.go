package command

import (
	"fmt"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
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

type Result int

const (
	ResultSuccess Result = iota
	ResultInvalidArgument
	ResultNotFound
	ResultDuplicateConflict
	ResultUnknownError
)

type ResultTranslator interface {
	Translate(err error) Result
}

type Bus interface {
	Publish(c Command) Result
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
	registry   map[Type]Handler
	logger     logger.Logger
	translator ResultTranslator
}

func (b *bus) Publish(c Command) Result {
	handler, ok := b.registry[c.Type()]
	if !ok {
		b.logger.Error(fmt.Sprintf("cannot find handler for %v", c.Type()))
		return ResultUnknownError
	}

	err := handler.Execute(c)
	result := b.translator.Translate(err)

	loggerWithFields := b.logger.WithError(err).With(logger.Fields{
		"command": c.Type(),
		"data":    c,
		"result":  result,
	})
	if result == ResultUnknownError {
		loggerWithFields.Error("command handled with error")
	} else {
		loggerWithFields.Info("command handled")
	}
	return result
}

func (b *bus) Register(h Handler) error {
	if _, exists := b.registry[h.Type()]; exists {
		return fmt.Errorf("handler is already set for %v", h.Type())
	}
	b.registry[h.Type()] = h
	return nil
}

func NewBusRegistry(translator ResultTranslator, loggerImpl logger.Logger) BusRegistry {
	return &bus{make(map[Type]Handler), loggerImpl, translator}
}
