package internal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// RestoreDatabaseToContainer восстанавливает базу данных из дампа в контейнер Docker.
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

	fmt.Println("Восстановление данных")
	loading := make(chan bool)
	go spinner(loading)

	// Создаем pipe для передачи данных между процессами
	pipeReader, pipeWriter := io.Pipe()

	// Формируем команду для восстановления базы данных через docker exec
	cmd := exec.Command(
		"docker", "exec", "-i", dbParams.Host, "/usr/bin/mysql",
		"-u", dbParams.User, "--password="+dbParams.Password, dbParams.Name,
	)

	// Перенаправляем вывод команды gzip в вход mysql
	cmd.Stdin = pipeReader

	// Запускаем команду восстановления базы данных
	go func() {
		// Закрываем pipe, когда команда будет завершена
		defer func(pw *io.PipeWriter) {
			err := pw.Close()
			if err != nil {

			}
		}(pipeWriter)

		// Запускаем gzip, чтобы разархивировать файл и передать данные в pipe
		gzipCmd := exec.Command("gzip", "-dc")
		gzipCmd.Stdin = file
		gzipCmd.Stdout = pipeWriter

		// Выполняем команду gzip
		if err := gzipCmd.Run(); err != nil {
			fmt.Printf("❌ ошибка при распаковке дампа: %v\n", err)
		}
	}()

	// Выполняем команду восстановления базы данных
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("❌ ошибка при восстановлении дампа: %v, output: %s", err, output)
	}

	loading <- true
	fmt.Println("\n✅ Восстановление завершено.")

	if RemoveFlag == true {
		err = os.Remove(dumpFile) // Удаляем файл дампа
		if err != nil {
			return fmt.Errorf("❌ не удалось удалить файл после восстановления: %v", err.Error())
		}

		fmt.Printf("🗑 Файл %s удален", dumpFile)
	}

	return nil
}

// RestoreDatabase восстанавливает базу данных из дампа базу по прямому подключению к ней.
func RestoreDatabase(dumpFile string) error {
	dbParams, err := InitParamsRestoreDatabaseFromFlags()
	if err != nil {
		return err
	}

	fmt.Println("Начинаю восстановление")
	loading := make(chan bool)
	go spinner(loading)

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

	loading <- true
	fmt.Println("Дамп успешно восстановлен:", stdout.String())
	return nil
}
