package pulsar

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/messaging"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
)

type eventbus struct {
	producer   pulsar.Producer
	serializer messaging.StoredEventSerializer
	ctx        context.Context
}

func (b *eventbus) Dispatch(event *storedevent.StoredEvent) error {
	eventMsg, err := b.serializer.Serialize(event)
	if err != nil {
		return err
	}
	msg := &pulsar.ProducerMessage{
		Payload: eventMsg,
	}
	_, err = b.producer.Send(b.ctx, msg)
	return err
}

func NewEventBus(con Connection, serializer messaging.StoredEventSerializer, ctx context.Context) (storedevent.Bus, error) {
	producer, err := con.CreateProducer(&ProducerConfig{Topic: domainEventTopic})
	if err != nil {
		return nil, err
	}
	return &eventbus{producer, serializer, ctx}, nil
}
