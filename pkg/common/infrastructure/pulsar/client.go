package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"time"
)

const connectionTimeout = time.Minute
const operationTimeout = time.Minute

type Config struct {
	URL string
}

type Connection interface {
	CreateProducer(pulsar.ProducerOptions) (pulsar.Producer, error)
	Subscribe(pulsar.ConsumerOptions) (pulsar.Consumer, error)
	Close()
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
	return c, nil
}
