// Package search реализует функции поиска файлов и по содержимому для файлового менеджера.
package search

import (
	"bufio"
	"file-manager/internal/i18n"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Searcher предоставляет функции для поиска файлов и содержимого
type Searcher struct{}

// NewSearcher создает новый экземпляр Searcher
func NewSearcher() *Searcher {
	return &Searcher{}
}

// SearchByName ищет файлы по шаблону имени
func (s *Searcher) SearchByName(root, pattern string) ([]string, error) {
	var matches []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Проверяем, соответствует ли имя файла шаблону
		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return err
		}

		if matched {
			matches = append(matches, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf(i18n.T("search_files"), err)
	}

	return matches, nil
}

// SearchByContent ищет файлы по содержимому (содержащие указанный текст)
func (s *Searcher) SearchByContent(root, content string) ([]string, error) {
	var matches []string
	var processedFilesMap = make(map[string]bool)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Пропускаем файлы, к которым нет доступа
		}

		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		// Пропускаем слишком большие файлы (более 10 МБ)
		if info.Size() > 10*1024*1024 {
			return nil
		}

		// Проверяем, обрабатывали ли мы уже этот файл
		if processedFilesMap[path] {
			return nil
		}
		processedFilesMap[path] = true

		// Проверяем права на чтение файла
		if info.Mode().Perm()&0400 == 0 {
			return nil // нет прав на чтение
		}

		// Открываем файл для чтения
		file, err := os.Open(path)
		if err != nil {
			return nil // Пропускаем файлы, которые не можем открыть
		}
		defer func() {
			err := file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ошибка при закрытии файла: %v\n", err)
			}
		}()

		// Сканируем файл построчно
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), content) {
				matches = append(matches, path)
				break // Достаточно одного совпадения в файле
			}
		}

		if err := scanner.Err(); err != nil {
			return nil // Пропускаем при ошибке сканирования
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf(i18n.T("search_content"), err)
	}

	return matches, nil
}

// SearchByRegex ищет файлы по содержимому с использованием регулярного выражения
func (s *Searcher) SearchByRegex(root, regexPattern string) ([]string, error) {
	var matches []string
	var processedFilesMap = make(map[string]bool)

	// Компилируем регулярное выражение
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("invalid_regex_pattern"), err)
	}

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Пропускаем файлы, к которым нет доступа
		}

		// Пропускаем директории
		if info.IsDir() {
			return nil
		}

		// Пропускаем слишком большие файлы (более 10 МБ)
		if info.Size() > 10*1024*1024 {
			return nil
		}

		// Проверяем, обрабатывали ли мы уже этот файл
		if processedFilesMap[path] {
			return nil
		}
		processedFilesMap[path] = true

		// Открываем файл для чтения
		file, err := os.Open(path)
		if err != nil {
			return nil // Пропускаем файлы, которые не можем открыть
		}
		defer func() {
			err := file.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ошибка при закрытии файла: %v\n", err)
			}
		}()

		// Сканируем файл построчно
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if regex.MatchString(scanner.Text()) {
				matches = append(matches, path)
				break // Достаточно одного совпадения в файле
			}
		}

		if err := scanner.Err(); err != nil {
			return nil // Пропускаем при ошибке сканирования
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf(i18n.T("search_regex"), err)
	}

	return matches, nil
}
