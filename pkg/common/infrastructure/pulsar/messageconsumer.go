package pulsar

import (
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
)

type messageID []byte

func (id messageID) Serialize() []byte {
	return id
}

type MessageConsumer struct {
	handler  messaging.NamedMessageHandler
	consumer pulsar.Consumer
	logger   logger.Logger
}

func (c *MessageConsumer) start() {
	for msg := range c.consumer.Chan() {
		typ, ok := msg.Properties()[propertyMessageType]
		if !ok {
			c.logger.Error(fmt.Sprintf("failed to get message type for %v", msg.ID().Serialize()))
			c.consumer.Ack(msg)
			continue
		}

		err := c.handler.Handle(&messaging.Message{
			ID:        msg.ID().Serialize(),
			Type:      typ,
			Data:      msg.Payload(),
			EventTime: msg.EventTime(),
		})
		if err != nil {
			c.logger.WithError(err).Error(fmt.Sprintf("failed to handle message %v", msg.ID().Serialize()))
			c.consumer.Nack(msg)
			continue
		}
		c.consumer.Ack(msg)
	}
}

func NewMessageConsumer(
	topic string,
	handler messaging.NamedMessageHandler,
	con Connection,
	loggerImpl logger.Logger,
) (*MessageConsumer, error) {
	offset, err := handler.LatestMessageID()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest message: %w", err)
	}
	initialPosition := pulsar.EarliestMessageID()
	if offset != nil {
		initialPosition = messageID(*offset)
	}

	pulsarConsumer, err := con.Subscribe(&ConsumerConfig{
		Topic:            topic,
		SubscriptionName: handler.Name(),
	})
	if err != nil {
		return nil, err
	}
	err = pulsarConsumer.Seek(initialPosition)
	if err != nil {
		return nil, err
	}

	consumer := &MessageConsumer{
		handler:  handler,
		consumer: pulsarConsumer,
		logger:   loggerImpl,
	}
	go consumer.start()
	return consumer, nil
}
