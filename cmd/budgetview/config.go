package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type config struct {
	DBName                         string
	DBAddress                      string
	DBUser                         string
	DBPassword                     string
	DBConnectionTimeout            time.Duration
	MessageBrokerAddress           string
	MessageBrokerConnectionTimeout time.Duration
	RealtimeAddress                string
	RealtimeAPIKey                 string
}

func parseEnvInt(key string, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	str, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("undefined environment variable %s", key)
	}
	integer, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("failed to convert environment variable %s to int, value: %s", key, str)
	}
	return integer, nil
}

func parseEnvString(key string, err error) (string, error) {
	if err != nil {
		return "", err
	}
	str, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("undefined environment variable %s", key)
	}
	return str, nil
}

func parseConfig() (*config, error) {
	var err error
	dbName, err := parseEnvString("DATABASE_NAME", err)
	dbAddress, err := parseEnvString("DATABASE_ADDRESS", err)
	dbUser, err := parseEnvString("DATABASE_USER", err)
	dbPassword, err := parseEnvString("DATABASE_PASSWORD", err)
	dbConnTimeout, err := parseEnvInt("DATABASE_CONNECTION_TIMEOUT", err)
	messageBrokerAddress, err := parseEnvString("MESSAGE_BROKER_ADDRESS", err)
	messageBrokerConnTimeout, err := parseEnvInt("MESSAGE_BROKER_CONNECTION_TIMEOUT", err)
	realtimeAddress, err := parseEnvString("REALTIME_ADDRESS", err)
	realtimeAPIKey, err := parseEnvString("REALTIME_API_KEY", err)

	if err != nil {
		return nil, err
	}

	return &config{
		dbName,
		dbAddress,
		dbUser,
		dbPassword,
		time.Duration(dbConnTimeout) * time.Second,
		messageBrokerAddress,
		time.Duration(messageBrokerConnTimeout) * time.Second,
		realtimeAddress,
		realtimeAPIKey,
	}, nil
}
