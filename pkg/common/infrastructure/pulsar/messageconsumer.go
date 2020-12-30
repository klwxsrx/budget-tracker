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

type MessageConsumer struct {
	handler  messaging.MessageHandler
	consumer pulsar.Consumer
	context  context.Context
	logger   logger.Logger
}

func (c *MessageConsumer) start() {
	for msg := range c.consumer.Chan() {
		err := c.handler.Handle(msg.Payload())
		if err != nil {
			c.logger.WithError(err).Error("failed to handle external event")
			c.consumer.NackID(msg.ID())
			continue
		}

		err = c.handler.SetOffset(msg.ID().Serialize())
		if err != nil {
			c.logger.WithError(err).Error(fmt.Sprintf("failed to set offset to consumer: %v", c.handler.GetName()))
		}
	}
}

func NewMessageConsumer(
	topic string,
	handler messaging.MessageHandler,
	con Connection,
	ctx context.Context,
	logger logger.Logger,
) (*MessageConsumer, error) {
	offset, err := handler.GetOffset()
	if err != nil {
		return nil, fmt.Errorf("failed to get initial offset: %v", err)
	}

	var initialPosition *pulsar.SubscriptionInitialPosition
	if offset == nil {
		earliest := pulsar.SubscriptionPositionEarliest
		initialPosition = &earliest
	}

	pulsarConsumer, err := con.Subscribe(&ConsumerConfig{
		Topic:            topic,
		SubscriptionName: handler.GetName(),
		InitialPosition:  initialPosition,
	})
	if err != nil {
		return nil, err
	}

	if offset != nil {
		err = pulsarConsumer.Seek(messageID(*offset))
		if err != nil {
			return nil, err
		}
	}

	consumer := &MessageConsumer{
		handler:  handler,
		consumer: pulsarConsumer,
		context:  ctx,
		logger:   logger,
	}
	go consumer.start()
	return consumer, nil
}
