package main

import (
	"context"
	appLogger "github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	infrastructureLogger "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mongo"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/pulsar"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := infrastructureLogger.New(initLogrus())
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

	client, err := getMongoClient(config, ctx, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup mongo connection")
	}
	defer client.Close()
	defer ctxCancel()
	logger.Info("app is ready")

	listenOSKillSignals()
}

func getPulsarClient(config *config, logger appLogger.Logger) (pulsar.Connection, error) {
	return pulsar.NewConnection(pulsar.Config{Address: config.MessageBrokerAddress}, logger)
}

func getMongoClient(config *config, ctx context.Context, logger appLogger.Logger) (mongo.Connection, error) {
	return mongo.NewConnection(mongo.Config{
		User:     config.DbUser,
		Password: config.DbPassword,
		Address:  config.DbAddress,
	}, ctx, logger)
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
