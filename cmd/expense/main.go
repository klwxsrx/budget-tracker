package main

import (
	"context"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	appLogger "github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	infrastructureLogger "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/pulsar"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/transport"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// TODO: use error package with stacktrace

func main() {
	logger := infrastructureLogger.New(initLogrus())
	config, err := ParseConfig()
	if err != nil {
		logger.WithError(err).Fatal("failed to parse config")
	}

	db, client, err := getReadyDatabaseClient(config, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup db connection")
	}
	defer db.CloseConnection()

	broker, err := getPulsarClient(config, logger)
	if err != nil {
		logger.WithError(err).Fatal("failed to setup broker connection")
	}
	defer broker.Close()

	container, err := infrastructure.NewContainer(client, broker, logger)
	if err != nil {
		logger.WithError(err).Fatal(err.Error())
	}
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

func getPulsarClient(config *Config, logger appLogger.Logger) (pulsar.Connection, error) {
	return pulsar.NewConnection(pulsar.Config{URL: config.MessageBrokerAddress}, logger)
}

func getReadyDatabaseClient(config *Config, logger appLogger.Logger) (mysql.Database, mysql.TransactionalClient, error) {
	db := mysql.NewDatabase(mysql.Dsn{
		User:     config.DbUser,
		Password: config.DbPassword,
		Address:  config.DbAddress,
		Database: config.DbName,
	}, config.DbMaxConnections)
	err := db.OpenConnection()
	if err != nil {
		return nil, nil, err
	}
	client, err := db.GetClient()
	if err != nil {
		db.CloseConnection()
		return nil, nil, err
	}
	eventStoreMigration, err := mysql.NewMigration(client, logger, config.EventStoreMigrationsDir)
	if err != nil {
		db.CloseConnection()
		return nil, nil, err
	}
	err = eventStoreMigration.Migrate()
	if err != nil {
		db.CloseConnection()
		return nil, nil, err
	}
	return db, client, nil
}

func startServer(bus command.Bus, logger appLogger.Logger) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: transport.NewHttpHandler(bus, logger),
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
