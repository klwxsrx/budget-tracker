package pulsar

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/storedevent"
)

type bus struct {
	producer pulsar.Producer
	ctx      context.Context
}

func (b *bus) Dispatch(event *storedevent.StoredEvent) error {
	msg := &pulsar.ProducerMessage{
		Payload: event.EventData,
	}
	_, err := b.producer.Send(b.ctx, msg)
	return err
}

func NewEventBus(con Connection, ctx context.Context) (storedevent.Bus, error) {
	p, err := con.CreateProducer(&ProducerConfig{Topic: "domain-event"})
	if err != nil {
		return nil, err
	}
	return &bus{p, ctx}, nil
}
