// Package fileops реализует операции с файлами и архивами для файлового менеджера.
package fileops

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Archiver предоставляет функции для работы с архивами
type Archiver struct{}

// NewArchiver создает новый экземпляр Archiver
func NewArchiver() *Archiver {
	return &Archiver{}
}

// ArchiveFiles создает архив из указанных файлов и директорий (только zip)
func (a *Archiver) ArchiveFiles(sources []string, destination string, format string) error {
	if format == "" {
		format = filepath.Ext(destination)
		if format != "" {
			format = format[1:]
		}
	}
	if strings.ToLower(format) != "zip" {
		return fmt.Errorf("поддерживается только zip-архивация (format=%s)", format)
	}
	zipFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("не удалось создать архив: %w", err)
	}
	defer func() {
		err := zipFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка при закрытии zip-файла: %v\n", err)
		}
	}()
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		err := zipWriter.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка при закрытии zipWriter: %v\n", err)
		}
	}()
	for _, src := range sources {
		err := addFileToZip(zipWriter, src, "")
		if err != nil {
			return fmt.Errorf("ошибка при добавлении %s: %w", src, err)
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, src, baseInZip string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			entryPath := filepath.Join(src, entry.Name())
			var entryBase string
			if baseInZip == "" {
				entryBase = entry.Name()
			} else {
				entryBase = filepath.Join(baseInZip, entry.Name())
			}
			err = addFileToZip(zipWriter, entryPath, entryBase)
			if err != nil {
				return err
			}
		}
		return nil
	}
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка при закрытии файла: %v\n", err)
		}
	}()
	zipHeader, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// baseInZip всегда относительный путь без ведущих слэшей
	nameInZip := baseInZip
	if nameInZip == "" {
		nameInZip = filepath.Base(src)
	}
	zipHeader.Name = filepath.ToSlash(nameInZip)
	zipHeader.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(zipHeader)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

// ExtractArchive безопасно распаковывает zip-архив (path traversal mitigation)
func (a *Archiver) ExtractArchive(source, destination string) error {
	if filepath.Ext(source) != ".zip" {
		return fmt.Errorf("поддерживается только распаковка zip-архивов")
	}
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("не удалось открыть архив: %w", err)
	}
	defer func() {
		err := zipReader.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка при закрытии zipReader: %v\n", err)
		}
	}()
	for _, f := range zipReader.File {
		if strings.Contains(f.Name, "..") || filepath.IsAbs(f.Name) {
			return fmt.Errorf("архив содержит небезопасный путь: %s", f.Name)
		}
		fpath := filepath.Join(destination, f.Name)
		if !strings.HasPrefix(filepath.Clean(fpath)+string(os.PathSeparator), filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("path traversal: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, f.Mode()); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			errClose := outFile.Close()
			if errClose != nil {
				fmt.Fprintf(os.Stderr, "ошибка при закрытии outFile: %v\n", errClose)
			}
			return err
		}
		_, err = io.Copy(outFile, rc)
		errClose1 := outFile.Close()
		errClose2 := rc.Close()
		if errClose1 != nil {
			fmt.Fprintf(os.Stderr, "ошибка при закрытии outFile: %v\n", errClose1)
		}
		if errClose2 != nil {
			fmt.Fprintf(os.Stderr, "ошибка при закрытии rc: %v\n", errClose2)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// ListArchiveContents выводит содержимое zip-архива
func (a *Archiver) ListArchiveContents(source string) ([]string, error) {
	if filepath.Ext(source) != ".zip" {
		return nil, fmt.Errorf("поддерживается только просмотр zip-архивов")
	}
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть архив: %w", err)
	}
	var files []string
	for _, f := range zipReader.File {
		files = append(files, f.Name)
	}
	return files, nil
}
