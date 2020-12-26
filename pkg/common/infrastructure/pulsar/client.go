package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"time"
)

const (
	connectionTimeout = time.Minute
	operationTimeout  = time.Minute
)

type Config struct {
	URL string
}

type ProducerConfig struct {
	Topic string
}

type ConsumerConfig struct {
	Topic            string
	SubscriptionName string
	InitialPosition  *pulsar.SubscriptionInitialPosition
}

type Connection interface {
	CreateProducer(config *ProducerConfig) (pulsar.Producer, error)
	Subscribe(config *ConsumerConfig) (pulsar.Consumer, error)
	Close()
}

type connection struct {
	client pulsar.Client
}

func (c *connection) CreateProducer(config *ProducerConfig) (pulsar.Producer, error) {
	return c.client.CreateProducer(pulsar.ProducerOptions{
		Topic: config.Topic,
	})
}

func (c *connection) Subscribe(config *ConsumerConfig) (pulsar.Consumer, error) {
	consumerConfig := pulsar.ConsumerOptions{
		Topic:            config.Topic,
		SubscriptionName: config.SubscriptionName,
	}
	if config.InitialPosition != nil {
		consumerConfig.SubscriptionInitialPosition = *config.InitialPosition
	}
	return c.client.Subscribe(consumerConfig)
}

func (c *connection) Close() {
	c.client.Close()
}

func NewConnection(config Config, logger logger.Logger) (Connection, error) {
	c, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               config.URL,
		ConnectionTimeout: connectionTimeout,
		OperationTimeout:  operationTimeout,
		Logger:            &loggerAdapter{logger},
	})
	if err != nil {
		return nil, err
	}
	return &connection{c}, nil
}
