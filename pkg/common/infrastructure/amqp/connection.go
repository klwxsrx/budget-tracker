package amqp

import (
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"github.com/streadway/amqp"
	"time"
)

const maxConnectionTime = time.Minute

type Config struct {
	User     string
	Password string
	Address  string
}

type Channel interface {
	Connect(connection *amqp.Connection) error
}

type Connection interface {
	AddChannel(ch Channel)
	Open() error
	Close()
}

func (c Config) String() string {
	return fmt.Sprintf("amqp://%s:%s@%s/", c.User, c.Password, c.Address)
}

type connection struct {
	config   Config
	con      *amqp.Connection
	channels []Channel
	logger   logger.Logger
}

func (c *connection) AddChannel(ch Channel) {
	c.channels = append(c.channels, ch)
}

func (c *connection) Open() error {
	if c.con != nil {
		c.Close()
	}

	err := backoff.Retry(func() error {
		con, err := amqp.Dial(c.config.String())
		if err != nil {
			return err
		}
		c.con = con
		return nil
	}, newOpenConnectionBackoff())
	if err != nil {
		return err
	}

	for _, ch := range c.channels {
		err := ch.Connect(c.con)
		if err != nil {
			c.Close()
			return err
		}
	}

	closeChan := c.con.NotifyClose(make(chan *amqp.Error))
	go c.processCloseEvent(closeChan)
	return nil
}

func (c *connection) Close() {
	_ = c.con.Close()
	c.con = nil
}

func (c *connection) processCloseEvent(closeCh chan *amqp.Error) {
	err := <-closeCh
	if err != nil {
		return
	}

	c.logger.Error("amqp connection failed, trying to reconnect")
	for {
		err := c.Open()
		if err == nil {
			c.logger.Info("amqp connection re-established")
			break
		} else {
			c.logger.Error("failed to reconnect to amqp")
		}
	}
}

func newOpenConnectionBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = maxConnectionTime
	return b
}

func NewConnection(config Config, logger logger.Logger) Connection {
	return &connection{config: config, logger: logger}
}
