package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DbName                  string
	DbAddress               string
	DbUser                  string
	DbPassword              string
	DbMaxConnections        int
	EventStoreMigrationsDir string
	AMQPUser                string
	AMQPPassword            string
	AMQPAddress             string
}

func parseEnvString(key string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	str, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.New(fmt.Sprintf("undefined environment variable %v", key))
	}
	return str, nil
}

func parseEnvInt(key string, err error) (int, error) {
	str, err := parseEnvString(key, err)
	if err != nil {
		return 0, err
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("%v must be integer", key))
	}
	return int(num), nil
}

func ParseConfig() (*Config, error) {
	var err error
	dbName, err := parseEnvString("DATABASE_NAME", err)
	dbAddress, err := parseEnvString("DATABASE_ADDRESS", err)
	dbUser, err := parseEnvString("DATABASE_USER", err)
	dbPassword, err := parseEnvString("DATABASE_PASSWORD", err)
	dbMaxConnections, err := parseEnvInt("DATABASE_MAX_CONNECTIONS", err)
	eventStoreMigrationsDir, err := parseEnvString("EVENT_STORE_MIGRATIONS_DIR", err)

	amqpUser, err := parseEnvString("AMQP_USER", err)
	amqpPassword, err := parseEnvString("AMQP_PASSWORD", err)
	amqpAddress, err := parseEnvString("AMQP_ADDRESS", err)

	if err != nil {
		return nil, err
	}

	return &Config{
		dbName,
		dbAddress,
		dbUser,
		dbPassword,
		dbMaxConnections,
		eventStoreMigrationsDir,
		amqpUser,
		amqpPassword,
		amqpAddress,
	}, nil
}
