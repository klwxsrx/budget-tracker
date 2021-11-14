package domain

import "github.com/klwxsrx/budget-tracker/pkg/common/domain"

type Aggregate interface {
	GetChanges() []domain.Event
}

type AggregateState interface {
	Apply(e domain.Event) error
	AggregateName() string
}
