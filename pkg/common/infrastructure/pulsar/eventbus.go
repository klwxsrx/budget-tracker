package pulsar

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/storedevent"
)

const (
	TopicDomainEvent = "domain_event"

	propertyMessageType = "type"
)

type eventBus struct {
	producer   pulsar.Producer
	serializer messaging.StoredEventSerializer
	ctx        context.Context
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
	_, err = b.producer.Send(b.ctx, msg)
	return err
}

func NewEventBus(ctx context.Context, con Connection, serializer messaging.StoredEventSerializer) (storedevent.Bus, error) {
	producer, err := con.CreateProducer(&ProducerConfig{Topic: TopicDomainEvent})
	if err != nil {
		return nil, err
	}
	return &eventBus{producer, serializer, ctx}, nil
}
