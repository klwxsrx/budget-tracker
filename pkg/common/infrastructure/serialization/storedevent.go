package serialization

import (
	"encoding/json"
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
)

type storedEventMessage struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type storedEventSerializer struct{}

func (s *storedEventSerializer) Serialize(event *storedevent.StoredEvent) ([]byte, error) {
	jsonObj := storedEventMessage{
		ID:   int(event.ID),
		Type: fmt.Sprintf("%v.%v", event.AggregateName, event.Type),
		Data: event.EventData,
	}
	result, err := json.Marshal(jsonObj)
	if err != nil {
		return nil, fmt.Errorf("can't serialize stored event - %s: %v", jsonObj, err)
	}
	return result, nil
}

func NewStoredEventSerializer() messaging.StoredEventSerializer {
	return &storedEventSerializer{}
}
