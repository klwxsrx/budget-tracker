package main

import (
	"context"
	appLogger "github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	infrastructureLogger "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/logger"
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

	_, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	logger.Info("app is ready")

	listenOSKillSignals()
}

func getPulsarClient(config *config, logger appLogger.Logger) (pulsar.Connection, error) {
	return pulsar.NewConnection(pulsar.Config{URL: config.MessageBrokerAddress}, logger)
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
