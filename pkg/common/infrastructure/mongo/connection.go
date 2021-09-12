package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectionTimeout    = time.Minute
	disconnectionTimeout = time.Minute
)

type Config struct {
	User     string
	Password string
	Address  string
}

type Connection interface {
	Database(name string) *mongo.Database
	Close()
}

type connection struct {
	client *mongo.Client
	ctx    context.Context
	logger logger.Logger
}

func (c *connection) Database(name string) *mongo.Database {
	return c.client.Database(name)
}

func (c *connection) Close() {
	ctx, cancelFunc := context.WithTimeout(c.ctx, disconnectionTimeout)
	defer cancelFunc()

	err := c.client.Disconnect(ctx)
	if err != nil {
		c.logger.WithError(err).Error("failed to close mongo db connection")
	}
}

func NewConnection(ctx context.Context, config Config, loggerImpl logger.Logger) (Connection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%v", config.Address)).SetAuth(options.Credential{
		Username: config.User,
		Password: config.Password,
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	ctx, cancelFunc := context.WithTimeout(ctx, connectionTimeout)
	defer cancelFunc()
	err = client.Connect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open mongo connection: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongo connection: %w", err)
	}

	return &connection{client, ctx, loggerImpl}, nil
}
