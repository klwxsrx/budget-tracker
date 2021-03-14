package serialization

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

type storedEventMessage struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
	Data []byte    `json:"data"`
}

type storedEventSerializer struct{}

func (s *storedEventSerializer) Serialize(event *storedevent.StoredEvent) ([]byte, error) {
	jsonObj := storedEventMessage{
		ID:   event.ID.UUID,
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
