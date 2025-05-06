package logger

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "logger_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем логгер с временным файлом журнала
	logFile := filepath.Join(tempDir, "test.log")
	logger := &Logger{
		LogFile:    logFile,
		MaxEntries: 5, // Маленькое значение для тестирования ограничения
		Level:      INFO,
		entries:    []LogEntry{},
	}

	// Тест на добавление записей
	t.Run("AddLogEntries", func(t *testing.T) {
		// Добавляем записи разных уровней
		logger.Debug("TEST", "/path/to/file", "Отладочное сообщение", nil)
		logger.Info("TEST", "/path/to/file", "Информационное сообщение", nil)
		logger.Warning("TEST", "/path/to/file", "Предупреждение", nil)
		logger.Error("TEST", "/path/to/file", "Ошибка", errors.New("тестовая ошибка"))

		// Проверяем, что у нас 3 записи (DEBUG не должен быть добавлен из-за уровня INFO)
		entries := logger.GetEntries(0)
		if len(entries) != 3 {
			t.Errorf("неверное количество записей: получено %d, ожидалось 3", len(entries))
		}
	})

	// Тест на ограничение количества записей
	t.Run("LimitLogEntries", func(t *testing.T) {
		// Добавляем еще 3 записи, чтобы превысить лимит
		logger.Info("TEST", "/path/to/file1", "Сообщение 1", nil)
		logger.Info("TEST", "/path/to/file2", "Сообщение 2", nil)
		logger.Info("TEST", "/path/to/file3", "Сообщение 3", nil)

		// Проверяем, что у нас только 5 записей (последние 5)
		entries := logger.GetEntries(0)
		if len(entries) != 5 {
			t.Errorf("неверное количество записей после ограничения: получено %d, ожидалось 5", len(entries))
		}

		// Проверяем, что первой записью является предупреждение
		if entries[0].Level != WARNING {
			t.Errorf("неверный уровень первой записи: получено %d, ожидалось %d", entries[0].Level, WARNING)
		}
	})

	// Тест на сохранение и загрузку журнала
	t.Run("SaveAndLoadLog", func(t *testing.T) {
		// Сохраняем журнал
		err := logger.SaveLog()
		if err != nil {
			t.Errorf("ошибка при сохранении журнала: %v", err)
		}

		// Создаем новый логгер, который загрузит данные из того же файла
		newLogger := &Logger{
			LogFile:    logFile,
			MaxEntries: 5,
			Level:      INFO,
			entries:    []LogEntry{},
		}

		// Загружаем журнал
		err = newLogger.LoadLog()
		if err != nil {
			t.Errorf("ошибка при загрузке журнала: %v", err)
		}

		// Проверяем, что количество записей совпадает
		if len(newLogger.entries) != len(logger.entries) {
			t.Errorf("неверное количество записей после загрузки: получено %d, ожидалось %d",
				len(newLogger.entries), len(logger.entries))
		}
	})

	// Тест на очистку журнала
	t.Run("ClearLog", func(t *testing.T) {
		// Очищаем журнал
		err := logger.ClearLog()
		if err != nil {
			t.Errorf("ошибка при очистке журнала: %v", err)
		}

		// Проверяем, что у нас 0 записей
		entries := logger.GetEntries(0)
		if len(entries) != 0 {
			t.Errorf("неверное количество записей после очистки: получено %d, ожидалось 0", len(entries))
		}
	})

	// Тест на форматирование уровня журнала
	t.Run("FormatLogLevel", func(t *testing.T) {
		levels := map[LogLevel]string{
			DEBUG:   "DEBUG",
			INFO:    "INFO",
			WARNING: "WARNING",
			ERROR:   "ERROR",
			5:       "UNKNOWN", // Неизвестный уровень
		}

		for level, expected := range levels {
			result := FormatLogLevel(level)
			if result != expected {
				t.Errorf("неверное форматирование уровня %d: получено %s, ожидалось %s", level, result, expected)
			}
		}
	})

	// Тест на форматирование записи для отображения
	t.Run("FormatEntryForDisplay", func(t *testing.T) {
		// Создаем тестовую запись
		entry := LogEntry{
			Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			Level:     INFO,
			Operation: "TEST",
			Path:      "/path/to/file",
			Message:   "Тестовое сообщение",
			Error:     "тестовая ошибка",
		}

		// Форматируем запись
		formatted := FormatEntryForDisplay(entry)

		// Проверяем, что все поля присутствуют в результате
		if !contains(formatted, "TEST") ||
			!contains(formatted, "Тестовое сообщение") ||
			!contains(formatted, "/path/to/file") ||
			!contains(formatted, "тестовая ошибка") {
			t.Errorf("не все поля присутствуют в отформатированной записи: %s", formatted)
		}
	})
}

func TestNewLogger(t *testing.T) {
	// Проверка успешного создания логгера
	logger, err := NewLogger()
	if err != nil {
		t.Errorf("ошибка при создании логгера: %v", err)
	}
	if logger == nil {
		t.Error("логгер не был создан")
	}

	// Проверка повторного создания (файл уже существует)
	logger2, err2 := NewLogger()
	if err2 != nil {
		t.Errorf("ошибка при повторном создании логгера: %v", err2)
	}
	if logger2 == nil {
		t.Error("логгер не был создан повторно")
	}
}

// contains проверяет, содержит ли строка подстроку
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
