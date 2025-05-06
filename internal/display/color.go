// Package display реализует отображение информации и цветовой вывод для файлового менеджера.
package display

import (
	"os"

	"github.com/fatih/color"

	"github.com/mattn/go-isatty"
)

var (
	// Проверяем, поддерживает ли терминал цвета
	useColors = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	// DirColor используется для отображения директорий
	DirColor = color.New(color.FgBlue, color.Bold)
	// FileColor используется для отображения обычных файлов
	FileColor = color.New(color.FgWhite)
	// ExecColor используется для отображения исполняемых файлов
	ExecColor = color.New(color.FgGreen, color.Bold)
	// SymlinkColor используется для отображения символических ссылок
	SymlinkColor = color.New(color.FgCyan, color.Bold)
	// HeaderColor используется для заголовков
	HeaderColor = color.New(color.FgYellow, color.Bold)
	// ErrorColor используется для ошибок
	ErrorColor = color.New(color.FgRed, color.Bold)
	// WarningColor используется для предупреждений
	WarningColor = color.New(color.FgYellow)
	// SuccessColor используется для успешных сообщений
	SuccessColor = color.New(color.FgGreen)
	// ArchiveColor используется для архивов
	ArchiveColor = color.New(color.FgMagenta)
	// ImageColor используется для изображений
	ImageColor = color.New(color.FgHiMagenta)
	// AudioColor используется для аудиофайлов
	AudioColor = color.New(color.FgHiBlue)
	// VideoColor используется для видеофайлов
	VideoColor = color.New(color.FgHiCyan)
	// DocumentColor используется для документов
	DocumentColor = color.New(color.FgHiYellow)
)

// EnableColors включает цветной вывод
func EnableColors() {
	color.NoColor = !useColors
}

// DisableColors отключает цветной вывод
func DisableColors() {
	color.NoColor = true
}

// IsColorEnabled проверяет, включен ли цветной вывод
func IsColorEnabled() bool {
	return !color.NoColor
}

// GetColorByFileType возвращает цвет в зависимости от типа файла
func GetColorByFileType(name string, isDir bool, isExec bool) *color.Color {
	if isDir {
		return DirColor
	}

	if isExec {
		return ExecColor
	}

	// Определение по расширению
	switch GetFileExtension(name) {
	case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".rar":
		return ArchiveColor
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico", ".webp":
		return ImageColor
	case ".mp3", ".wav", ".flac", ".ogg", ".m4a", ".aac":
		return AudioColor
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm":
		return VideoColor
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".ods", ".odp", ".md", ".txt":
		return DocumentColor
	default:
		return FileColor
	}
}
