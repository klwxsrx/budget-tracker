package main

import (
	"errors"
	"fmt"
	"os"
)

type config struct {
	MessageBrokerAddress string
	DbAddress            string
	DbUser               string
	DbPassword           string
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

func parseConfig() (*config, error) {
	var err error
	messageBrokerAddress, err := parseEnvString("MESSAGE_BROKER_ADDRESS", err)
	dbAddress, err := parseEnvString("DATABASE_ADDRESS", err)
	dbUser, err := parseEnvString("DATABASE_USER", err)
	dbPassword, err := parseEnvString("DATABASE_PASSWORD", err)

	if err != nil {
		return nil, err
	}

	return &config{
		messageBrokerAddress,
		dbAddress,
		dbUser,
		dbPassword,
	}, nil
}
