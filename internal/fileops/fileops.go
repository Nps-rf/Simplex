package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		return fmt.Errorf("не удалось создать файл %s: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("ошибка при закрытии файла %s: %w", path, err))
		}
	}()
	return nil
}

// CreateDirectory создает новую директорию
func (f *FileOperator) CreateDirectory(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать директорию %s: %w", path, err)
	}
	return nil
}

// CopyFile копирует файл из source в destination
func (f *FileOperator) CopyFile(source, destination string) error {
	// Открываем исходный файл
	src, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("не удалось открыть исходный файл %s: %w", source, err)
	}
	defer func() {
		if err := src.Close(); err != nil {
			panic(fmt.Errorf("ошибка при закрытии исходного файла %s: %w", source, err))
		}
	}()

	// Создаем файл назначения
	dst, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("не удалось создать файл назначения %s: %w", destination, err)
	}
	defer func() {
		if err := dst.Close(); err != nil {
			panic(fmt.Errorf("ошибка при закрытии файла назначения %s: %w", destination, err))
		}
	}()

	// Копируем содержимое
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("не удалось скопировать содержимое: %w", err)
	}

	// Получаем информацию о разрешениях исходного файла
	srcInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о файле %s: %w", source, err)
	}

	// Копируем разрешения
	err = os.Chmod(destination, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("не удалось скопировать разрешения файла: %w", err)
	}

	return nil
}

// CopyDirectory рекурсивно копирует директорию из source в destination
func (f *FileOperator) CopyDirectory(source, destination string) error {
	// Получаем информацию об исходной директории
	srcInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о директории %s: %w", source, err)
	}

	// Создаем директорию назначения с теми же разрешениями
	err = os.MkdirAll(destination, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("не удалось создать директорию %s: %w", destination, err)
	}

	// Читаем содержимое исходной директории
	entries, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf("не удалось прочитать директорию %s: %w", source, err)
	}

	// Копируем каждую запись
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destPath := filepath.Join(destination, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return fmt.Errorf("не удалось получить информацию о файле %s: %w", sourcePath, err)
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
		return fmt.Errorf("не удалось переместить %s в %s: %w", source, destination, err)
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
		return fmt.Errorf("не удалось удалить директорию %s: %w", path, err)
	}
	return nil
}
