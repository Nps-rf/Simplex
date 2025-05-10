package navigation

import (
	"file-manager/internal/i18n"
	"fmt"
	"io/fs"
	"os"
	"sort"
)

// Navigator предоставляет функции для навигации по файловой системе
type Navigator struct {
	CurrentDir string
}

// NewNavigator создает новый экземпляр Navigator
func NewNavigator() (*Navigator, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("не удалось получить текущую директорию: %w", err)
	}
	return &Navigator{
		CurrentDir: currentDir,
	}, nil
}

// ListDirectory отображает содержимое текущей директории
func (n *Navigator) ListDirectory() ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(n.CurrentDir)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("nav_readdir"), n.CurrentDir, err)
	}

	// Сортировка: сначала директории, затем файлы
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() && !entries[j].IsDir() {
			return true
		}
		if !entries[i].IsDir() && entries[j].IsDir() {
			return false
		}
		return entries[i].Name() < entries[j].Name()
	})

	return entries, nil
}

// ChangeDirectory изменяет текущую директорию
func (n *Navigator) ChangeDirectory(targetPath string) error {
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf(i18n.T("nav_stat"), targetPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf(i18n.T("nav_notdir"), targetPath)
	}
	err = os.Chdir(targetPath)
	if err != nil {
		return fmt.Errorf(i18n.T("nav_chdir"), targetPath, err)
	}
	n.CurrentDir = targetPath
	return nil
}

// GetCurrentDirectory возвращает текущую директорию
func (n *Navigator) GetCurrentDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf(i18n.T("nav_getwd"), err)
	}
	return dir, nil
}
