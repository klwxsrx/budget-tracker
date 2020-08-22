package main

import (
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logger := initLogger()
	config, err := ParseConfig()
	if err != nil {
		logger.Fatalf("failed to parse config: %v", err)
	}
	dbClient, err := getReadyDatabaseClient(config)
	if err != nil {
		logger.Fatalf("failed to setup db connection: %v", err)
	}
	_ = infrastructure.NewContainer(dbClient).CommandBus()
}

func initLogger() *logrus.Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)
	return l
}

func getReadyDatabaseClient(config *Config) (mysql.TransactionalClient, error) {
	db := mysql.NewDatabase(mysql.Dsn{
		User:     config.DbUser,
		Password: config.DbPassword,
		Address:  config.DbAddress,
		Database: config.DbName,
	}, config.DbMaxConnections)
	err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	client, err := db.GetClient()
	if err != nil {
		return nil, err
	}
	eventStoreMigration, err := mysql.NewMigration(client, config.EventStoreMigrationsDir)
	if err != nil {
		return nil, err
	}
	err = eventStoreMigration.Migrate()
	if err != nil {
		return nil, err
	}
	return client, nil
}
