package serialization

import (
	"encoding/json"
	"fmt"
	appEvent "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	domain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

type eventSerializer struct{}

func (s *eventSerializer) Serialize(event domain.Event) (string, error) {
	result, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("can't serialize event - %s: %v", event, err)
	}
	return string(result), nil
}

func NewEventSerializer() appEvent.Serializer {
	return &eventSerializer{}
}
