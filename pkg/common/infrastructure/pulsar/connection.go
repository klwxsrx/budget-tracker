package pulsar

import (
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/cenkalti/backoff"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"time"
)

const (
	maxConnectionTime = time.Minute
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

func newClientExponentialBackoff(config Config, logger logger.Logger) (pulsar.Client, error) {
	var client pulsar.Client
	err := backoff.Retry(func() error {
		var err error
		client, err = pulsar.NewClient(pulsar.ClientOptions{
			URL:    fmt.Sprintf("pulsar://%v", config.Address),
			Logger: &loggerAdapter{logger},
		})
		return err
	}, newOpenConnectionBackoff())
	return client, err
}

func newOpenConnectionBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = maxConnectionTime
	return b
}

func NewConnection(config Config, logger logger.Logger) (Connection, error) {
	c, err := newClientExponentialBackoff(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %v", err)
	}
	return &connection{c}, nil
}
