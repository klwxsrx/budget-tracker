package pulsar

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/messaging"
)

type messageID []byte

func (id messageID) Serialize() []byte {
	return id
}

type ExternalEventReader struct {
	subscriber messaging.Subscriber
	consumer   pulsar.Consumer
	context    context.Context
	logger     logger.Logger
}

func (r *ExternalEventReader) consume() {
	for message := range r.consumer.Chan() {
		err := r.subscriber.Handle(message.Payload())
		if err != nil {
			r.logger.WithError(err).Error("failed to handle external event")
			r.consumer.NackID(message.ID())
			continue
		}

		err = r.subscriber.SetOffset(message.ID().Serialize())
		if err != nil {
			r.logger.WithError(err).Error(fmt.Sprintf("failed to set offset to subscriber: %v", r.subscriber.GetName()))
		}
	}
}

func NewExternalEventReader(
	subscriber messaging.Subscriber,
	con Connection,
	ctx context.Context,
	logger logger.Logger,
) (*ExternalEventReader, error) {
	var initialPosition *pulsar.SubscriptionInitialPosition
	if subscriber.GetOffset() == nil {
		earliest := pulsar.SubscriptionPositionEarliest
		initialPosition = &earliest
	}

	consumer, err := con.Subscribe(&ConsumerConfig{
		Topic:            domainEventTopic,
		SubscriptionName: subscriber.GetName(),
		InitialPosition:  initialPosition,
	})
	if err != nil {
		return nil, err
	}

	if subscriber.GetOffset() != nil {
		err = consumer.Seek(messageID(*subscriber.GetOffset()))
		if err != nil {
			return nil, err
		}
	}

	reader := &ExternalEventReader{subscriber, consumer, ctx, logger}
	go reader.consume()
	return reader, nil
}
