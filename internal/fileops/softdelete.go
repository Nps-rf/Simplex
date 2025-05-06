package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// SoftDeleter определяет интерфейс для soft-delete (корзины)
type SoftDeleter interface {
	MoveToTrash(path string) error
	RestoreFromTrash(fileName string) error
	EmptyTrash() error
	ListTrash() ([]string, error)
}

// GetSoftDeleter возвращает платформозависимую реализацию soft-delete
func GetSoftDeleter() SoftDeleter {
	switch runtime.GOOS {
	case "windows":
		return &windowsSoftDeleter{}
	case "darwin":
		return &macSoftDeleter{}
	default:
		return &linuxSoftDeleter{}
	}
}

// --- Linux ---
type linuxSoftDeleter struct{}

func (l *linuxSoftDeleter) MoveToTrash(path string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("не удалось определить домашнюю директорию: %w", err)
	}
	trashDir := filepath.Join(home, ".local", "share", "Trash", "files")
	infoDir := filepath.Join(home, ".local", "share", "Trash", "info")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку Trash: %w", err)
	}
	if err := os.MkdirAll(infoDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку info: %w", err)
	}
	fileName := filepath.Base(path)
	baseName := fileName
	suffix := 1
	for {
		if _, err := os.Stat(filepath.Join(trashDir, fileName)); os.IsNotExist(err) {
			break
		}
		fileName = fmt.Sprintf("%s_%d", baseName, suffix)
		suffix++
	}
	dest := filepath.Join(trashDir, fileName)
	if err := os.Rename(path, dest); err != nil {
		return fmt.Errorf("не удалось переместить файл в корзину: %w", err)
	}
	// Создаём .trashinfo
	trashInfo := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", path, time.Now().Format("2006-01-02T15:04:05"))
	infoPath := filepath.Join(infoDir, fileName+".trashinfo")
	if err := os.WriteFile(infoPath, []byte(trashInfo), 0644); err != nil {
		return fmt.Errorf("не удалось создать trashinfo: %w", err)
	}
	return nil
}

func (l *linuxSoftDeleter) RestoreFromTrash(fileName string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	trashDir := filepath.Join(home, ".local", "share", "Trash", "files")
	infoDir := filepath.Join(home, ".local", "share", "Trash", "info")
	infoPath := filepath.Join(infoDir, fileName+".trashinfo")
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать trashinfo: %w", err)
	}
	lines := strings.Split(string(data), "\n")
	var origPath string
	for _, line := range lines {
		if strings.HasPrefix(line, "Path=") {
			origPath = strings.TrimPrefix(line, "Path=")
			break
		}
	}
	if origPath == "" {
		return fmt.Errorf("оригинальный путь не найден в trashinfo")
	}
	filePath := filepath.Join(trashDir, fileName)
	if err := os.Rename(filePath, origPath); err != nil {
		return fmt.Errorf("не удалось восстановить файл: %w", err)
	}
	if err := os.Remove(infoPath); err != nil {
		return err
	}
	return nil
}

func (l *linuxSoftDeleter) EmptyTrash() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	trashDir := filepath.Join(home, ".local", "share", "Trash", "files")
	infoDir := filepath.Join(home, ".local", "share", "Trash", "info")
	if err := os.RemoveAll(trashDir); err != nil {
		return err
	}
	if err := os.RemoveAll(infoDir); err != nil {
		return err
	}
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(infoDir, 0755); err != nil {
		return err
	}
	return nil
}

func (l *linuxSoftDeleter) ListTrash() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	trashDir := filepath.Join(home, ".local", "share", "Trash", "files")
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files, nil
}

// --- macOS ---
type macSoftDeleter struct{}

func (m *macSoftDeleter) MoveToTrash(path string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("не удалось определить домашнюю директорию: %w", err)
	}
	trashDir := filepath.Join(home, ".Trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку .Trash: %w", err)
	}
	fileName := filepath.Base(path)
	baseName := fileName
	suffix := 1
	for {
		if _, err := os.Stat(filepath.Join(trashDir, fileName)); os.IsNotExist(err) {
			break
		}
		fileName = fmt.Sprintf("%s_%d", baseName, suffix)
		suffix++
	}
	dest := filepath.Join(trashDir, fileName)
	return os.Rename(path, dest)
}

func (m *macSoftDeleter) RestoreFromTrash(fileName string) error {
	// Для macOS: восстановление не реализовано (нет .trashinfo)
	return fmt.Errorf("восстановление не поддерживается на macOS")
}

func (m *macSoftDeleter) EmptyTrash() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	trashDir := filepath.Join(home, ".Trash")
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(trashDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func (m *macSoftDeleter) ListTrash() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	trashDir := filepath.Join(home, ".Trash")
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files, nil
}

// --- Windows ---
type windowsSoftDeleter struct{}

func (w *windowsSoftDeleter) MoveToTrash(path string) error {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return fmt.Errorf("не удалось определить USERPROFILE")
	}
	trashDir := filepath.Join(userProfile, "Recycle.Bin")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf("не удалось создать папку корзины: %w", err)
	}
	fileName := filepath.Base(path)
	baseName := fileName
	suffix := 1
	for {
		if _, err := os.Stat(filepath.Join(trashDir, fileName)); os.IsNotExist(err) {
			break
		}
		fileName = fmt.Sprintf("%s_%d", baseName, suffix)
		suffix++
	}
	dest := filepath.Join(trashDir, fileName)
	return os.Rename(path, dest)
}

func (w *windowsSoftDeleter) RestoreFromTrash(fileName string) error {
	// Для Windows: восстановление не реализовано (нет .trashinfo)
	return fmt.Errorf("восстановление не поддерживается на Windows")
}

func (w *windowsSoftDeleter) EmptyTrash() error {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return fmt.Errorf("не удалось определить USERPROFILE")
	}
	trashDir := filepath.Join(userProfile, "Recycle.Bin")
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := os.RemoveAll(filepath.Join(trashDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func (w *windowsSoftDeleter) ListTrash() ([]string, error) {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return nil, fmt.Errorf("не удалось определить USERPROFILE")
	}
	trashDir := filepath.Join(userProfile, "Recycle.Bin")
	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files, nil
}
