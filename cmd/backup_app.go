package main

import (
	"backup/internal"
	"fmt"
	"os"
)

func main() {
	internal.InitFlags()

	backupFile, err := runBackup()
	if err != nil {
		fmt.Println("Ошибка при создании дампа:", err.Error())
		return
	}

	if internal.RestoreFlag != true {
		fmt.Println("Use --restore to create backup")
		os.Exit(1)
	}

	err = runRestore(backupFile)
	if err != nil {
		fmt.Println("Ошибка при восстановления:", err.Error())
	}

	fmt.Println("\n🚀 Хорошего дня!")
}

func runBackup() (string, error) {
	var backupDir = internal.Dir
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

	return backupFile, nil
}

func runRestore(backupFile string) error {
	if internal.RestoreToContainerFlag == true {
		return internal.RestoreDatabaseToContainer(backupFile)
	} else {
		return internal.RestoreDatabase(backupFile)
	}
}
