package navigation

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
		return nil, fmt.Errorf("не удалось прочитать директорию %s: %w", n.CurrentDir, err)
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
func (n *Navigator) ChangeDirectory(path string) error {
	// Обработка специальных случаев
	if path == ".." {
		// Переход на уровень выше
		parent := filepath.Dir(n.CurrentDir)
		n.CurrentDir = parent
		return nil
	}

	// Проверяем, является ли путь абсолютным или относительным
	var targetPath string
	if filepath.IsAbs(path) {
		targetPath = path
	} else {
		targetPath = filepath.Join(n.CurrentDir, path)
	}

	// Проверяем, существует ли директория
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о пути %s: %w", targetPath, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s не является директорией", targetPath)
	}

	n.CurrentDir = targetPath
	return nil
}

// GetCurrentDirectory возвращает текущую директорию
func (n *Navigator) GetCurrentDirectory() string {
	return n.CurrentDir
}
