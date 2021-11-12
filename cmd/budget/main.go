package main

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/transport"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	commonapplogger "github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	commoninfrastructurelogger "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/pulsar"

	"github.com/sirupsen/logrus"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := commoninfrastructurelogger.New(initLogrus())
	config, err := parseConfig()
	if err != nil {
		logger.WithError(err).Fatal("failed to parse config")
	}

	db, client, err := getReadyDatabaseClient(config, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup db connection")
	}
	defer db.Close()

	broker, err := getPulsarClient(config, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup broker connection")
	}
	defer broker.Close()

	container, err := infrastructure.NewContainer(client, broker, logger)
	if err != nil {
		logger.WithError(err).Fatal(err.Error())
	}
	defer container.Stop()

	server := startServer(container.CommandBus(), logger)
	logger.Info("app is ready")

	listenOSKillSignals()
	_ = server.Shutdown(context.Background())
}

func initLogrus() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)
	return l
}

func getPulsarClient(config *config, logger commonapplogger.Logger) (pulsar.Connection, error) {
	return pulsar.NewConnection(pulsar.Config{
		Address:           config.MessageBrokerAddress,
		ConnectionTimeout: config.MessageBrokerConnectionTimeout,
	}, logger)
}

func getReadyDatabaseClient(config *config, logger commonapplogger.Logger) (mysql.Connection, mysql.TransactionalClient, error) {
	db, err := mysql.NewConnection(mysql.Config{
		DSN: mysql.Dsn{
			User:     config.DBUser,
			Password: config.DBPassword,
			Address:  config.DBAddress,
			Database: config.DBName,
		},
		ConnectionTimeout: config.DBConnectionTimeout,
	}, logger)
	if err != nil {
		return nil, nil, err
	}
	client, err := db.Client()
	if err != nil {
		db.Close()
		return nil, nil, err
	}
	migration, err := mysql.NewMigration(client, logger, config.DBMigrationsDir)
	if err != nil {
		db.Close()
		return nil, nil, err
	}
	err = migration.Migrate()
	if err != nil {
		db.Close()
		return nil, nil, err
	}
	return db, client, nil
}

func startServer(bus command.Bus, logger commonapplogger.Logger) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: transport.NewHTTPHandler(bus, logger),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.WithError(err).Fatal("unable to start the server")
		}
	}()
	return srv
}

func listenOSKillSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
