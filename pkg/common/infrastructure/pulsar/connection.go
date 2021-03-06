package pulsar

import (
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cenkalti/backoff"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"time"
)

const (
	connectionTimeout     = time.Minute
	maxTestConnectionTime = time.Minute
)

type Config struct {
	Address string
}

type ProducerConfig struct {
	Topic string
}

type ConsumerConfig struct {
	Topic            string
	SubscriptionName string
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
	return c.client.Subscribe(consumerConfig)
}

func (c *connection) Close() {
	c.client.Close()
}

func testCreateProducer(client pulsar.Client) error {
	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = maxTestConnectionTime

	err := backoff.Retry(func() error {
		p, err := client.CreateProducer(pulsar.ProducerOptions{
			Topic: "non-persistent://public/default/test-topic",
		})
		if err != nil {
			return err
		}
		p.Close()
		return nil
	}, exponentialBackOff)
	if err != nil {
		return fmt.Errorf("failed to create test producer: %v", err)
	}
	return nil
}

func NewConnection(config Config, logger logger.Logger) (Connection, error) {
	c, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               fmt.Sprintf("pulsar://%v", config.Address),
		ConnectionTimeout: connectionTimeout,
		Logger:            &loggerAdapter{logger},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %v", err)
	}

	err = testCreateProducer(c)
	if err != nil {
		return nil, err
	}
	return &connection{c}, nil
}
