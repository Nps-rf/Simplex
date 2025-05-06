package fileops

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestFileOperator(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "fileops_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("ошибка при удалении временной директории: %v", err)
		}
	}()

	// Инициализируем FileOperator
	fileOperator := NewFileOperator()

	// Тест на создание файла
	t.Run("CreateFile", func(t *testing.T) {
		testFilePath := filepath.Join(tempDir, "test_create.txt")

		err := fileOperator.CreateFile(testFilePath)
		if err != nil {
			t.Errorf("не удалось создать файл: %v", err)
		}

		// Проверяем, что файл создан
		if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
			t.Error("файл не был создан")
		}
	})

	// Тест на создание директории
	t.Run("CreateDirectory", func(t *testing.T) {
		testDirPath := filepath.Join(tempDir, "test_dir")

		err := fileOperator.CreateDirectory(testDirPath)
		if err != nil {
			t.Errorf("не удалось создать директорию: %v", err)
		}

		// Проверяем, что директория создана
		fileInfo, err := os.Stat(testDirPath)
		if os.IsNotExist(err) {
			t.Error("директория не была создана")
		}
		if !fileInfo.IsDir() {
			t.Error("созданный объект не является директорией")
		}
	})

	// Тест на копирование файла
	t.Run("CopyFile", func(t *testing.T) {
		// Создаем исходный файл
		sourceFile := filepath.Join(tempDir, "source.txt")
		err := os.WriteFile(sourceFile, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("не удалось создать исходный файл: %v", err)
		}

		// Копируем файл
		destFile := filepath.Join(tempDir, "dest.txt")
		err = fileOperator.CopyFile(sourceFile, destFile)
		if err != nil {
			t.Errorf("не удалось скопировать файл: %v", err)
		}

		// Проверяем, что файл скопирован
		if _, err := os.Stat(destFile); os.IsNotExist(err) {
			t.Error("файл назначения не был создан")
		}

		// Проверяем содержимое скопированного файла
		sourceContent, err := os.ReadFile(sourceFile)
		if err != nil {
			t.Fatalf("не удалось прочитать исходный файл: %v", err)
		}
		destContent, err := os.ReadFile(destFile)
		if err != nil {
			t.Fatalf("не удалось прочитать файл назначения: %v", err)
		}

		if string(sourceContent) != string(destContent) {
			t.Error("содержимое скопированного файла не соответствует исходному")
		}
	})

	// Тест на копирование директории
	t.Run("CopyDirectory", func(t *testing.T) {
		// Создаем исходную директорию и файл в ней
		sourceDir := filepath.Join(tempDir, "source_dir")
		err := os.Mkdir(sourceDir, 0755)
		if err != nil {
			t.Fatalf("не удалось создать исходную директорию: %v", err)
		}

		sourceFile := filepath.Join(sourceDir, "file.txt")
		err = os.WriteFile(sourceFile, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("не удалось создать файл в исходной директории: %v", err)
		}

		// Копируем директорию
		destDir := filepath.Join(tempDir, "dest_dir")
		err = fileOperator.CopyDirectory(sourceDir, destDir)
		if err != nil {
			t.Errorf("не удалось скопировать директорию: %v", err)
		}

		// Проверяем, что директория скопирована
		if _, err := os.Stat(destDir); os.IsNotExist(err) {
			t.Error("директория назначения не была создана")
		}

		// Проверяем, что файл в директории скопирован
		destFile := filepath.Join(destDir, "file.txt")
		if _, err := os.Stat(destFile); os.IsNotExist(err) {
			t.Error("файл в директории назначения не был создан")
		}
	})

	// Тест на перемещение файла
	t.Run("MoveFile", func(t *testing.T) {
		// Создаем исходный файл
		sourceFile := filepath.Join(tempDir, "move_source.txt")
		err := os.WriteFile(sourceFile, []byte("test content for move"), 0644)
		if err != nil {
			t.Fatalf("не удалось создать исходный файл: %v", err)
		}

		// Перемещаем файл
		destFile := filepath.Join(tempDir, "move_dest.txt")
		err = fileOperator.MoveFile(sourceFile, destFile)
		if err != nil {
			t.Errorf("не удалось переместить файл: %v", err)
		}

		// Проверяем, что исходный файл больше не существует
		if _, err := os.Stat(sourceFile); !os.IsNotExist(err) {
			t.Error("исходный файл все еще существует после перемещения")
		}

		// Проверяем, что файл назначения создан
		if _, err := os.Stat(destFile); os.IsNotExist(err) {
			t.Error("файл назначения не был создан")
		}

		// Проверяем содержимое перемещенного файла
		destContent, err := os.ReadFile(destFile)
		if err != nil {
			t.Fatalf("не удалось прочитать файл назначения: %v", err)
		}

		if string(destContent) != "test content for move" {
			t.Error("содержимое перемещенного файла не соответствует исходному")
		}
	})

	// Тест на удаление файла
	t.Run("DeleteFile", func(t *testing.T) {
		// Создаем файл для удаления
		fileToDelete := filepath.Join(tempDir, "to_delete.txt")
		err := os.WriteFile(fileToDelete, []byte("delete me"), 0644)
		if err != nil {
			t.Fatalf("не удалось создать файл для удаления: %v", err)
		}

		// Удаляем файл
		err = fileOperator.DeleteFile(fileToDelete)
		if err != nil {
			t.Errorf("не удалось удалить файл: %v", err)
		}

		// Проверяем, что файл удален
		if _, err := os.Stat(fileToDelete); !os.IsNotExist(err) {
			t.Error("файл не был удален")
		}
	})

	// Тест на удаление директории
	t.Run("DeleteDirectory", func(t *testing.T) {
		// Создаем директорию для удаления
		dirToDelete := filepath.Join(tempDir, "dir_to_delete")
		err := os.Mkdir(dirToDelete, 0755)
		if err != nil {
			t.Fatalf("не удалось создать директорию для удаления: %v", err)
		}

		// Создаем файл в директории
		fileInDir := filepath.Join(dirToDelete, "file.txt")
		err = os.WriteFile(fileInDir, []byte("delete me too"), 0644)
		if err != nil {
			t.Fatalf("не удалось создать файл в директории: %v", err)
		}

		// Удаляем директорию
		err = fileOperator.DeleteDirectory(dirToDelete)
		if err != nil {
			t.Errorf("не удалось удалить директорию: %v", err)
		}

		// Проверяем, что директория удалена
		if _, err := os.Stat(dirToDelete); !os.IsNotExist(err) {
			t.Error("директория не была удалена")
		}
	})
}

func TestFileViewer(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "fileviewer_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("ошибка при удалении временной директории: %v", err)
		}
	}()

	// Создаем тестовый текстовый файл
	testFile := filepath.Join(tempDir, "test.txt")
	content := "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\n"
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("не удалось создать тестовый файл: %v", err)
	}

	// Инициализируем FileViewer
	fileViewer := NewFileViewer()

	// Тест на просмотр всего файла
	t.Run("ViewTextFile_All", func(t *testing.T) {
		lines, err := fileViewer.ViewTextFile(testFile, 0, 10)
		if err != nil {
			t.Errorf("не удалось просмотреть файл: %v", err)
		}

		expectedLines := 5
		if len(lines) != expectedLines {
			t.Errorf("неверное количество строк: получено %d, ожидалось %d", len(lines), expectedLines)
		}

		if lines[0] != "Line 1" {
			t.Errorf("неверное содержимое первой строки: получено %s, ожидалось %s", lines[0], "Line 1")
		}
	})

	// Тест на просмотр части файла с указанной строки
	t.Run("ViewTextFile_PartialOffset", func(t *testing.T) {
		lines, err := fileViewer.ViewTextFile(testFile, 2, 2)
		if err != nil {
			t.Errorf("не удалось просмотреть файл: %v", err)
		}

		expectedLines := 2
		if len(lines) != expectedLines {
			t.Errorf("неверное количество строк: получено %d, ожидалось %d", len(lines), expectedLines)
		}

		if lines[0] != "Line 3" {
			t.Errorf("неверное содержимое первой строки: получено %s, ожидалось %s", lines[0], "Line 3")
		}
	})

	// Тест на подсчет строк
	t.Run("GetTotalLines", func(t *testing.T) {
		totalLines, err := fileViewer.GetTotalLines(testFile)
		if err != nil {
			t.Errorf("не удалось подсчитать строки: %v", err)
		}

		expectedLines := 5
		if totalLines != expectedLines {
			t.Errorf("неверное количество строк: получено %d, ожидалось %d", totalLines, expectedLines)
		}
	})

	// Тест на обработку бинарного файла
	t.Run("ViewTextFile_Binary", func(t *testing.T) {
		// Создаем бинарный файл
		binaryFile := filepath.Join(tempDir, "binary.dat")
		binaryData := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0, 0, 0}
		err := os.WriteFile(binaryFile, binaryData, 0644)
		if err != nil {
			t.Fatalf("не удалось создать бинарный файл: %v", err)
		}

		_, err = fileViewer.ViewTextFile(binaryFile, 0, 10)
		// Должна быть ошибка при попытке просмотреть бинарный файл как текст
		if err == nil {
			t.Error("ожидалась ошибка при просмотре бинарного файла как текста")
		}
	})

	// Тест на форматирование текстового содержимого
	t.Run("FormatTextContent", func(t *testing.T) {
		lines := []string{"Line 1", "Line 2"}
		formatted := fileViewer.FormatTextContent(lines, 0)

		// В сформатированном тексте должны быть номера строк и сами строки
		expectedPart1 := "    0 | Line 1"
		expectedPart2 := "    1 | Line 2"

		if !strings.Contains(formatted, expectedPart1) || !strings.Contains(formatted, expectedPart2) {
			t.Error("неверное форматирование текстового содержимого")
		}
	})
}

func TestArchiver(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "archiver_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("ошибка при удалении временной директории: %v", err)
		}
	}()

	// Создаем тестовые файлы и директории для архивации
	testFiles := []struct {
		path    string
		content string
		isDir   bool
	}{
		{filepath.Join(tempDir, "file1.txt"), "Content of file 1", false},
		{filepath.Join(tempDir, "file2.txt"), "Content of file 2", false},
		{filepath.Join(tempDir, "subdir"), "", true},
		{filepath.Join(tempDir, "subdir", "file3.txt"), "Content of file 3", false},
	}

	for _, tf := range testFiles {
		if tf.isDir {
			if err := os.Mkdir(tf.path, 0755); err != nil {
				t.Fatalf("не удалось создать директорию %s: %v", tf.path, err)
			}
		} else {
			if err := os.WriteFile(tf.path, []byte(tf.content), 0644); err != nil {
				t.Fatalf("не удалось создать файл %s: %v", tf.path, err)
			}
		}
	}

	// Инициализируем Archiver
	archiver := NewArchiver()

	// Тест на создание zip-архива
	t.Run("ArchiveFiles_Zip", func(t *testing.T) {
		// Пути для архивации
		sources := []string{
			filepath.Join(tempDir, "file1.txt"),
			filepath.Join(tempDir, "file2.txt"),
		}

		// Путь архива
		zipFile := filepath.Join(tempDir, "archive.zip")

		// Архивируем файлы
		err := archiver.ArchiveFiles(sources, zipFile, "zip")
		if err != nil {
			t.Errorf("не удалось создать zip-архив: %v", err)
		}

		// Проверяем, что архив создан
		if _, err := os.Stat(zipFile); os.IsNotExist(err) {
			t.Error("zip-архив не был создан")
		}

		// Проверяем, что исходные файлы существуют
		for _, src := range sources {
			if _, err := os.Stat(src); err != nil {
				t.Fatalf("Исходный файл %s не существует: %v", src, err)
			}
		}
	})

	// Тест на извлечение содержимого архива
	t.Run("ListArchiveContents", func(t *testing.T) {
		zipFile := filepath.Join(tempDir, "archive.zip")

		// Получаем список содержимого архива
		contents, err := archiver.ListArchiveContents(zipFile)
		if err != nil {
			t.Errorf("не удалось получить содержимое архива: %v", err)
		}

		// Проверяем, что в архиве есть только два файла верхнего уровня
		if len(contents) != 2 {
			t.Errorf("неверное количество элементов в архиве: получено %d, ожидалось 2", len(contents))
		}
		// Проверяем наличие только file1.txt и file2.txt
		found1 := false
		found2 := false
		for _, item := range contents {
			if item == "file1.txt" {
				found1 = true
			}
			if item == "file2.txt" {
				found2 = true
			}
		}
		if !found1 || !found2 {
			t.Error("не все файлы найдены в содержимом архива (ожидаются только file1.txt и file2.txt)")
		}

		// ВРЕМЕННО: выводим содержимое архива для отладки
		t.Logf("Содержимое архива: %v", contents)
	})

	// Тест на распаковку архива
	t.Run("ExtractArchive", func(t *testing.T) {
		zipFile := filepath.Join(tempDir, "archive.zip")
		extractDir := filepath.Join(tempDir, "unzip_dir")

		// Распаковываем архив
		err := archiver.ExtractArchive(zipFile, extractDir)
		if err != nil {
			t.Errorf("не удалось распаковать архив: %v", err)
		}

		// Проверяем, что директория для распаковки создана
		if _, err := os.Stat(extractDir); os.IsNotExist(err) {
			t.Error("директория для распаковки не была создана")
		}

		// Проверяем, что распакованы только file1.txt и file2.txt
		extractedFile1 := filepath.Join(extractDir, "file1.txt")
		extractedFile2 := filepath.Join(extractDir, "file2.txt")

		if _, err := os.Stat(extractedFile1); os.IsNotExist(err) {
			t.Error("файл file1.txt не был распакован")
		}
		if _, err := os.Stat(extractedFile2); os.IsNotExist(err) {
			t.Error("файл file2.txt не был распакован")
		}

		// Проверяем содержимое распакованных файлов
		content1, err := os.ReadFile(extractedFile1)
		if err != nil {
			t.Fatalf("не удалось прочитать распакованный файл: %v", err)
		}

		if string(content1) != "Content of file 1" {
			t.Error("содержимое распакованного файла не соответствует исходному")
		}
	})
}

func TestPermissionsManager(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "permissions_test")
	if err != nil {
		t.Fatalf("не удалось создать временную директорию: %v", err)
	}
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Errorf("ошибка при удалении временной директории: %v", err)
		}
	}()

	// Создаем тестовый файл
	testFile := filepath.Join(tempDir, "test_permissions.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("не удалось создать тестовый файл: %v", err)
	}

	// Инициализируем PermissionsManager
	permissionsManager := NewPermissionsManager()

	// Тест на изменение прав доступа
	t.Run("ChangePermissions", func(t *testing.T) {
		// Изменяем права на "только чтение"
		err := permissionsManager.ChangePermissions(testFile, "444")
		if err != nil {
			t.Errorf("не удалось изменить права доступа: %v", err)
		}

		// Проверяем, что права изменились
		fileInfo, err := os.Stat(testFile)
		if err != nil {
			t.Fatalf("не удалось получить информацию о файле: %v", err)
		}

		// Сравниваем права доступа (в восьмеричной системе)
		expectedMode := os.FileMode(0444)
		if fileInfo.Mode().Perm() != expectedMode {
			t.Errorf("неверные права доступа: получено %o, ожидалось %o", fileInfo.Mode().Perm(), expectedMode)
		}
	})

	// Тест на получение прав доступа
	t.Run("GetPermissions", func(t *testing.T) {
		// Сначала устанавливаем известные права доступа
		err := os.Chmod(testFile, 0644)
		if err != nil {
			t.Fatalf("не удалось установить права доступа: %v", err)
		}

		// Получаем права доступа через менеджер
		permissions, err := permissionsManager.GetPermissions(testFile)
		if err != nil {
			t.Errorf("не удалось получить права доступа: %v", err)
		}

		// Проверяем, что права не пустые и соответствуют формату восьмеричного числа
		if len(permissions) == 0 {
			t.Error("права доступа не должны быть пустыми")
		}

		// Проверяем, что права можно преобразовать в восьмеричное число
		_, err = strconv.ParseUint(permissions, 8, 32)
		if err != nil {
			t.Errorf("права доступа %s не являются корректным восьмеричным числом: %v", permissions, err)
		}
	})

	// Тест на форматирование прав доступа
	t.Run("FormatPermissions", func(t *testing.T) {
		// Получаем информацию о файле
		fileInfo, err := os.Stat(testFile)
		if err != nil {
			t.Fatalf("не удалось получить информацию о файле: %v", err)
		}

		// Форматируем права доступа
		formatted := permissionsManager.FormatPermissions(fileInfo.Mode())

		// Проверяем, что формат соответствует Unix-подобному (например, "-rw-r--r--")
		if len(formatted) != 10 || formatted[0] != '-' {
			t.Errorf("неверный формат прав доступа: %s", formatted)
		}
	})
}
