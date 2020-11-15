package main

import (
	"context"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	appLogger "github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/amqp"
	infrastructureLogger "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/logger"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/transport"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := infrastructureLogger.New(initLogrus())
	config, err := ParseConfig()
	if err != nil {
		logger.With(appLogger.Fields{
			"error": err,
		}).Fatal("failed to parse config")
	}

	db, client, err := getReadyDatabaseClient(config)
	if err != nil {
		logger.With(appLogger.Fields{
			"error": err,
		}).Fatal("failed to setup db connection")
		os.Exit(1)
	}
	defer db.CloseConnection()

	container := infrastructure.NewContainer(client, logger)

	amqpConn, err := getReadyAmqpConnection(config, container, logger)
	if err != nil {
		logger.With(appLogger.Fields{
			"error": err,
		}).Fatal("failed to setup amqp connection")
		os.Exit(1)
	}
	defer amqpConn.Close()

	server := startServer(container.CommandBus(), logger)
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

func getReadyAmqpConnection(config *Config, container infrastructure.Container, logger appLogger.Logger) (amqp.Connection, error) {
	amqpConfig := amqp.Config{
		User:     config.AMQPUser,
		Password: config.AMQPPassword,
		Address:  config.AMQPAddress,
	}
	conn := amqp.NewConnection(amqpConfig, logger)
	conn.AddChannel(container.EventNotifierChannel())
	err := conn.Open()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func getReadyDatabaseClient(config *Config) (mysql.Database, mysql.TransactionalClient, error) {
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
	eventStoreMigration, err := mysql.NewMigration(client, config.EventStoreMigrationsDir)
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
			logger.With(appLogger.Fields{
				"error": err,
			}).Fatal("unable to start the server")
		}
	}()
	return srv
}

func listenOSKillSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
