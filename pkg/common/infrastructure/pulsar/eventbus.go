package pulsar

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

const (
	EventTopicsPattern = ".+" + eventTopicNamePostfix

	eventTopicNamePostfix  = "_domain_event"
	eventTopicNameTemplate = "%v" + eventTopicNamePostfix

	propertyEventType = "type"
)

type eventBus struct {
	producer   pulsar.Producer
	serializer messaging.StoredEventSerializer
}

func (b *eventBus) Publish(event *storedevent.StoredEvent) error {
	eventMsg, err := b.serializer.Serialize(event)
	if err != nil {
		return err
	}
	sequenceID := int64(event.SurrogateID)
	msg := &pulsar.ProducerMessage{
		Payload:    eventMsg,
		Properties: map[string]string{propertyEventType: event.Type},
		EventTime:  event.CreatedAt,
		SequenceID: &sequenceID,
	}
	_, err = b.producer.Send(context.Background(), msg)
	return err
}

func NewEventBus(
	connection Connection,
	moduleName string,
	serializer messaging.StoredEventSerializer,
) (storedevent.Bus, error) {
	producer, err := connection.CreateProducer(&ProducerConfig{Topic: fmt.Sprintf(eventTopicNameTemplate, moduleName)})
	if err != nil {
		return nil, err
	}
	return &eventBus{producer, serializer}, nil
}
