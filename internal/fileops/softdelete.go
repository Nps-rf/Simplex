package fileops

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"errors"
	"file-manager/internal/i18n"
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
		return fmt.Errorf(i18n.T("softdelete_home_error"), err)
	}
	trashDir := filepath.Join(home, ".local", "share", "Trash", "files")
	infoDir := filepath.Join(home, ".local", "share", "Trash", "info")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf(i18n.T("softdelete_trashdir_error"), err)
	}
	if err := os.MkdirAll(infoDir, 0755); err != nil {
		return fmt.Errorf(i18n.T("softdelete_infodir_error"), err)
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
		return fmt.Errorf(i18n.T("softdelete_move_error"), err)
	}
	// Создаём .trashinfo
	trashInfo := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n", path, time.Now().Format("2006-01-02T15:04:05"))
	infoPath := filepath.Join(infoDir, fileName+".trashinfo")
	if err := os.WriteFile(infoPath, []byte(trashInfo), 0644); err != nil {
		return fmt.Errorf(i18n.T("softdelete_trashinfo_error"), err)
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
		return fmt.Errorf(i18n.T("softdelete_read_trashinfo_error"), err)
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
		return errors.New(i18n.T("softdelete_origpath_not_found"))
	}
	filePath := filepath.Join(trashDir, fileName)
	if err := os.Rename(filePath, origPath); err != nil {
		return fmt.Errorf(i18n.T("softdelete_restore_error"), err)
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
		return fmt.Errorf(i18n.T("softdelete_home_error"), err)
	}
	trashDir := filepath.Join(home, ".Trash")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf(i18n.T("softdelete_trashdir_error"), err)
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

func (m *macSoftDeleter) RestoreFromTrash(_ string) error {
	return errors.New(i18n.T("softdelete_restore_unsupported_mac"))
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
		return fmt.Errorf(i18n.T("softdelete_userprofile_error"))
	}
	trashDir := filepath.Join(userProfile, "Recycle.Bin")
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return fmt.Errorf(i18n.T("softdelete_trashdir_error"), err)
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

func (w *windowsSoftDeleter) RestoreFromTrash(_ string) error {
	return fmt.Errorf(i18n.T("softdelete_restore_unsupported_win"))
}

func (w *windowsSoftDeleter) EmptyTrash() error {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		return fmt.Errorf(i18n.T("softdelete_userprofile_error"))
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
		return nil, fmt.Errorf(i18n.T("softdelete_userprofile_error"))
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
