package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogLevel определяет уровень журналирования
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// LogEntry представляет одну запись в журнале
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     LogLevel  `json:"level"`
	Operation string    `json:"operation"`
	Path      string    `json:"path"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
}

// Logger предоставляет функциональность журналирования
type Logger struct {
	LogFile    string
	MaxEntries int
	Level      LogLevel
	entries    []LogEntry
}

// NewLogger создает новый экземпляр Logger
func NewLogger() (*Logger, error) {
	// Определение пути к файлу журнала
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить домашнюю директорию: %w", err)
	}

	configDir := filepath.Join(homeDir, ".filemanager")
	logFile := filepath.Join(configDir, "operations.log")

	// Создаем директорию конфигурации, если она не существует
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf("не удалось создать директорию конфигурации: %w", err)
		}
	}

	logger := &Logger{
		LogFile:    logFile,
		MaxEntries: 1000, // По умолчанию храним 1000 записей
		Level:      INFO, // По умолчанию логируем INFO и выше
		entries:    []LogEntry{},
	}

	// Загружаем существующие записи, если файл существует
	if _, err := os.Stat(logFile); err == nil {
		err = logger.LoadLog()
		if err != nil {
			return logger, fmt.Errorf("не удалось загрузить журнал: %w", err)
		}
	}

	return logger, nil
}

// Log добавляет новую запись в журнал
func (l *Logger) Log(level LogLevel, operation, path, message string, err error) {
	// Проверяем уровень журналирования
	if level < l.Level {
		return
	}

	// Создаем новую запись
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Operation: operation,
		Path:      path,
		Message:   message,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	// Добавляем запись в память
	l.entries = append(l.entries, entry)

	// Если достигли максимального количества записей, удаляем самые старые
	if len(l.entries) > l.MaxEntries {
		l.entries = l.entries[len(l.entries)-l.MaxEntries:]
	}

	// Сохраняем журнал на диск
	l.SaveLog()
}

// Debug логирует отладочное сообщение
func (l *Logger) Debug(operation, path, message string, err error) {
	l.Log(DEBUG, operation, path, message, err)
}

// Info логирует информационное сообщение
func (l *Logger) Info(operation, path, message string, err error) {
	l.Log(INFO, operation, path, message, err)
}

// Warning логирует предупреждение
func (l *Logger) Warning(operation, path, message string, err error) {
	l.Log(WARNING, operation, path, message, err)
}

// Error логирует ошибку
func (l *Logger) Error(operation, path, message string, err error) {
	l.Log(ERROR, operation, path, message, err)
}

// GetEntries возвращает список записей журнала
func (l *Logger) GetEntries(maxEntries int) []LogEntry {
	if maxEntries <= 0 || maxEntries > len(l.entries) {
		return l.entries
	}

	return l.entries[len(l.entries)-maxEntries:]
}

// SaveLog сохраняет журнал в файл
func (l *Logger) SaveLog() error {
	// Сериализуем записи в JSON
	data, err := json.MarshalIndent(l.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать журнал: %w", err)
	}

	// Записываем данные в файл
	err = os.WriteFile(l.LogFile, data, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать журнал в файл: %w", err)
	}

	return nil
}

// LoadLog загружает журнал из файла
func (l *Logger) LoadLog() error {
	// Проверяем, существует ли файл
	if _, err := os.Stat(l.LogFile); os.IsNotExist(err) {
		return nil // Файл не существует, но это не ошибка
	}

	// Читаем данные из файла
	data, err := os.ReadFile(l.LogFile)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл журнала: %w", err)
	}

	// Десериализуем JSON
	if len(data) > 0 {
		err = json.Unmarshal(data, &l.entries)
		if err != nil {
			return fmt.Errorf("не удалось десериализовать журнал: %w", err)
		}
	}

	return nil
}

// ClearLog очищает журнал
func (l *Logger) ClearLog() error {
	l.entries = []LogEntry{}
	return l.SaveLog()
}

// FormatLogLevel возвращает строковое представление уровня журналирования
func FormatLogLevel(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// FormatEntryForDisplay форматирует запись журнала для отображения
func FormatEntryForDisplay(entry LogEntry) string {
	timestamp := entry.Timestamp.Format("02.01.2006 15:04:05")
	level := FormatLogLevel(entry.Level)

	result := fmt.Sprintf("[%s] [%s] %s: %s", timestamp, level, entry.Operation, entry.Message)

	if entry.Path != "" {
		result += fmt.Sprintf(" (путь: %s)", entry.Path)
	}

	if entry.Error != "" {
		result += fmt.Sprintf(" [ошибка: %s]", entry.Error)
	}

	return result
}
