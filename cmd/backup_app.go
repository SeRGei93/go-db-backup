package main

import (
	"backup/internal"
	"flag"
	"fmt"
	"os"
)

func main() {
	internal.InitFlags()

	backupFlag := flag.Lookup("backup").Value
	restoreFlag := flag.Lookup("restore").Value
	restoreToContainerFlag := flag.Lookup("docker").Value

	if backupFlag == nil {
		fmt.Println("Use --backup to create backup")
		os.Exit(1)
	}

	backupFile, err := runBackup()
	if err != nil {
		fmt.Println("Ошибка создания дампа:", err.Error())
		return
	}

	if restoreFlag == nil {
		fmt.Println("Use --restore to create backup")
		os.Exit(1)
	}

	if restoreToContainerFlag != nil {
		err = internal.RestoreDatabaseToContainer(backupFile)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	} else {
		err = internal.RestoreDatabase(backupFile)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func runBackup() (string, error) {
	var backupDir = "./backups"
	// Получаем параметры подключения
	sshParams, err := internal.InitParamsSSHFromFlags()
	if err != nil {
		return "", err
	}

	dbParams, err := internal.InitParamsDatabaseFromFlags()
	if err != nil {
		return "", err
	}

	sshClient, err := internal.InitSshClient(sshParams)
	if err != nil {
		return "", err
	}

	backupFile, err := internal.GetBackupFileName(backupDir)
	if err != nil {
		return "", err
	}

	err = internal.BackupDatabase(sshClient, dbParams, backupFile)
	if err != nil {
		return "", err
	}
	fmt.Println("Дамп успешно создан: ", backupFile)

	return backupFile, nil
}
