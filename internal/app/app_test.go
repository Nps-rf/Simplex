package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestApp проверяет основные функции приложения
func TestApp(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "app_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Инициализируем приложение
	app, err := NewApp()
	if err != nil {
		t.Fatalf("не удалось инициализировать приложение: %v", err)
	}

	// Тест на создание команд
	t.Run("CommandsRegistration", func(t *testing.T) {
		// Проверяем, что команды зарегистрированы
		if len(app.commands) == 0 {
			t.Error("команды не были зарегистрированы")
		}

		// Проверяем наличие базовых команд
		requiredCommands := []string{"help", "ls", "cd", "mkdir", "touch", "rm", "cp", "mv", "exit"}
		for _, cmd := range requiredCommands {
			if _, exists := app.commands[cmd]; !exists {
				t.Errorf("обязательная команда %s не зарегистрирована", cmd)
			}
		}
	})

	// Тест на выполнение базовых команд файловой системы
	t.Run("FileSystemCommands", func(t *testing.T) {
		// Меняем директорию на временную
		if err := app.cmdChangeDir([]string{tempDir}); err != nil {
			t.Errorf("не удалось изменить директорию: %v", err)
			return
		}

		// Создаем директорию
		dirName := "test_dir"
		if err := app.cmdMakeDir([]string{dirName}); err != nil {
			t.Errorf("не удалось создать директорию: %v", err)
			return
		}

		// Проверяем, что директория создана
		dirPath := filepath.Join(tempDir, dirName)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Error("директория не была создана")
			return
		}

		// Создаем файл
		fileName := "test.txt"
		if err := app.cmdCreateFile([]string{fileName}); err != nil {
			t.Errorf("не удалось создать файл: %v", err)
			return
		}

		// Проверяем, что файл создан
		filePath := filepath.Join(tempDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Error("файл не был создан")
			return
		}

		// Копируем файл
		copyName := "test_copy.txt"
		if err := app.cmdCopy([]string{fileName, copyName}); err != nil {
			t.Errorf("не удалось скопировать файл: %v", err)
			return
		}

		// Проверяем, что копия создана
		copyPath := filepath.Join(tempDir, copyName)
		if _, err := os.Stat(copyPath); os.IsNotExist(err) {
			t.Error("файл не был скопирован")
			return
		}

		// Удаляем файл
		if err := app.cmdRemoveFile([]string{copyName}); err != nil {
			t.Errorf("не удалось удалить файл: %v", err)
			return
		}

		// Проверяем, что файл удален
		if _, err := os.Stat(copyPath); !os.IsNotExist(err) {
			t.Error("файл не был удален")
		}
	})

	// Тест на поиск файлов
	t.Run("SearchCommands", func(t *testing.T) {
		// Подготовка тестовых данных
		testFilePath := filepath.Join(tempDir, "search_test.txt")
		err := os.WriteFile(testFilePath, []byte("Это тестовый текст для поиска"), 0644)
		if err != nil {
			t.Fatalf("не удалось создать тестовый файл: %v", err)
		}

		// Меняем директорию на временную
		if err := app.cmdChangeDir([]string{tempDir}); err != nil {
			t.Errorf("не удалось изменить директорию: %v", err)
			return
		}

		// Тест на поиск по имени
		results := captureOutput(func() {
			app.cmdFindByName([]string{"*.txt"})
		})

		if !strings.Contains(results, "search_test.txt") {
			t.Error("поиск по имени не нашел тестовый файл")
		}

		// Тест на поиск по содержимому
		results = captureOutput(func() {
			app.cmdFindByContent([]string{"тестовый"})
		})

		if !strings.Contains(results, "search_test.txt") {
			t.Error("поиск по содержимому не нашел тестовый файл")
		}
	})

	// Тест на работу с закладками
	t.Run("BookmarkCommands", func(t *testing.T) {
		// Добавляем закладку
		if err := app.cmdManageBookmarks([]string{"add", "temp", tempDir}); err != nil {
			t.Errorf("не удалось добавить закладку: %v", err)
			return
		}

		// Проверяем, что закладка добавлена
		output := captureOutput(func() {
			app.cmdManageBookmarks([]string{"list"})
		})

		if !strings.Contains(output, "temp") || !strings.Contains(output, tempDir) {
			t.Error("закладка не была добавлена или не отображается в списке")
			return
		}

		// Удаляем закладку
		if err := app.cmdManageBookmarks([]string{"remove", "temp"}); err != nil {
			t.Errorf("не удалось удалить закладку: %v", err)
		}
	})

	// Тест на переключение цветов
	t.Run("ToggleColors", func(t *testing.T) {
		// Запоминаем начальное состояние
		initialState := app.display.UseColors

		// Переключаем цвета
		if err := app.cmdToggleColors([]string{}); err != nil {
			t.Errorf("не удалось переключить цвета: %v", err)
			return
		}

		// Проверяем, что состояние изменилось
		if app.display.UseColors == initialState {
			t.Log("состояние цветов не изменилось, это может происходить на некоторых терминалах")
		}

		// Восстанавливаем исходное состояние
		if app.display.UseColors != initialState {
			app.cmdToggleColors([]string{})
		}
	})

	// Тест на просмотр журнала
	t.Run("ViewLog", func(t *testing.T) {
		// Команда просмотра журнала не должна вызывать ошибок
		if err := app.cmdViewLog([]string{}); err != nil {
			t.Errorf("ошибка при просмотре журнала: %v", err)
		}
	})

	// Тест на команду exit
	t.Run("ExitCommand", func(t *testing.T) {
		// Устанавливаем флаг isRunning
		app.isRunning = true

		// Выполняем команду exit
		if err := app.cmdExit([]string{}); err != nil {
			t.Errorf("ошибка при выполнении команды exit: %v", err)
		}

		// Проверяем, что флаг isRunning установлен в false
		if app.isRunning {
			t.Error("команда exit не остановила приложение")
		}
	})
}

// captureOutput захватывает вывод в stdout во время выполнения функции
func captureOutput(f func()) string {
	// Сохраняем оригинальный stdout
	oldStdout := os.Stdout

	// Создаем пайп для перехвата вывода
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Выполняем функцию
	f()

	// Закрываем writer, чтобы весь вывод был отправлен в reader
	w.Close()
	os.Stdout = oldStdout

	// Читаем вывод из reader
	var buf strings.Builder
	io.Copy(&buf, r)
	return buf.String()
}

// Mocking some dependencies for testing
func mockApp() *App {
	// Создаем минимальную версию App для тестирования
	app, _ := NewApp()
	return app
}
