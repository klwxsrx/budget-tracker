package pulsar

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

const (
	EventTopicsPattern = ".+" + eventTopicNamePostfix

	eventTopicNamePostfix  = "_domain_event"
	eventTopicNameTemplate = "%s" + eventTopicNamePostfix
)

type eventBus struct {
	producer pulsar.Producer
}

func (b *eventBus) Publish(event *storedevent.StoredEvent) error {
	msg := &pulsar.ProducerMessage{
		Payload:    event.EventData,
		Properties: map[string]string{propertyMessageType: event.Type},
		EventTime:  event.CreatedAt,
	}
	_, err := b.producer.Send(context.Background(), msg)
	return err
}

func NewEventBus(
	connection Connection,
	moduleName string,
) (storedevent.Bus, error) {
	producer, err := connection.CreateProducer(&ProducerConfig{Topic: fmt.Sprintf(eventTopicNameTemplate, moduleName)})
	if err != nil {
		return nil, err
	}
	return &eventBus{producer}, nil
}
