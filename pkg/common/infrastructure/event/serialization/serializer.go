package serialization

import (
	"encoding/json"
	"fmt"
	app "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	domain "github.com/klwxsrx/expense-tracker/pkg/common/domain/event"
)

type serializer struct{}

func (s *serializer) Serialize(event domain.Event) (string, error) {
	result, err := json.Marshal(event)
	if err != nil {
		return "", fmt.Errorf("can't serialize event - %s: %v", event, err)
	}
	return string(result), nil
}

func NewSerializer() app.Serializer {
	return &serializer{}
}
