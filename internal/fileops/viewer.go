package fileops

import (
	"bufio"
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
	// Открываем файл для чтения
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("ошибка при закрытии файла %s: %w", path, err))
		}
	}()

	// Проверяем, является ли файл бинарным
	if isBinaryFile(file) {
		return nil, fmt.Errorf("невозможно отобразить бинарный файл как текст")
	}

	// Сбрасываем позицию файла в начало после проверки
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("ошибка сброса позиции файла: %w", err)
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
		return nil, fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	return lines, nil
}

// GetTotalLines возвращает общее количество строк в файле
func (v *FileViewer) GetTotalLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("не удалось открыть файл %s: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("ошибка при закрытии файла %s: %w", path, err))
		}
	}()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("ошибка при подсчете строк: %w", err)
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
