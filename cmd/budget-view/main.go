package main

import (
	commonapplogger "github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	commoninfrastructurelogger "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mongo"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/pulsar"

	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := commoninfrastructurelogger.New(initLogrus())
	config, err := parseConfig()
	if err != nil {
		logger.WithError(err).Fatal("failed to parse config")
	}

	broker, err := getPulsarClient(config, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup broker connection")
	}
	defer broker.Close()

	ctx, ctxCancel := context.WithCancel(context.Background())

	client, err := getMongoClient(ctx, config, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup mongo connection")
	}
	defer client.Close()
	defer ctxCancel()
	logger.Info("app is ready")

	listenOSKillSignals()
}

func getPulsarClient(config *config, logger commonapplogger.Logger) (pulsar.Connection, error) {
	return pulsar.NewConnection(pulsar.Config{Address: config.MessageBrokerAddress}, logger)
}

func getMongoClient(ctx context.Context, config *config, logger commonapplogger.Logger) (mongo.Connection, error) {
	return mongo.NewConnection(ctx,
		mongo.Config{
			User:     config.DbUser,
			Password: config.DbPassword,
			Address:  config.DbAddress,
		}, logger,
	)
}

func initLogrus() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)
	return l
}

func listenOSKillSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
