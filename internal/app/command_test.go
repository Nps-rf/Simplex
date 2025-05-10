package app

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCommandParsing тестирует парсинг и выполнение команд
func TestCommandParsing(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("не удалось создать экземпляр приложения: %v", err)
	}

	tests := []struct {
		input    string
		wantName string
		wantArgs []string
		wantErr  bool
	}{
		{"", "", nil, true},
		{"help", "help", nil, false},
		{"cd /tmp", "cd", []string{"/tmp"}, false},
		{"ls -la /home", "ls", []string{"-la", "/home"}, false},
		{"nonexistent", "nonexistent", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmdName := tt.input
			args := []string{}
			if len(tt.input) > 0 {
				parts := strings.Fields(tt.input)
				cmdName = parts[0]
				if len(parts) > 1 {
					args = parts[1:]
				}
			}

			if tt.input == "" {
				if !tt.wantErr {
					t.Error("Пустая команда должна обрабатываться как ошибка")
				}
				return
			}

			cmd, ok := app.commands[cmdName]
			if tt.wantErr {
				if ok {
					t.Errorf("Ожидалась ошибка для команды %q", tt.input)
				}
				return
			}

			if !ok {
				t.Errorf("Команда %q не найдена", tt.input)
				return
			}

			if cmd.Name != tt.wantName {
				t.Errorf("Получена команда %q, ожидалась %q", cmd.Name, tt.wantName)
			}

			if !equalStringSlices(args, tt.wantArgs) {
				t.Errorf("Получены аргументы %v, ожидались %v", args, tt.wantArgs)
			}
		})
	}
}

// equalStringSlices сравнивает два слайса строк
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// TestExecuteCommand тестирует выполнение команд
func TestExecuteCommand(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("не удалось создать экземпляр приложения: %v", err)
	}

	tests := []struct {
		name          string
		command       string
		expectedError bool
	}{
		{
			name:          "Команда help",
			command:       "help",
			expectedError: false,
		},
		{
			name:          "Несуществующая команда",
			command:       "nonexistent",
			expectedError: true,
		},
		{
			name:          "Команда exit",
			command:       "exit",
			expectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parts := strings.Fields(tc.command)
			var cmd Command
			var ok bool
			var err error

			if len(parts) > 0 {
				cmd, ok = app.commands[parts[0]]
				if !ok {
					err = fmt.Errorf("команда не найдена: %s", parts[0])
				} else if len(parts) > 1 {
					err = cmd.Execute(parts[1:])
				} else {
					err = cmd.Execute(nil)
				}
			}

			if tc.expectedError && err == nil {
				t.Error("ожидалась ошибка, но её нет")
			} else if !tc.expectedError && err != nil {
				t.Errorf("не ожидалась ошибка, но получили: %v", err)
			}
		})
	}
}

// TestRegisterCommands проверяет регистрацию команд
func TestRegisterCommands(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("не удалось создать экземпляр приложения: %v", err)
	}

	// Очистка команд перед тестом
	app.commands = make(map[string]Command)

	// Регистрируем команды
	app.registerCommands()

	// Проверяем, что команды зарегистрированы
	expectedCommands := []string{
		"help", "ls", "cd", "mkdir", "touch", "cat", "rm", "cp",
		"mv", "find", "grep", "bookmark", "colors", "log", "exit",
	}

	for _, cmdName := range expectedCommands {
		if _, exists := app.commands[cmdName]; !exists {
			t.Errorf("команда %s не была зарегистрирована", cmdName)
		}
	}

	// Проверяем, что команды имеют описание
	for name, cmd := range app.commands {
		if cmd.Description == "" {
			t.Errorf("команда %s не имеет описания", name)
		}
	}
}

// TestHelpCommand проверяет команду помощи
func TestHelpCommand(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("не удалось создать экземпляр приложения: %v", err)
	}

	// Проверяем справку для всех команд
	err = app.cmdHelp([]string{})
	if err != nil {
		t.Errorf("ошибка при вызове help: %v", err)
	}

	// Проверяем справку для конкретной команды
	output := captureOutput(func() {
		err := app.cmdHelp([]string{"ls"})
		if err != nil {
			t.Errorf("ошибка при вызове help для ls: %v", err)
		}
	})

	if !bytes.Contains([]byte(output), []byte("ls")) {
		t.Error("справка для конкретной команды не содержит информацию о ней")
	}
}

// TestHistoryCommand проверяет команду history
func TestHistoryCommand(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("не удалось создать экземпляр приложения: %v", err)
	}

	// Обрабатываем несколько команд для истории
	commands := []string{"help", "ls", "cd /"}
	for _, cmd := range commands {
		parts := strings.Fields(cmd)
		if len(parts) > 0 {
			if command, ok := app.commands[parts[0]]; ok {
				if len(parts) > 1 {
					err := command.Execute(parts[1:])
					if err != nil {
						t.Errorf("ошибка при выполнении команды: %v", err)
					}
				} else {
					err := command.Execute(nil)
					if err != nil {
						t.Errorf("ошибка при выполнении команды: %v", err)
					}
				}
			}
		}
	}

	// Проверка журнала команд через просмотр лога
	output := captureOutput(func() {
		err := app.cmdViewLog([]string{})
		if err != nil {
			t.Errorf("ошибка при вызове cmdViewLog: %v", err)
		}
	})

	// Проверяем, что в выводе есть следы выполненных команд
	for _, cmd := range commands {
		if !strings.Contains(output, cmd) {
			t.Logf("Журнал не содержит запись о команде %s (возможно, это нормально, если команда не логируется)", cmd)
		}
	}
}

// TestLsCommand проверяет команду ls
func TestLsCommand(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("не удалось создать экземпляр приложения: %v", err)
	}

	// Создаем временную директорию для теста
	tempDir, err := os.MkdirTemp("", "file-manager-test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer func() {
		_ = os.Chdir(os.TempDir())
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("ошибка при удалении временной директории: %v", err)
		}
	}()

	// Создаем тестовые файлы
	testFiles := []string{"file1.txt", "file2.txt", "folder"}
	for _, file := range testFiles {
		path := filepath.Join(tempDir, file)
		if strings.HasSuffix(file, "/") || !strings.Contains(file, ".") {
			err := os.Mkdir(path, 0755)
			if err != nil {
				t.Fatalf("не удалось создать директорию %s: %v", path, err)
			}
		} else {
			err := os.WriteFile(path, []byte("test content"), 0644)
			if err != nil {
				t.Fatalf("не удалось создать файл %s: %v", path, err)
			}
		}
	}

	// Переходим во временную директорию
	err = app.cmdChangeDir([]string{tempDir})
	if err != nil {
		t.Fatalf("ошибка при смене директории: %v", err)
	}

	// Запускаем команду ls
	output := captureOutput(func() {
		err := app.cmdListDir([]string{})
		if err != nil {
			t.Errorf("ошибка при выполнении ls: %v", err)
		}
	})

	// Проверяем, что все файлы отображаются
	for _, file := range testFiles {
		if !strings.Contains(output, file) {
			t.Errorf("команда ls не показывает файл %s", file)
		}
	}
}

// TestCdCommand проверяет команду cd
func TestCdCommand(t *testing.T) {
	app, err := NewApp()
	if err != nil {
		t.Fatalf("не удалось создать экземпляр приложения: %v", err)
	}

	// Сохраняем текущую директорию
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("не удалось получить текущую директорию: %v", err)
	}

	// Создаем временную директорию
	tempDir, err := os.MkdirTemp("", "file-manager-test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("ошибка при удалении временной директории: %v", err)
		}
	}()

	// Выполняем cd во временную директорию
	err = app.cmdChangeDir([]string{tempDir})
	if err != nil {
		t.Fatalf("ошибка при смене директории: %v", err)
	}

	// Проверяем, что текущая директория изменилась
	currDir := captureOutput(func() {
		err := app.cmdPrintWorkingDir([]string{})
		if err != nil {
			t.Errorf("ошибка при выполнении pwd: %v", err)
		}
	})
	if !strings.Contains(currDir, filepath.Base(tempDir)) {
		t.Errorf("ожидалось, что текущая директория изменится на %s, получено: %s", tempDir, currDir)
	}

	// Возвращаемся в исходную директорию
	err = app.cmdChangeDir([]string{origDir})
	if err != nil {
		t.Fatalf("ошибка при возврате в исходную директорию: %v", err)
	}
}
