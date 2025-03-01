package internal

import (
	"fmt"
	"os"
	"strconv"
)

type SshConnectionParams struct {
	host     string
	port     int
	user     string
	password string
}

type DatabaseConnectionParams struct {
	Name     string
	User     string
	Password string
	Host     string
}

func InitParamsSSH() (SshConnectionParams, error) {
	// Получаем значения из окружения
	host, err := getEnv("SSH_HOST")
	if err != nil {
		return SshConnectionParams{}, err
	}

	port, err := getEnvAsInt("SSH_PORT")
	if err != nil {
		return SshConnectionParams{}, err
	}

	user, err := getEnv("SSH_USER")
	if err != nil {
		return SshConnectionParams{}, err
	}

	password, err := getEnv("SSH_PASSWORD")
	if err != nil {
		return SshConnectionParams{}, err
	}

	// Инициализируем глобальную переменную SshParams
	return SshConnectionParams{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}, nil
}

func InitParamsDatabase() (DatabaseConnectionParams, error) {
	// Получаем значения из окружения
	name, err := getEnv("DB_NAME")
	if err != nil {
		return DatabaseConnectionParams{}, err
	}

	host, err := getEnv("DB_HOST")
	if err != nil {
		return DatabaseConnectionParams{}, err
	}

	user, err := getEnv("DB_USER")
	if err != nil {
		return DatabaseConnectionParams{}, err
	}

	password, err := getEnv("DB_PASSWORD")
	if err != nil {
		return DatabaseConnectionParams{}, err
	}

	// Инициализируем глобальную переменную donorDatabase
	return DatabaseConnectionParams{
		Name:     name,
		User:     user,
		Password: password,
		Host:     host,
	}, nil
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s not found", key)
}

func getEnvAsInt(key string) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue, nil
		}
		return 0, fmt.Errorf("invalid value for %s: %v (must be a valid integer)", key, value)
	}
	return 0, fmt.Errorf("environment variable %s not found", key)
}
