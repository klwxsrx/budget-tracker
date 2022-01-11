package pulsar

import (
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/messaging"
)

type MessageConsumer struct {
	handler  messaging.MessageHandler
	consumer pulsar.Consumer
	logger   logger.Logger
	stopChan chan struct{}
}

func (c *MessageConsumer) Stop() {
	c.stopChan <- struct{}{}
}

func (c *MessageConsumer) run() {
	for {
		select {
		case msg, ok := <-c.consumer.Chan():
			if !ok {
				return
			}
			c.processMessage(msg)
		case <-c.stopChan:
			return
		}
	}
}
func (c *MessageConsumer) processMessage(msg pulsar.ConsumerMessage) {
	typ, ok := msg.Properties()[propertyEventType]
	if !ok {
		c.logger.Error(fmt.Sprintf("failed to get message type for %v", msg.ID().Serialize()))
		c.consumer.Ack(msg)
		return
	}

	err := c.handler.Handle(messaging.Message{
		ID:        msg.ID().Serialize(),
		Type:      typ,
		Data:      msg.Payload(),
		EventTime: msg.EventTime(),
	})
	if err != nil {
		c.logger.WithError(err).Error(fmt.Sprintf("failed to handle message %s", msg.Payload()))
		c.consumer.Nack(msg)
		return
	}
	c.consumer.Ack(msg)
	c.logger.Info(fmt.Sprintf("message with type %s and id %v successfully handled", typ, msg.ID().Serialize()))
}

func NewMessageConsumer(
	topicsPattern string,
	handler messaging.NamedMessageHandler,
	resetSubscription bool,
	connection Connection,
	loggerImpl logger.Logger,
) (*MessageConsumer, error) {
	pulsarConsumer, err := connection.Subscribe(&ConsumerConfig{
		TopicsPattern:    topicsPattern,
		SubscriptionName: handler.Name(),
	})
	if err != nil {
		return nil, err
	}

	if resetSubscription {
		err = pulsarConsumer.Seek(pulsar.EarliestMessageID())
		if err != nil {
			return nil, err
		}
	}

	consumer := &MessageConsumer{
		handler:  handler,
		consumer: pulsarConsumer,
		logger:   loggerImpl,
	}
	go consumer.run()
	return consumer, nil
}
