package search

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestEnvironment(t *testing.T) (string, func()) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "search-tests")
	if err != nil {
		t.Fatalf("Не удалось создать временную директорию: %v", err)
	}

	// Создаем тестовую структуру файлов
	testFiles := map[string]string{
		"file1.txt":         "Это тестовый файл с текстом для поиска",
		"file2.log":         "Другой тестовый файл с текстом",
		"subdir/file3.txt":  "Файл в поддиректории с текстом для поиска",
		"subdir/file4.conf": "Конфигурационный файл",
		"subdir2/file5.txt": "Еще один текстовый файл",
		"subdir2/file6.bin": "Бинарный файл с специальным содержимым",
		"file7_special.txt": "Специальный файл для тестирования регулярных выражений 123-456-789",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)

		// Создаем поддиректории при необходимости
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Не удалось создать директорию %s: %v", dir, err)
		}

		// Создаем файл и записываем содержимое
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Не удалось создать тестовый файл %s: %v", fullPath, err)
		}
	}

	// Возвращаем функцию очистки
	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

func TestSearchByName(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	searcher := NewSearcher()

	tests := []struct {
		name            string
		pattern         string
		expectedMatches int
	}{
		{"ПоискВсехТекстовыхФайлов", "*.txt", 4},
		{"ПоискЛогФайлов", "*.log", 1},
		{"ПоискКонфигурационныхФайлов", "*.conf", 1},
		{"ПоискНесуществующихФайлов", "*.xyz", 0},
		{"ПоискПоЧастиИмени", "file[1-3].*", 3},
		{"ПоискСпециальныхФайлов", "*special*", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := searcher.SearchByName(tempDir, tt.pattern)
			if err != nil {
				t.Fatalf("SearchByName вернул ошибку: %v", err)
			}

			if len(results) != tt.expectedMatches {
				t.Errorf("SearchByName для шаблона %s нашел %d совпадений, ожидалось %d",
					tt.pattern, len(results), tt.expectedMatches)
				t.Logf("Найденные файлы: %v", results)
			}
		})
	}
}

func TestSearchByContent(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	searcher := NewSearcher()

	tests := []struct {
		name            string
		content         string
		expectedMatches int
	}{
		{"ПоискПоФразеДляПоиска", "для поиска", 2},
		{"ПоискПоСловуТестовый", "тестовый", 2},
		{"ПоискПоСловуКонфигурационный", "Конфигурационный", 1},
		{"ПоискНесуществующегоСодержимого", "этого текста нет нигде", 0},
		{"ПоискСпециальногоСодержимого", "специальным содержимым", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := searcher.SearchByContent(tempDir, tt.content)
			if err != nil {
				t.Fatalf("SearchByContent вернул ошибку: %v", err)
			}

			if len(results) != tt.expectedMatches {
				t.Errorf("SearchByContent для текста '%s' нашел %d совпадений, ожидалось %d",
					tt.content, len(results), tt.expectedMatches)
				t.Logf("Найденные файлы: %v", results)
			}
		})
	}
}

func TestSearchByRegex(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	searcher := NewSearcher()

	tests := []struct {
		name            string
		regexPattern    string
		expectedMatches int
	}{
		{"ПоискПоРегулярномуВыражениюЦифры", `\d{3}-\d{3}-\d{3}`, 1},
		{"ПоискПоРегулярномуВыражениюСлова", `файл\s+с\s+текстом`, 2},
		{"ПоискПоРегулярномуВыражениюНачалоСтроки", `^Это`, 1},
		{"ПоискНесуществующегоШаблона", `xyz\d{10}`, 0},
		{"ПоискСловаСпециальный", `[Сс]пециальн(ый|ое|ым)`, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := searcher.SearchByRegex(tempDir, tt.regexPattern)
			if err != nil {
				t.Fatalf("SearchByRegex вернул ошибку: %v", err)
			}

			if len(results) != tt.expectedMatches {
				t.Errorf("SearchByRegex для шаблона '%s' нашел %d совпадений, ожидалось %d",
					tt.regexPattern, len(results), tt.expectedMatches)
				t.Logf("Найденные файлы: %v", results)
			}
		})
	}
}

func TestInvalidRegex(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	searcher := NewSearcher()

	// Проверка обработки некорректного регулярного выражения
	_, err := searcher.SearchByRegex(tempDir, "[неправильное выражение")
	if err == nil {
		t.Error("SearchByRegex должен вернуть ошибку для некорректного регулярного выражения")
	}
}
