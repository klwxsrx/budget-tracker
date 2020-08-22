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

func ParseConfig() (*Config, error) {
	var err error
	dbName, err := parseEnvString("DATABASE_NAME", err)
	dbAddress, err := parseEnvString("DATABASE_ADDRESS", err)
	dbUser, err := parseEnvString("DATABASE_USER", err)
	dbPassword, err := parseEnvString("DATABASE_PASSWORD", err)
	dbMaxConnectionsStr, err := parseEnvString("DATABASE_MAX_CONNECTIONS", err)
	eventStoreMigrationsDir, err := parseEnvString("EVENT_STORE_MIGRATIONS_DIR", err)
	if err != nil {
		return nil, err
	}
	dbMaxConnections, err := strconv.ParseInt(dbMaxConnectionsStr, 10, 64)
	if err != nil {
		return nil, errors.New("DATABASE_MAX_CONNECTIONS must be unsigned int")
	}

	return &Config{
		dbName,
		dbAddress,
		dbUser,
		dbPassword,
		int(dbMaxConnections),
		eventStoreMigrationsDir,
	}, nil
}
