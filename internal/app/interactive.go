package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"file-manager/internal/i18n"
)

// Start запускает интерактивный режим файлового менеджера
func (a *App) Start() {
	// Инициализируем i18n
	if err := i18n.LoadLocale("ru"); err != nil {
		fmt.Fprintf(os.Stderr, "Предупреждение: не удалось загрузить русскую локаль: %v\n", err)
	}

	a.isRunning = true
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(i18n.T("app_started"))

	for a.isRunning {
		dir, err := a.navigator.GetCurrentDirectory()
		if err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("error")+"\n", err)
			break
		}
		fmt.Printf("\n%s> ", dir)
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "" {
			continue
		}

		err = a.processCommand(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("error")+"\n", err)
		}
	}

	fmt.Println(i18n.T("app_stopped"))
}

// ExecuteCommand выполняет одну команду и завершает работу
func (a *App) ExecuteCommand(command string) error {
	// Инициализируем i18n
	if err := i18n.LoadLocale("ru"); err != nil {
		fmt.Fprintf(os.Stderr, "Предупреждение: не удалось загрузить русскую локаль: %v\n", err)
	}

	if command == "" {
		return fmt.Errorf("пустая команда")
	}

	// Обрабатываем команду
	err := a.processCommand(command)

	// Возвращаем ошибку, если была
	return err
}

// processCommand обрабатывает введенную пользователем команду
func (a *App) processCommand(input string) error {
	parts := strings.Split(input, " ")
	cmdName := parts[0]
	args := parts[1:]

	cmd, exists := a.commands[cmdName]
	if !exists {
		errMsg := fmt.Sprintf(i18n.T("unknown_command"), cmdName)
		fmt.Println(errMsg)
		return errors.New(errMsg)
	}

	err := cmd.Execute(args)
	dir, dirErr := a.navigator.GetCurrentDirectory()
	if err != nil {
		fmt.Printf(i18n.T("error")+"\n", err)
		if dirErr == nil {
			a.logger.Error(cmdName, dir, fmt.Sprintf("Выполнение команды '%s' с аргументами %v", cmdName, args), err)
		}
	} else {
		if dirErr == nil {
			a.logger.Info(cmdName, dir, fmt.Sprintf("Успешное выполнение команды '%s' с аргументами %v", cmdName, args), nil)
		}
	}

	return err
}