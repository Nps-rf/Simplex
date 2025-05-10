package fileops

import (
	"bufio"
	"file-manager/internal/i18n"
	"fmt"
	"io"
	"os"
	"strings"
)

// FileViewer предоставляет функции для просмотра содержимого файлов
type FileViewer struct {
	MaxLineLength int // Максимальная длина строки для отображения
}

// NewFileViewer создает новый экземпляр FileViewer
func NewFileViewer() *FileViewer {
	return &FileViewer{
		MaxLineLength: 100, // По умолчанию ограничиваем длину строки
	}
}

// ViewTextFile показывает содержимое текстового файла
func (v *FileViewer) ViewTextFile(path string, startLine, maxLines int) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("viewer_open"), path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf(i18n.T("viewer_close_error"), path, err))
		}
	}()

	if isBinaryFile(file) {
		return nil, fmt.Errorf(i18n.T("viewer_binary_error"))
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf(i18n.T("viewer_seek_error"), err)
	}

	// Читаем файл построчно
	scanner := bufio.NewScanner(file)
	var lines []string
	lineNum := 0

	// Пропускаем строки до startLine
	for lineNum < startLine && scanner.Scan() {
		lineNum++
	}

	// Читаем maxLines строк или до конца файла
	for scanner.Scan() && (maxLines <= 0 || len(lines) < maxLines) {
		line := scanner.Text()

		// Обрезаем слишком длинные строки
		if len(line) > v.MaxLineLength {
			line = line[:v.MaxLineLength] + "..."
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(i18n.T("viewer_read_error"), err)
	}

	return lines, nil
}

// GetTotalLines возвращает общее количество строк в файле
func (v *FileViewer) GetTotalLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf(i18n.T("viewer_open"), path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf(i18n.T("viewer_close_error"), path, err))
		}
	}()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf(i18n.T("viewer_count_error"), err)
	}

	return lineCount, nil
}

// isBinaryFile проверяет, является ли файл бинарным
func isBinaryFile(file io.Reader) bool {
	// Читаем первые 512 байт для определения
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return true // В случае ошибки считаем файл бинарным
	}
	buf = buf[:n]

	// Подсчитываем количество нулевых байт и непечатаемых символов
	zeroCount := 0
	nonPrintable := 0
	for _, b := range buf {
		if b == 0 {
			zeroCount++
		} else if b < 32 && b != 9 && b != 10 && b != 13 { // Исключаем табуляцию, новую строку и возврат каретки
			nonPrintable++
		}
	}

	// Если доля нулевых байт или непечатаемых символов превышает порог, считаем файл бинарным
	threshold := 0.1 // 10%
	if n > 0 && (float64(zeroCount)/float64(n) > threshold || float64(nonPrintable)/float64(n) > threshold) {
		return true
	}

	return false
}

// FormatTextContent форматирует содержимое текста для отображения с номерами строк
func (v *FileViewer) FormatTextContent(lines []string, startLine int) string {
	var sb strings.Builder

	for i, line := range lines {
		sb.WriteString(fmt.Sprintf("%5d | %s\n", startLine+i, line))
	}

	return sb.String()
}
