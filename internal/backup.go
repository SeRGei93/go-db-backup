package internal

import (
	"bytes"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

func BackupDatabase(sshClient *ssh.Client, dbParams DatabaseConnectionParams, backupFile string) error {
	fmt.Println("Создаем дамп")

	remoteFile := fmt.Sprintf("/tmp/%s.sql.gz", dbParams.Name)

	cmd := fmt.Sprintf(
		"mysqldump -h %s -u %s --password=\"%s\" %s | gzip > %s",
		dbParams.Host, dbParams.User, dbParams.Password, dbParams.Name, remoteFile,
	)

	// Создаем SSH-сессию
	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("ошибка создания SSH-сессии: %w", err)
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {

		}
	}(session)

	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("ошибка выполнения mysqldump: %w, stderr: %s", err, stderrBuf.String())
	}

	fmt.Println("Дамп готов")

	return DownloadDumpFile(sshClient, remoteFile, backupFile)
}

func GetBackupFileName(dir string) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("ошибка создания каталога: %w", err)
	}
	return fmt.Sprintf("%s/backup_%s.sql.gz", dir, time.Now().Format("20060102_150405")), nil
}

func DownloadDumpFile(sshClient *ssh.Client, remoteFile string, backupFile string) error {
	fmt.Println("Скачиваем дамп")

	// Открываем SFTP-соединение
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("ошибка создания SFTP-клиента: %w", err)
	}
	defer func(sftpClient *sftp.Client) {
		err := sftpClient.Close()
		if err != nil {

		}
	}(sftpClient)

	// Открываем удалённый файл
	remoteFileHandle, err := sftpClient.Open(remoteFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия удалённого файла: %w", err)
	}
	defer func(remoteFileHandle *sftp.File) {
		err := remoteFileHandle.Close()
		if err != nil {

		}
	}(remoteFileHandle)

	// Создаём локальный файл
	localFile, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("ошибка создания локального файла: %w", err)
	}
	defer func(localFile *os.File) {
		err := localFile.Close()
		if err != nil {

		}
	}(localFile)

	// Копируем файл с сервера на локальную машину
	_, err = io.Copy(localFile, remoteFileHandle)
	if err != nil {
		return fmt.Errorf("ошибка копирования бекапа: %w", err)
	}

	// Удаляем временный файл на сервере
	err = sftpClient.Remove(remoteFile)
	if err != nil {
		fmt.Printf("⚠ Не удалось удалить временный файл %s: %v\n", remoteFile, err)
	}

	fmt.Println("Скачивание дампа завершено")

	return nil
}
