// Package main — точка входа в файловый менеджер.
package main

import (
	"fmt"
	"os"
	"strings"

	"file-manager/internal/app"
)

func main() {
	fileManager, err := app.NewApp()
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Ошибка при инициализации приложения: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}

	// Проверяем, есть ли аргументы командной строки
	args := os.Args[1:]

	if len(args) == 0 {
		// Если нет аргументов, запускаем интерактивный режим
		fileManager.Start()
	} else {
		// Если есть аргументы, выполняем их как одну команду
		command := strings.Join(args, " ")
		err := fileManager.ExecuteCommand(command)
		if err != nil {
			os.Exit(1)
		}
	}
}
