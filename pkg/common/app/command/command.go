package command

import (
	"errors"
	"fmt"
)

type Type string

type Command interface {
	GetType() Type
}

type Bus interface {
	Publish(c Command) error
}

type Handler interface {
	Execute(c Command) error
	GetType() Type
}

type BusRegistry interface {
	Register(h Handler) error
}

type bus struct {
	registry map[Type]Handler
}

func (b *bus) Publish(c Command) error {
	handler, ok := b.registry[c.GetType()]
	if !ok {
		return errors.New(fmt.Sprintf("cannot find handler for %v", c.GetType()))
	}
	_ = handler.Execute(c) // TODO: log
	return nil
}

func (b *bus) Register(h Handler) error {
	if _, exists := b.registry[h.GetType()]; exists {
		return errors.New(fmt.Sprintf("handler is already set for %v", h.GetType()))
	}
	b.registry[h.GetType()] = h
	return nil
}

var busImpl = &bus{make(map[Type]Handler)}

func NewBus() Bus {
	return busImpl
}

func NewBusRegistry() BusRegistry {
	return busImpl
}
