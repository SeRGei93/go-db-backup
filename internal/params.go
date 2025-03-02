package internal

import (
	"flag"
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

func InitFlags() {
	flag.Bool("backup", false, "Выполнить бэкап базы данных")
	flag.Bool("restore", false, "Выполнить восстановление дампа")
	flag.Bool("docker", false, "Восстановление в докер контейнер")

	// Определяем флаги для SSH
	flag.String("ssh_host", "", "SSH хост")
	flag.Int("ssh_port", 22, "SSH порт")
	flag.String("ssh_user", "", "SSH пользователь")
	flag.String("ssh_password", "", "SSH пароль")

	// Определяем флаги для базы данных
	flag.String("db_name", "", "Имя базы данных")
	flag.String("db_user", "", "Пользователь базы данных")
	flag.String("db_password", "", "Пароль базы данных")
	flag.String("db_host", "", "Хост базы данных")

	// Определяем флаги для восстановления базы данных
	flag.String("restore_db_name", "", "Имя базы данных для восстановления")
	flag.String("restore_db_user", "", "Пользователь базы данных для восстановления")
	flag.String("restore_db_password", "", "Пароль базы данных для восстановления")
	flag.String("restore_db_host", "", "Хост базы данных для восстановления")

	// Парсим все флаги
	flag.Parse()
}

func InitParamsSSHFromFlags() (SshConnectionParams, error) {
	host := flag.Lookup("ssh_host").Value.String()
	portStr := flag.Lookup("ssh_port").Value.String()
	user := flag.Lookup("ssh_user").Value.String()
	password := flag.Lookup("ssh_password").Value.String()

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return SshConnectionParams{}, fmt.Errorf("неверный порт SSH: %v", err)
	}

	// Проверяем, что обязательные параметры заданы
	if host == "" || user == "" || password == "" {
		return SshConnectionParams{}, fmt.Errorf("ошибка: все параметры SSH (host, user, password) должны быть заданы")
	}

	return SshConnectionParams{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}, nil
}

func InitParamsDatabaseFromFlags() (DatabaseConnectionParams, error) {
	name := flag.Lookup("db_name").Value.String()
	user := flag.Lookup("db_user").Value.String()
	password := flag.Lookup("db_password").Value.String()
	host := flag.Lookup("db_host").Value.String()

	// Проверяем, что обязательные параметры заданы
	if name == "" || user == "" || password == "" || host == "" {
		return DatabaseConnectionParams{}, fmt.Errorf("ошибка: все параметры базы данных (name, user, password, host) должны быть заданы")
	}

	return DatabaseConnectionParams{
		Name:     name,
		User:     user,
		Password: password,
		Host:     host,
	}, nil
}

func InitParamsRestoreDatabaseFromFlags() (DatabaseConnectionParams, error) {
	name := flag.Lookup("restore_db_name").Value.String()
	user := flag.Lookup("restore_db_user").Value.String()
	password := flag.Lookup("restore_db_password").Value.String()
	host := flag.Lookup("restore_db_host").Value.String()

	// Проверяем, что обязательные параметры заданы
	if name == "" || user == "" || password == "" || host == "" {
		return DatabaseConnectionParams{}, fmt.Errorf("ошибка: все параметры восстановления базы данных (name, user, password, host) должны быть заданы")
	}

	return DatabaseConnectionParams{
		Name:     name,
		User:     user,
		Password: password,
		Host:     host,
	}, nil
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

func InitParamsDatabase(
	envKeyName string,
	envKeyUser string,
	envKeyPass string,
	envKeyHost string,
) (DatabaseConnectionParams, error) {
	// Получаем значения из окружения
	name, err := getEnv(envKeyName)
	if err != nil {
		return DatabaseConnectionParams{}, err
	}

	host, err := getEnv(envKeyHost)
	if err != nil {
		return DatabaseConnectionParams{}, err
	}

	user, err := getEnv(envKeyUser)
	if err != nil {
		return DatabaseConnectionParams{}, err
	}

	password, err := getEnv(envKeyPass)
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
