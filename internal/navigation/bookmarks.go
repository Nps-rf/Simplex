// Package navigation реализует навигацию и работу с закладками для файлового менеджера.
package navigation

import (
	"encoding/json"
	"file-manager/internal/i18n"
	"fmt"
	"os"
	"path/filepath"
)

// Bookmark представляет одну закладку
type Bookmark struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// BookmarkManager управляет закладками директорий
type BookmarkManager struct {
	Bookmarks     []Bookmark
	BookmarksFile string
}

// NewBookmarkManager создает новый экземпляр BookmarkManager
func NewBookmarkManager() (*BookmarkManager, error) {
	// Определение пути к файлу закладок
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf(i18n.T("bm_home"), err)
	}

	configDir := filepath.Join(homeDir, ".filemanager")
	bookmarksFile := filepath.Join(configDir, "bookmarks.json")

	manager := &BookmarkManager{
		Bookmarks:     []Bookmark{},
		BookmarksFile: bookmarksFile,
	}

	// Создаем директорию конфигурации, если она не существует
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, fmt.Errorf(i18n.T("bm_dir"), err)
		}
	}

	// Загружаем закладки, если файл существует
	if _, err := os.Stat(bookmarksFile); err == nil {
		err = manager.LoadBookmarks()
		if err != nil {
			return nil, fmt.Errorf(i18n.T("bm_load"), err)
		}
	}

	return manager, nil
}

// AddBookmark добавляет новую закладку
func (bm *BookmarkManager) AddBookmark(name, path string) error {
	// Проверяем, существует ли директория
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf(i18n.T("bm_path"), err)
	}

	if !info.IsDir() {
		return fmt.Errorf(i18n.T("bm_dir_not"), path)
	}

	// Проверяем уникальность имени
	for _, bookmark := range bm.Bookmarks {
		if bookmark.Name == name {
			return fmt.Errorf(i18n.T("bm_exists"), name)
		}
	}

	// Добавляем закладку
	bm.Bookmarks = append(bm.Bookmarks, Bookmark{
		Name: name,
		Path: path,
	})

	// Сохраняем изменения
	return bm.SaveBookmarks()
}

// RemoveBookmark удаляет закладку по имени
func (bm *BookmarkManager) RemoveBookmark(name string) error {
	found := false
	newBookmarks := []Bookmark{}

	for _, bookmark := range bm.Bookmarks {
		if bookmark.Name != name {
			newBookmarks = append(newBookmarks, bookmark)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf(i18n.T("bm_not_found"), name)
	}

	bm.Bookmarks = newBookmarks
	return bm.SaveBookmarks()
}

// GetBookmarkPath возвращает путь к закладке по имени
func (bm *BookmarkManager) GetBookmarkPath(name string) (string, error) {
	for _, bookmark := range bm.Bookmarks {
		if bookmark.Name == name {
			return bookmark.Path, nil
		}
	}

	return "", fmt.Errorf(i18n.T("bm_not_found"), name)
}

// ListBookmarks возвращает список всех закладок
func (bm *BookmarkManager) ListBookmarks() []Bookmark {
	return bm.Bookmarks
}

// SaveBookmarks сохраняет закладки в файл
func (bm *BookmarkManager) SaveBookmarks() error {
	// Сериализуем закладки в JSON
	data, err := json.MarshalIndent(bm.Bookmarks, "", "  ")
	if err != nil {
		return fmt.Errorf(i18n.T("bm_marshal"), err)
	}

	// Записываем данные в файл
	err = os.WriteFile(bm.BookmarksFile, data, 0644)
	if err != nil {
		return fmt.Errorf(i18n.T("bm_write"), err)
	}

	return nil
}

// LoadBookmarks загружает закладки из файла
func (bm *BookmarkManager) LoadBookmarks() error {
	// Проверяем, существует ли файл
	if _, err := os.Stat(bm.BookmarksFile); os.IsNotExist(err) {
		return nil // Файл не существует, но это не ошибка
	}

	// Читаем данные из файла
	data, err := os.ReadFile(bm.BookmarksFile)
	if err != nil {
		return fmt.Errorf(i18n.T("bm_read"), err)
	}

	// Десериализуем JSON
	if len(data) > 0 {
		err = json.Unmarshal(data, &bm.Bookmarks)
		if err != nil {
			return fmt.Errorf(i18n.T("bm_unmarshal"), err)
		}
	}

	return nil
}
