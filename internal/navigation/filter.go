package navigation

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FilterOptions содержит опции для фильтрации файлов
type FilterOptions struct {
	// Фильтр по расширению
	Extensions []string

	// Фильтр по имени (поддерживает шаблоны * и ?)
	NamePattern string

	// Фильтр по размеру
	MinSize int64
	MaxSize int64

	// Фильтр по дате изменения
	ModifiedAfter  time.Time
	ModifiedBefore time.Time

	// Тип файла
	ShowDirs   bool
	ShowFiles  bool
	ShowHidden bool
}

// NewFilterOptions создает новый экземпляр FilterOptions с значениями по умолчанию
func NewFilterOptions() *FilterOptions {
	return &FilterOptions{
		Extensions:  []string{},
		NamePattern: "",
		MinSize:     -1,
		MaxSize:     -1,
		ShowDirs:    true,
		ShowFiles:   true,
		ShowHidden:  false,
	}
}

// Filter фильтрует список записей директории согласно опциям фильтрации
func Filter(entries []os.DirEntry, basePath string, options *FilterOptions) ([]os.DirEntry, error) {
	var result []os.DirEntry

	for _, entry := range entries {
		// Проверка скрытых файлов
		if !options.ShowHidden && isHidden(entry.Name()) {
			continue
		}

		// Проверка типа (файл/директория)
		if entry.IsDir() && !options.ShowDirs {
			continue
		}

		if !entry.IsDir() && !options.ShowFiles {
			continue
		}

		// Проверка имени файла
		if options.NamePattern != "" {
			match, err := filepath.Match(options.NamePattern, entry.Name())
			if err != nil {
				return nil, err
			}
			if !match {
				continue
			}
		}

		// Проверка расширения (только для файлов)
		if len(options.Extensions) > 0 {
			if entry.IsDir() {
				continue // Пропускаем директории при фильтрации по расширению
			}

			ext := strings.ToLower(filepath.Ext(entry.Name()))
			if ext != "" {
				ext = ext[1:] // Убираем точку в начале
			}

			matchExtension := false
			for _, allowedExt := range options.Extensions {
				if ext == strings.ToLower(allowedExt) {
					matchExtension = true
					break
				}
			}

			if !matchExtension {
				continue
			}
		}

		// Проверка размера (только для файлов)
		if options.MinSize >= 0 || options.MaxSize >= 0 {
			// Пропускаем директории при фильтрации по размеру
			if entry.IsDir() {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				// Пропускаем файлы, к которым нет доступа
				continue
			}

			size := info.Size()
			// Размер должен быть >= MinSize, если MinSize указан
			if options.MinSize >= 0 && size < options.MinSize {
				continue
			}
			// Размер должен быть <= MaxSize, если MaxSize указан
			if options.MaxSize >= 0 && size > options.MaxSize {
				continue
			}
		}

		// Проверка даты изменения
		if !options.ModifiedAfter.IsZero() || !options.ModifiedBefore.IsZero() {
			info, err := entry.Info()
			if err != nil {
				// Пропускаем файлы, к которым нет доступа
				continue
			}

			modTime := info.ModTime()
			if !options.ModifiedAfter.IsZero() && modTime.Before(options.ModifiedAfter) {
				continue
			}
			if !options.ModifiedBefore.IsZero() && modTime.After(options.ModifiedBefore) {
				continue
			}
		}

		// Если прошли все фильтры, добавляем запись в результат
		result = append(result, entry)
	}

	return result, nil
}

// isHidden проверяет, является ли файл скрытым
func isHidden(name string) bool {
	// В Unix/Linux скрытые файлы начинаются с точки
	return len(name) > 0 && name[0] == '.'
}
