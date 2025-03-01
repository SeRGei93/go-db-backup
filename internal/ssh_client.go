package internal

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func InitSshClient(credentials SshConnectionParams) (*ssh.Client, error) {
	// Конфигурация клиента SSH
	config := &ssh.ClientConfig{
		User: credentials.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(credentials.password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Подключение к серверу
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", credentials.host, credentials.port), config)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к серверу по SSH: %w", err)
	}

	fmt.Println("Соединение с сервером установлено")
	return sshClient, nil
}
