package storedevent

import (
	"encoding/json"
	"fmt"
	domain "github.com/klwxsrx/budget-tracker/pkg/common/domain/event"
)

type Serializer interface { // TODO: implement json dto of events
	Serialize(event domain.Event) (string, error)
}

type eventSerializer struct{}

func (s *eventSerializer) Serialize(event domain.Event) (string, error) {
	result, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("can't serialize event - %s: %v", event, err)
	}
	return string(result), nil
}

func NewSerializer() Serializer {
	return &eventSerializer{}
}
