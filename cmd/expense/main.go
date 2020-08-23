package main

import (
	"context"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	appLogger "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/logger"
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
	logger := initLogger()
	config, err := ParseConfig()
	if err != nil {
		logger.Fatalf("failed to parse config: %v", err)
	}

	db, client, err := getReadyDatabaseClient(config)
	if err != nil {
		logger.Fatalf("failed to setup db connection: %v", err)
	}
	defer db.CloseConnection()

	bus := infrastructure.NewContainer(client, appLogger.New(logger)).CommandBus()
	server := startServer(bus, logger)

	listenOSKillSignals()
	_ = server.Shutdown(context.Background())
}

func initLogger() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)
	return l
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

func startServer(bus command.Bus, logger *logrus.Logger) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: transport.NewHttpHandler(bus, appLogger.New(logger)),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalf("unable to start the server: %v", err)
		}
	}()
	return srv
}

func listenOSKillSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
