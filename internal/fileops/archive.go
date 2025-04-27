package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

// Archiver предоставляет функции для работы с архивами
type Archiver struct{}

// NewArchiver создает новый экземпляр Archiver
func NewArchiver() *Archiver {
	return &Archiver{}
}

// ArchiveFiles создает архив из указанных файлов и директорий
func (a *Archiver) ArchiveFiles(sources []string, destination string, format string) error {
	// Определяем формат архива по расширению, если не указан
	if format == "" {
		format = filepath.Ext(destination)
		if format != "" {
			format = format[1:] // Убираем точку в начале
		}
	}

	// Проверяем существование источников
	for _, source := range sources {
		if _, err := os.Stat(source); err != nil {
			return fmt.Errorf("источник %s не существует: %w", source, err)
		}
	}

	// Создаем архиватор в зависимости от формата
	var arc archiver.Archiver
	switch strings.ToLower(format) {
	case "zip":
		arc = archiver.NewZip()
	case "tar":
		arc = archiver.NewTar()
	case "tar.gz", "tgz":
		arc = archiver.NewTarGz()
	case "tar.bz2", "tbz2":
		arc = archiver.NewTarBz2()
	case "tar.xz", "txz":
		arc = archiver.NewTarXz()
	case "rar":
		return fmt.Errorf("создание RAR архивов не поддерживается")
	case "7z":
		return fmt.Errorf("создание 7z архивов не поддерживается")
	default:
		return fmt.Errorf("неподдерживаемый формат архива: %s", format)
	}

	// Добавляем расширение к имени файла, если оно отсутствует
	if filepath.Ext(destination) == "" {
		destination = destination + "." + format
	}

	// Создаем архив
	err := arc.Archive(sources, destination)
	if err != nil {
		return fmt.Errorf("не удалось создать архив: %w", err)
	}

	return nil
}

// ExtractArchive распаковывает архив в указанную директорию
func (a *Archiver) ExtractArchive(source, destination string) error {
	// Проверяем существование архива
	if _, err := os.Stat(source); err != nil {
		return fmt.Errorf("архив %s не существует: %w", source, err)
	}

	// Создаем директорию назначения, если она не существует
	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("не удалось создать директорию %s: %w", destination, err)
	}

	// Определяем формат архива по расширению
	format := filepath.Ext(source)
	if format != "" {
		format = format[1:] // Убираем точку в начале
	}

	// Для обработки .tar.gz и подобных
	if strings.HasSuffix(source, ".tar.gz") || strings.HasSuffix(source, ".tar.bz2") ||
		strings.HasSuffix(source, ".tar.xz") {
		parts := strings.Split(filepath.Base(source), ".")
		if len(parts) >= 3 {
			format = parts[len(parts)-2] + "." + parts[len(parts)-1]
		}
	}

	// Распаковываем архив
	err := archiver.Unarchive(source, destination)
	if err != nil {
		return fmt.Errorf("не удалось распаковать архив: %w", err)
	}

	return nil
}

// ListArchiveContents выводит содержимое архива
func (a *Archiver) ListArchiveContents(source string) ([]string, error) {
	// Проверяем существование архива
	if _, err := os.Stat(source); err != nil {
		return nil, fmt.Errorf("архив %s не существует: %w", source, err)
	}

	// Открываем архив для чтения
	var files []string
	err := archiver.Walk(source, func(f archiver.File) error {
		files = append(files, f.Name())
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать содержимое архива: %w", err)
	}

	return files, nil
}
