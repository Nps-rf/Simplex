package display

import (
	"github.com/fatih/color"
	"os"

	"github.com/mattn/go-isatty"
)

var (
	// Проверяем, поддерживает ли терминал цвета
	useColors = isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())

	// Определяем различные цвета для разных типов файлов
	DirColor      = color.New(color.FgBlue, color.Bold)
	FileColor     = color.New(color.FgWhite)
	ExecColor     = color.New(color.FgGreen, color.Bold)
	SymlinkColor  = color.New(color.FgCyan, color.Bold)
	HeaderColor   = color.New(color.FgYellow, color.Bold)
	ErrorColor    = color.New(color.FgRed, color.Bold)
	WarningColor  = color.New(color.FgYellow)
	SuccessColor  = color.New(color.FgGreen)
	ArchiveColor  = color.New(color.FgMagenta)
	ImageColor    = color.New(color.FgHiMagenta)
	AudioColor    = color.New(color.FgHiBlue)
	VideoColor    = color.New(color.FgHiCyan)
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
