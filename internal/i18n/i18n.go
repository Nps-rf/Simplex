// Package i18n реализует функции интернационализации (i18n) для файлового менеджера.
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	translations map[string]string
	currentLang  = "ru"
	mu           sync.RWMutex
)

func findProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	for dir := cwd; dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
	}
	return ""
}

// LoadLocale загружает языковой файл (JSON) по коду языка
func LoadLocale(lang string) error {
	mu.Lock()
	defer mu.Unlock()

	var triedPaths []string
	filename := lang + ".json"
	pathsToTry := []string{
		filepath.Join("internal", "i18n", filename),             // относительный путь
		filepath.Join("..", "..", "internal", "i18n", filename), // путь для тестов из поддиректорий
	}
	// Путь относительно исполняемого файла
	exe, exeErr := os.Executable()
	if exeErr == nil {
		root := filepath.Dir(exe)
		pathsToTry = append(pathsToTry, filepath.Join(root, "internal", "i18n", filename))
	}
	// Путь относительно текущей рабочей директории
	if cwd, err := os.Getwd(); err == nil {
		pathsToTry = append(pathsToTry, filepath.Join(cwd, "internal", "i18n", filename))
	}
	// Путь через переменную окружения FILEMANAGER_ROOT
	if envRoot := os.Getenv("FILEMANAGER_ROOT"); envRoot != "" {
		pathsToTry = append(pathsToTry, filepath.Join(envRoot, "internal", "i18n", filename))
	}
	// Путь относительно корня проекта (ищем go.mod вверх по дереву)
	if projectRoot := findProjectRoot(); projectRoot != "" {
		pathsToTry = append(pathsToTry, filepath.Join(projectRoot, "internal", "i18n", filename))
	}

	var file *os.File
	var err error
	for _, path := range pathsToTry {
		file, err = os.Open(path)
		triedPaths = append(triedPaths, path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return fmt.Errorf("не удалось найти языковой файл %s. Пробовал пути: %v. Последняя ошибка: %v", filename, triedPaths, err)
	}
	defer func() { _ = file.Close() }()

	decoder := json.NewDecoder(file)
	trans := make(map[string]string)
	if err := decoder.Decode(&trans); err != nil {
		return err
	}

	translations = trans
	currentLang = lang
	return nil
}

// T возвращает перевод по ключу, если нет — сам ключ
func T(key string, args ...interface{}) string {
	mu.RLock()
	defer mu.RUnlock()
	if translations == nil {
		return key // Явно возвращаем ключ, если локаль не загружена
	}
	msg, ok := translations[key]
	if !ok {
		return key
	}
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// GetCurrentLang возвращает текущий язык
func GetCurrentLang() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}
