package display

import (
	"path/filepath"
	"strings"
)

// GetFileExtension возвращает расширение файла в нижнем регистре
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

// FormatPath форматирует путь для отображения (сокращает длинные пути)
func FormatPath(path string, maxLength int) string {
	if len(path) <= maxLength {
		return path
	}

	// Сокращаем путь, сохраняя начало и конец
	halfLen := (maxLength - 3) / 2
	return path[:halfLen] + "..." + path[len(path)-halfLen:]
}
