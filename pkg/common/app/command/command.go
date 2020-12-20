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

type Result int

const (
	ResultSuccess Result = iota
	ResultInvalidArgument
	ResultNotFound
	ResultDuplicateConflict
	ResultUnknownError
)

type Bus interface {
	Publish(c Command) Result
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
	registry  map[Type]Handler
	logger    logger.Logger
	resultMap map[error]Result
}

func (b *bus) Publish(c Command) Result {
	handler, ok := b.registry[c.GetType()]
	if !ok {
		b.logger.Error(fmt.Sprintf("cannot find handler for %v", c.GetType()))
		return ResultUnknownError
	}

	err := handler.Execute(c)
	result := b.getResultByError(err)

	loggerWithFields := b.logger.WithError(err).With(logger.Fields{
		"command": c.GetType(),
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

func (b *bus) getResultByError(err error) Result {
	if err == nil {
		return ResultSuccess
	}

	for e, r := range b.resultMap {
		if errors.Is(err, e) {
			return r
		}
	}
	return ResultUnknownError
}

func (b *bus) Register(h Handler) error {
	if _, exists := b.registry[h.GetType()]; exists {
		return errors.New(fmt.Sprintf("handler is already set for %v", h.GetType()))
	}
	b.registry[h.GetType()] = h
	return nil
}

func NewBusRegistry(resultMap map[error]Result, logger logger.Logger) BusRegistry {
	return &bus{make(map[Type]Handler), logger, resultMap}
}
