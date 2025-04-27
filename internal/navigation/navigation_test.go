package navigation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNavigator(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "navigation_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем вложенные директории для тестирования
	nestedDir := filepath.Join(tempDir, "nested")
	if err := os.Mkdir(nestedDir, 0755); err != nil {
		t.Fatalf("не удалось создать вложенную директорию: %v", err)
	}

	// Инициализируем навигатор
	navigator, err := NewNavigator()
	if err != nil {
		t.Fatalf("не удалось создать навигатор: %v", err)
	}

	// Тест на получение текущей директории
	t.Run("GetCurrentDirectory", func(t *testing.T) {
		dir := navigator.GetCurrentDirectory()
		if dir == "" {
			t.Error("текущая директория не должна быть пустой")
		}
	})

	// Тест на изменение директории (абсолютный путь)
	t.Run("ChangeDirectory_Absolute", func(t *testing.T) {
		err := navigator.ChangeDirectory(tempDir)
		if err != nil {
			t.Errorf("не удалось изменить директорию на %s: %v", tempDir, err)
		}

		currentDir := navigator.GetCurrentDirectory()
		if currentDir != tempDir {
			t.Errorf("текущая директория (%s) не соответствует ожидаемой (%s)", currentDir, tempDir)
		}
	})

	// Тест на изменение директории (относительный путь)
	t.Run("ChangeDirectory_Relative", func(t *testing.T) {
		// Сначала переходим во временную директорию
		err := navigator.ChangeDirectory(tempDir)
		if err != nil {
			t.Errorf("не удалось изменить директорию на %s: %v", tempDir, err)
		}

		// Теперь переходим в nested используя относительный путь
		err = navigator.ChangeDirectory("nested")
		if err != nil {
			t.Errorf("не удалось изменить директорию на nested: %v", err)
		}

		currentDir := navigator.GetCurrentDirectory()
		if currentDir != nestedDir {
			t.Errorf("текущая директория (%s) не соответствует ожидаемой (%s)", currentDir, nestedDir)
		}
	})

	// Тест на переход в родительскую директорию
	t.Run("ChangeDirectory_Parent", func(t *testing.T) {
		// Сначала переходим во вложенную директорию
		err := navigator.ChangeDirectory(nestedDir)
		if err != nil {
			t.Errorf("не удалось изменить директорию на %s: %v", nestedDir, err)
		}

		// Теперь переходим в родительскую директорию
		err = navigator.ChangeDirectory("..")
		if err != nil {
			t.Errorf("не удалось изменить директорию на родительскую: %v", err)
		}

		currentDir := navigator.GetCurrentDirectory()
		if currentDir != tempDir {
			t.Errorf("текущая директория (%s) не соответствует ожидаемой (%s)", currentDir, tempDir)
		}
	})

	// Тест на получение списка содержимого директории
	t.Run("ListDirectory", func(t *testing.T) {
		// Создаем тестовый файл
		testFilePath := filepath.Join(tempDir, "test.txt")
		testFile, err := os.Create(testFilePath)
		if err != nil {
			t.Fatalf("не удалось создать тестовый файл: %v", err)
		}
		testFile.Close()

		// Переходим во временную директорию
		err = navigator.ChangeDirectory(tempDir)
		if err != nil {
			t.Errorf("не удалось изменить директорию на %s: %v", tempDir, err)
		}

		// Получаем список содержимого
		entries, err := navigator.ListDirectory()
		if err != nil {
			t.Errorf("не удалось получить содержимое директории: %v", err)
		}

		// Проверяем, что в списке есть директория nested и файл test.txt
		foundNested := false
		foundTestFile := false
		for _, entry := range entries {
			if entry.Name() == "nested" && entry.IsDir() {
				foundNested = true
			}
			if entry.Name() == "test.txt" && !entry.IsDir() {
				foundTestFile = true
			}
		}

		if !foundNested {
			t.Error("директория 'nested' не найдена в списке содержимого")
		}
		if !foundTestFile {
			t.Error("файл 'test.txt' не найден в списке содержимого")
		}
	})

	// Тест на обработку ошибок при изменении директории
	t.Run("ChangeDirectory_Error", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "non_existent")
		err := navigator.ChangeDirectory(nonExistentDir)
		if err == nil {
			t.Error("ожидалась ошибка при переходе в несуществующую директорию")
		}
	})
}

func TestBookmarkManager(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "bookmark_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем вложенную директорию
	nestedDir := filepath.Join(tempDir, "nested")
	if err := os.Mkdir(nestedDir, 0755); err != nil {
		t.Fatalf("не удалось создать вложенную директорию: %v", err)
	}

	// Инициализируем менеджер закладок
	bookmarkManager, err := NewBookmarkManager()
	if err != nil {
		t.Fatalf("не удалось создать менеджер закладок: %v", err)
	}

	// Тест на добавление закладки
	t.Run("AddBookmark", func(t *testing.T) {
		err := bookmarkManager.AddBookmark("temp", tempDir)
		if err != nil {
			t.Errorf("не удалось добавить закладку: %v", err)
		}

		// Проверяем, что закладка добавлена
		bookmarks := bookmarkManager.ListBookmarks()
		found := false
		for _, bookmark := range bookmarks {
			if bookmark.Name == "temp" && bookmark.Path == tempDir {
				found = true
				break
			}
		}
		if !found {
			t.Error("добавленная закладка не найдена в списке")
		}
	})

	// Тест на получение пути по имени закладки
	t.Run("GetBookmarkPath", func(t *testing.T) {
		// Добавляем закладку, если её ещё нет
		if err := bookmarkManager.AddBookmark("nested", nestedDir); err != nil {
			t.Fatalf("не удалось добавить закладку: %v", err)
		}

		path, err := bookmarkManager.GetBookmarkPath("nested")
		if err != nil {
			t.Errorf("не удалось получить путь закладки: %v", err)
		}
		if path != nestedDir {
			t.Errorf("полученный путь (%s) не соответствует ожидаемому (%s)", path, nestedDir)
		}
	})

	// Тест на удаление закладки
	t.Run("RemoveBookmark", func(t *testing.T) {
		// Добавляем временную закладку для удаления
		if err := bookmarkManager.AddBookmark("to_remove", tempDir); err != nil {
			t.Fatalf("не удалось добавить закладку: %v", err)
		}

		// Удаляем закладку
		err := bookmarkManager.RemoveBookmark("to_remove")
		if err != nil {
			t.Errorf("не удалось удалить закладку: %v", err)
		}

		// Проверяем, что закладка удалена
		_, err = bookmarkManager.GetBookmarkPath("to_remove")
		if err == nil {
			t.Error("закладка не была удалена")
		}
	})

	// Тест на получение списка закладок
	t.Run("ListBookmarks", func(t *testing.T) {
		// Очищаем существующие закладки
		bookmarks := bookmarkManager.ListBookmarks()
		for _, bookmark := range bookmarks {
			bookmarkManager.RemoveBookmark(bookmark.Name)
		}

		// Добавляем две закладки
		if err := bookmarkManager.AddBookmark("temp1", tempDir); err != nil {
			t.Fatalf("не удалось добавить закладку: %v", err)
		}
		if err := bookmarkManager.AddBookmark("temp2", nestedDir); err != nil {
			t.Fatalf("не удалось добавить закладку: %v", err)
		}

		// Получаем список закладок
		bookmarks = bookmarkManager.ListBookmarks()
		if len(bookmarks) != 2 {
			t.Errorf("неверное количество закладок: получено %d, ожидалось 2", len(bookmarks))
		}
	})

	// Тест на обработку ошибок при работе с закладками
	t.Run("BookmarkErrors", func(t *testing.T) {
		// Попытка добавить закладку с существующим именем
		if err := bookmarkManager.AddBookmark("temp1", nestedDir); err == nil {
			t.Error("ожидалась ошибка при добавлении закладки с существующим именем")
		}

		// Попытка получить несуществующую закладку
		_, err := bookmarkManager.GetBookmarkPath("non_existent")
		if err == nil {
			t.Error("ожидалась ошибка при получении несуществующей закладки")
		}

		// Попытка удалить несуществующую закладку
		err = bookmarkManager.RemoveBookmark("non_existent")
		if err == nil {
			t.Error("ожидалась ошибка при удалении несуществующей закладки")
		}
	})
}

func TestFilter(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "filter_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Создаем тестовые файлы и директории
	testFiles := []struct {
		name string
		dir  bool
		size int64
	}{
		{"file1.txt", false, 100},
		{"file2.txt", false, 200},
		{"image.jpg", false, 300},
		{"hidden.txt", false, 50},
		{".hidden", false, 30},
		{"subdir", true, 0},
	}

	for _, tf := range testFiles {
		path := filepath.Join(tempDir, tf.name)
		if tf.dir {
			if err := os.Mkdir(path, 0755); err != nil {
				t.Fatalf("не удалось создать директорию %s: %v", path, err)
			}
		} else {
			file, err := os.Create(path)
			if err != nil {
				t.Fatalf("не удалось создать файл %s: %v", path, err)
			}
			defer file.Close()

			// Устанавливаем размер файла
			if err := file.Truncate(tf.size); err != nil {
				t.Fatalf("не удалось установить размер файла %s: %v", path, err)
			}
		}
	}

	// Читаем содержимое директории
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("не удалось прочитать содержимое директории: %v", err)
	}

	// Тест на фильтрацию по расширению
	t.Run("FilterByExtension", func(t *testing.T) {
		options := NewFilterOptions()
		options.Extensions = []string{"txt"}

		filtered, err := Filter(entries, tempDir, options)
		if err != nil {
			t.Errorf("ошибка при фильтрации: %v", err)
		}

		expectedCount := 3 // file1.txt, file2.txt, hidden.txt
		if len(filtered) != expectedCount {
			t.Errorf("неверное количество файлов после фильтрации: получено %d, ожидалось %d", len(filtered), expectedCount)
		}

		// Проверяем, что в результате только файлы .txt
		for _, entry := range filtered {
			if filepath.Ext(entry.Name()) != ".txt" {
				t.Errorf("файл %s не должен быть в результате фильтрации по расширению .txt", entry.Name())
			}
		}
	})

	// Тест на фильтрацию по шаблону имени
	t.Run("FilterByNamePattern", func(t *testing.T) {
		options := NewFilterOptions()
		options.NamePattern = "file*"

		filtered, err := Filter(entries, tempDir, options)
		if err != nil {
			t.Errorf("ошибка при фильтрации: %v", err)
		}

		expectedCount := 2 // file1.txt, file2.txt
		if len(filtered) != expectedCount {
			t.Errorf("неверное количество файлов после фильтрации: получено %d, ожидалось %d", len(filtered), expectedCount)
		}

		// Проверяем, что в результате только файлы, начинающиеся с "file"
		for _, entry := range filtered {
			if !strings.HasPrefix(entry.Name(), "file") {
				t.Errorf("файл %s не должен быть в результате фильтрации по шаблону 'file*'", entry.Name())
			}
		}
	})

	// Тест на фильтрацию по размеру
	t.Run("FilterBySize", func(t *testing.T) {
		options := NewFilterOptions()
		options.MinSize = 100
		options.MaxSize = 200

		filtered, err := Filter(entries, tempDir, options)
		if err != nil {
			t.Errorf("ошибка при фильтрации: %v", err)
		}

		// Выводим для отладки список файлов и их размеры
		t.Logf("Файлы, прошедшие фильтрацию по размеру:")
		for _, entry := range filtered {
			info, _ := entry.Info()
			t.Logf("Файл: %s, Размер: %d", entry.Name(), info.Size())
		}

		expectedCount := 2 // file1.txt (100), file2.txt (200)
		if len(filtered) != expectedCount {
			t.Errorf("неверное количество файлов после фильтрации: получено %d, ожидалось %d", len(filtered), expectedCount)
		}
	})

	// Тест на фильтрацию по типу
	t.Run("FilterByType", func(t *testing.T) {
		options := NewFilterOptions()
		options.ShowDirs = true
		options.ShowFiles = false

		filtered, err := Filter(entries, tempDir, options)
		if err != nil {
			t.Errorf("ошибка при фильтрации: %v", err)
		}

		expectedCount := 1 // subdir
		if len(filtered) != expectedCount {
			t.Errorf("неверное количество директорий после фильтрации: получено %d, ожидалось %d", len(filtered), expectedCount)
		}

		// Проверяем, что в результате только директории
		for _, entry := range filtered {
			if !entry.IsDir() {
				t.Errorf("файл %s не должен быть в результате фильтрации по типу 'директория'", entry.Name())
			}
		}
	})

	// Тест на фильтрацию скрытых файлов
	t.Run("FilterHidden", func(t *testing.T) {
		options := NewFilterOptions()
		options.ShowHidden = true

		filtered, err := Filter(entries, tempDir, options)
		if err != nil {
			t.Errorf("ошибка при фильтрации: %v", err)
		}

		// Проверяем, что скрытые файлы включены в результат
		hasHidden := false
		for _, entry := range filtered {
			if entry.Name() == ".hidden" {
				hasHidden = true
				break
			}
		}
		if !hasHidden {
			t.Error("скрытый файл не включен в результат фильтрации")
		}
	})
}
