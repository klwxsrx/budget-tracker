package centrifugo

import (
	"context"
	"fmt"

	"github.com/centrifugal/gocent/v3"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/realtime"
)

type Config struct {
	Address string
	APIKey  string
}

type realtimeClient struct {
	client *gocent.Client
}

func (c *realtimeClient) PublishMessage(channel string, data []byte) error {
	_, err := c.client.Publish(context.Background(), channel, data)
	return err
}

func NewRealtimeClient(config Config) realtime.Client {
	client := gocent.New(gocent.Config{
		Addr: fmt.Sprintf("http://%s/api", config.Address),
		Key:  config.APIKey,
	})
	return &realtimeClient{client}
}
