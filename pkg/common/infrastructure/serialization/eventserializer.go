package serialization

import (
	"encoding/json"
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/serialization"
	"github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

type eventSerializer struct{}

func (es *eventSerializer) Serialize(event event.Event) (string, error) {
	result, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("can't serialize event - %s, %v", event, err)
	}
	return string(result), nil
}

func NewEventSerializer() serialization.EventSerializer {
	return &eventSerializer{}
}
