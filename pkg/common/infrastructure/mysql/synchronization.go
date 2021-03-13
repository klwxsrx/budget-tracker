package mysql

import (
	"fmt"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/persistence"
)

type synchronization struct {
	client TransactionalClient
}

func (s *synchronization) CriticalSection(name string, f func() error) error {
	l := NewLock(s.client, name)
	err := l.Get()
	if err != nil {
		return fmt.Errorf("can't create lock: %v", err)
	}
	defer l.Release()
	return f()
}

func NewSynchronization(client TransactionalClient) persistence.Synchronization {
	return &synchronization{client}
}
