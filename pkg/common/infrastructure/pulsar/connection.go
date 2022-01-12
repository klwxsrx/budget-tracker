package pulsar

import (
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cenkalti/backoff"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
)

type Config struct {
	Address           string
	ConnectionTimeout time.Duration
}

type ProducerConfig struct {
	Topic string
}

type ConsumerConfig struct {
	TopicsPattern    string
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
		TopicsPattern:    config.TopicsPattern,
		SubscriptionName: config.SubscriptionName,
		Type:             pulsar.Failover,
	}
	return c.client.Subscribe(consumerConfig)
}

func (c *connection) Close() {
	c.client.Close()
}

func testCreateProducer(client pulsar.Client, testConnectionTimeout time.Duration) error {
	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = testConnectionTimeout

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
		return fmt.Errorf("failed to create test producer: %w", err)
	}
	return nil
}

func NewConnection(config Config, logger log.Logger) (Connection, error) {
	c, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               fmt.Sprintf("pulsar://%v", config.Address),
		ConnectionTimeout: config.ConnectionTimeout,
		Logger:            &loggerAdapter{logger},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %w", err)
	}

	err = testCreateProducer(c, config.ConnectionTimeout)
	if err != nil {
		return nil, err
	}
	return &connection{c}, nil
}
