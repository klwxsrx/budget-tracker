package main

import (
	"context"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure"
	"github.com/klwxsrx/budget-tracker/pkg/budget/infrastructure/transport"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	appLogger "github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	infrastructureLogger "github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/pulsar"
	"github.com/sirupsen/logrus"
	"net/http"
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

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	container, err := infrastructure.NewContainer(client, broker, logger, ctx)
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

func getPulsarClient(config *config, logger appLogger.Logger) (pulsar.Connection, error) {
	return pulsar.NewConnection(pulsar.Config{Address: config.MessageBrokerAddress}, logger)
}

func getReadyDatabaseClient(config *config, logger appLogger.Logger) (mysql.Connection, mysql.TransactionalClient, error) {
	db, err := mysql.NewConnection(mysql.Dsn{
		User:     config.DbUser,
		Password: config.DbPassword,
		Address:  config.DbAddress,
		Database: config.DbName,
	}, logger)
	if err != nil {
		return nil, nil, err
	}
	client, err := db.Client()
	if err != nil {
		db.Close()
		return nil, nil, err
	}
	migration, err := mysql.NewMigration(client, logger, config.DbMigrationsDir)
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
