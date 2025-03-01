package main

import (
	"backup/internal"
	"fmt"
)

func main() {
	var backupDir = "./backups"
	// Получаем параметры подключения
	sshParams, err := internal.InitParamsSSH()
	if err != nil {
		fmt.Println("Error ssh credentials:", err.Error())
		return
	}

	dbParams, err := internal.InitParamsDatabase()
	if err != nil {
		fmt.Println("Error database credentials:", err.Error())
		return
	}

	sshSession, err := internal.InitSshClient(sshParams)
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err.Error())
		return
	}

	backupFile, err := internal.GetBackupFileName(backupDir)
	if err != nil {
		fmt.Println("Ошибка создания файла:", err.Error())
		return
	}

	err = internal.BackupDatabase(sshSession, dbParams, backupFile)
	if err != nil {
		fmt.Println("Ошибка создания дампа:", err.Error())
		return
	}
	fmt.Println("Дамп успешно создан: ", backupFile)
}
