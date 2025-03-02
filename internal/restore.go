package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func RestoreDatabase(dumpFile string) error {
	dbParams, err := InitParamsDatabase(
		"RESTORE_DB_NAME",
		"RESTORE_DB_USER",
		"RESTORE_DB_PASSWORD",
		"RESTORE_DB_HOST",
	)
	if err != nil {
		return err
	}
	// Проверяем, существует ли файл дампа
	if _, err := os.Stat(dumpFile); os.IsNotExist(err) {
		return fmt.Errorf("файл дампа не найден: %s", dumpFile)
	}

	restoreCmd := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("gunzip -c %s | mariadb -h %s -u %s -p%s %s",
			dumpFile, dbParams.Host, dbParams.User, dbParams.Password, dbParams.Name))

	// Выполняем команду
	var stdout, stderr bytes.Buffer
	restoreCmd.Stdout = &stdout
	restoreCmd.Stderr = &stderr

	if err := restoreCmd.Run(); err != nil {
		return fmt.Errorf("ошибка при восстановлении дампа: %s", stderr.String())
	}

	fmt.Println("Дамп успешно восстановлен:", stdout.String())
	return nil
}

func RestoreDatabaseToContainer(dumpFile string) error {
	dbParams, err := InitParamsRestoreDatabaseFromFlags()
	if err != nil {
		return err
	}

	file, err := os.Open(dumpFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл дампа: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Создаем pipe для передачи данных между процессами
	pr, pw := io.Pipe()

	// Формируем команду для восстановления базы данных через docker exec
	cmd := exec.Command(
		"docker", "exec", "-i", dbParams.Host, "/usr/bin/mysql",
		"-u", dbParams.User, "--password="+dbParams.Password, dbParams.Name,
	)

	// Перенаправляем вывод команды gzip в вход mysql
	cmd.Stdin = pr

	// Запускаем команду восстановления базы данных
	go func() {
		// Закрываем pipe, когда команда будет завершена
		defer func(pw *io.PipeWriter) {
			err := pw.Close()
			if err != nil {

			}
		}(pw)

		// Запускаем gzip, чтобы разархивировать файл и передать данные в pipe
		gzipCmd := exec.Command("gzip", "-dc")
		gzipCmd.Stdin = file
		gzipCmd.Stdout = pw

		// Выполняем команду gzip
		if err := gzipCmd.Run(); err != nil {
			fmt.Printf("ошибка при разархивировании дампа: %v\n", err)
		}
	}()

	// Выполняем команду восстановления базы данных
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка при восстановлении дампа: %v, output: %s", err, output)
	}

	fmt.Println("Восстановление завершено успешно.")

	err = os.Remove(dumpFile) // Удаляем файл дампа
	if err != nil {
		return fmt.Errorf("не удалось удалить файл после восстановления: %v", err)
	}

	return nil
}
