package main

import (
	"fmt"
	"os"
)

type config struct {
	DbName               string
	DbAddress            string
	DbUser               string
	DbPassword           string
	DbMigrationsDir      string
	MessageBrokerAddress string
}

func parseEnvString(key string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	str, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("undefined environment variable %v", key)
	}
	return str, nil
}

func parseConfig() (*config, error) {
	var err error
	dbName, err := parseEnvString("DATABASE_NAME", err)
	dbAddress, err := parseEnvString("DATABASE_ADDRESS", err)
	dbUser, err := parseEnvString("DATABASE_USER", err)
	dbPassword, err := parseEnvString("DATABASE_PASSWORD", err)
	dbMigrationsDir, err := parseEnvString("DATABASE_MIGRATIONS_DIR", err)
	messageBrokerAddress, err := parseEnvString("MESSAGE_BROKER_ADDRESS", err)

	if err != nil {
		return nil, err
	}

	return &config{
		dbName,
		dbAddress,
		dbUser,
		dbPassword,
		dbMigrationsDir,
		messageBrokerAddress,
	}, nil
}
