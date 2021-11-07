package pulsar

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

const (
	propertyMessageType       = "type"
	eventBusTopicNameTemplate = "%v_domain_event"
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
		Properties: map[string]string{propertyMessageType: event.Type},
		EventTime:  event.CreatedAt,
		SequenceID: &sequenceID,
	}
	_, err = b.producer.Send(context.Background(), msg)
	return err
}

func NewEventBus(
	con Connection,
	moduleName string,
	serializer messaging.StoredEventSerializer,
) (storedevent.Bus, error) {
	producer, err := con.CreateProducer(&ProducerConfig{Topic: fmt.Sprintf(eventBusTopicNameTemplate, moduleName)})
	if err != nil {
		return nil, err
	}
	return &eventBus{producer, serializer}, nil
}
