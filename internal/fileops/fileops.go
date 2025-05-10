package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"file-manager/internal/i18n"
)

// FileOperator предоставляет функции для работы с файлами и директориями
type FileOperator struct {
	SoftDeleter SoftDeleter
}

// NewFileOperator создает новый экземпляр FileOperator
func NewFileOperator() *FileOperator {
	return &FileOperator{
		SoftDeleter: GetSoftDeleter(),
	}
}

// CreateFile создает новый файл
func (f *FileOperator) CreateFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_create_file_error"), path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf(i18n.T("fileops_close_file_error"), path, err))
		}
	}()
	return nil
}

// CreateDirectory создает новую директорию
func (f *FileOperator) CreateDirectory(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_create_dir_error"), path, err)
	}
	return nil
}

// CopyFile копирует файл из source в destination
func (f *FileOperator) CopyFile(source, destination string) error {
	// Открываем исходный файл
	src, err := os.Open(source)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_open_source_error"), source, err)
	}
	defer func() {
		if err := src.Close(); err != nil {
			panic(fmt.Errorf(i18n.T("fileops_close_source_error"), source, err))
		}
	}()

	// Создаем файл назначения
	dst, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_create_dest_error"), destination, err)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			panic(fmt.Errorf(i18n.T("fileops_close_dest_error"), destination, err))
		}
	}()

	// Копируем содержимое
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_copy_content_error"), err)
	}

	// Получаем информацию о разрешениях исходного файла
	srcInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_stat_error"), source, err)
	}

	// Копируем разрешения
	err = os.Chmod(destination, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_chmod_error"), err)
	}

	return nil
}

// CopyDirectory рекурсивно копирует директорию из source в destination
func (f *FileOperator) CopyDirectory(source, destination string) error {
	// Получаем информацию об исходной директории
	srcInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_stat_dir_error"), source, err)
	}

	// Создаем директорию назначения с теми же разрешениями
	err = os.MkdirAll(destination, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_create_dir_error"), destination, err)
	}

	// Читаем содержимое исходной директории
	entries, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_read_dir_error"), source, err)
	}

	// Копируем каждую запись
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(destination, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return fmt.Errorf(i18n.T("fileops_stat_error"), sourcePath, err)
		}

		if fileInfo.IsDir() {
			// Рекурсивно копируем директорию
			if err = f.CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		} else {
			// Копируем файл
			if err = f.CopyFile(sourcePath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveFile перемещает файл или директорию
func (f *FileOperator) MoveFile(source, destination string) error {
	err := os.Rename(source, destination)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_move_error"), source, destination, err)
	}
	return nil
}

// DeleteFile удаляет файл (soft-delete)
func (f *FileOperator) DeleteFile(path string) error {
	return f.SoftDeleter.MoveToTrash(path)
}

// DeleteDirectory рекурсивно удаляет директорию
func (f *FileOperator) DeleteDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf(i18n.T("fileops_delete_dir_error"), path, err)
	}
	return nil
}
