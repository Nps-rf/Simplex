package display

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"file-manager/internal/i18n"

	"gopkg.in/djherbis/times.v1"
)

// FileInfo хранит информацию о файле для отображения
type FileInfo struct {
	Name         string
	Path         string
	Size         int64
	IsDir        bool
	Mode         os.FileMode
	LastModified time.Time
	CreatedAt    time.Time
	IsExecutable bool
}

// Display предоставляет функции для отображения информации о файлах
type Display struct {
	UseColors bool
}

// NewDisplay создает новый экземпляр Display
func NewDisplay() *Display {
	EnableColors() // По умолчанию включаем цвета
	return &Display{
		UseColors: IsColorEnabled(),
	}
}

// GetFileInfo получает подробную информацию о файле или директории
func (d *Display) GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("display_stat_error"), path, err)
	}

	// Проверяем, является ли файл исполняемым
	isExecutable := false
	if !info.IsDir() {
		isExecutable = info.Mode()&0111 != 0
	}

	// Получаем время создания файла, если возможно
	var createdAt time.Time
	t, err := times.Stat(path)
	if err == nil && t.HasBirthTime() {
		createdAt = t.BirthTime()
	} else {
		// Если не удалось получить дату создания, используем дату изменения
		createdAt = info.ModTime()
	}

	fileInfo := &FileInfo{
		Name:         info.Name(),
		Path:         path,
		Size:         info.Size(),
		IsDir:        info.IsDir(),
		Mode:         info.Mode(),
		LastModified: info.ModTime(),
		CreatedAt:    createdAt,
		IsExecutable: isExecutable,
	}

	return fileInfo, nil
}

// FormatFileInfo форматирует информацию о файле для вывода
func (d *Display) FormatFileInfo(fileInfo *FileInfo) string {
	var sb strings.Builder

	fileType := i18n.T("file")
	if fileInfo.IsDir {
		fileType = i18n.T("directory")
	}

	if d.UseColors {
		_, _ = HeaderColor.Fprintf(&sb, i18n.T("file_info")+"\n", fileInfo.Name)
	} else {
		sb.WriteString(fmt.Sprintf(i18n.T("file_info")+"\n", fileInfo.Name))
	}

	sb.WriteString(fmt.Sprintf(i18n.T("path")+": %s\n", fileInfo.Path))
	sb.WriteString(fmt.Sprintf(i18n.T("type")+": %s\n", fileType))

	if !fileInfo.IsDir {
		sb.WriteString(fmt.Sprintf(i18n.T("size")+"\n", formatSize(fileInfo.Size)))
	}

	sb.WriteString(fmt.Sprintf(i18n.T("permissions")+"\n", fileInfo.Mode.String()))
	sb.WriteString(fmt.Sprintf(i18n.T("last_modified")+"\n", fileInfo.LastModified.Format("02.01.2006 15:04:05")))
	sb.WriteString(fmt.Sprintf(i18n.T("created_at")+"\n", fileInfo.CreatedAt.Format("02.01.2006 15:04:05")))

	if fileInfo.IsExecutable {
		sb.WriteString(i18n.T("executable") + "\n")
	}

	return sb.String()
}

// FormatDirEntry форматирует запись директории для вывода списка файлов
func (d *Display) FormatDirEntry(entry os.DirEntry, basePath string) (string, error) {
	fullPath := filepath.Join(basePath, entry.Name())

	info, err := entry.Info()
	if err != nil {
		return "", fmt.Errorf(i18n.T("display_stat_error"), fullPath, err)
	}

	var prefix string
	if entry.IsDir() {
		prefix = i18n.T("dir_prefix")
	} else {
		prefix = i18n.T("file_prefix")
	}

	size := ""
	if !entry.IsDir() {
		size = formatSize(info.Size())
	}

	isExec := info.Mode()&0111 != 0

	result := fmt.Sprintf("%s %-30s %-10s %s",
		prefix,
		entry.Name(),
		size,
		info.ModTime().Format("02.01.2006 15:04:05"))

	if d.UseColors {
		color := GetColorByFileType(entry.Name(), entry.IsDir(), isExec)
		return color.Sprint(result), nil
	}

	return result, nil
}

// FormatSearchResults форматирует результаты поиска
func (d *Display) FormatSearchResults(results []string, query string) string {
	var sb strings.Builder

	if d.UseColors {
		_, _ = HeaderColor.Fprintf(&sb, i18n.T("search_results")+"\n", query)
	} else {
		sb.WriteString(fmt.Sprintf(i18n.T("search_results")+"\n", query))
	}

	sb.WriteString(fmt.Sprintf(i18n.T("found_items")+"\n\n", len(results)))

	for i, result := range results {
		if d.UseColors {
			isDir, _ := isDirectory(result)
			color := GetColorByFileType(result, isDir, false)
			sb.WriteString(fmt.Sprintf("%d. ", i+1))
			sb.WriteString(color.Sprint(result) + "\n")
		} else {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, result))
		}
	}

	return sb.String()
}

// Вспомогательная функция для форматирования размера файла
func formatSize(size int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f ТБ", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f ГБ", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f МБ", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f КБ", float64(size)/KB)
	default:
		return fmt.Sprintf("%d Б", size)
	}
}

// Вспомогательная функция для проверки, является ли путь директорией
func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// ToggleColors переключает использование цветов
func (d *Display) ToggleColors() {
	if d.UseColors {
		DisableColors()
	} else {
		EnableColors()
	}
	d.UseColors = IsColorEnabled()
}
