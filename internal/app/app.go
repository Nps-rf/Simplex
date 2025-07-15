// Package app реализует основную логику командного интерфейса файлового менеджера.
package app

import (
	"fmt"

	"file-manager/internal/display"
	"file-manager/internal/fileops"
	"file-manager/internal/logger"
	"file-manager/internal/navigation"
	"file-manager/internal/search"
)

// Command представляет команду файлового менеджера
type Command struct {
	Name        string
	Description string
	Execute     func(args []string) error
}

// App представляет основное приложение файлового менеджера
type App struct {
	navigator          *navigation.Navigator
	fileOperator       *fileops.FileOperator
	searcher           *search.Searcher
	display            *display.Display
	fileViewer         *fileops.FileViewer
	archiver           *fileops.Archiver
	permissionsManager *fileops.PermissionsManager
	bookmarkManager    *navigation.BookmarkManager
	logger             *logger.Logger
	commands           map[string]Command
	isRunning          bool
	filterOptions      *navigation.FilterOptions
}

// NewApp создает новый экземпляр App
func NewApp() (*App, error) {
	navigator, err := navigation.NewNavigator()
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать навигатор: %w", err)
	}

	bookmarkManager, err := navigation.NewBookmarkManager()
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать менеджер закладок: %w", err)
	}

	log, err := logger.NewLogger()
	if err != nil {
		return nil, fmt.Errorf("не удалось инициализировать журнал: %w", err)
	}

	app := &App{
		navigator:          navigator,
		fileOperator:       fileops.NewFileOperator(),
		searcher:           search.NewSearcher(),
		display:            display.NewDisplay(),
		fileViewer:         fileops.NewFileViewer(),
		archiver:           fileops.NewArchiver(),
		permissionsManager: fileops.NewPermissionsManager(),
		bookmarkManager:    bookmarkManager,
		logger:             log,
		commands:           make(map[string]Command),
		isRunning:          false,
		filterOptions:      navigation.NewFilterOptions(),
	}

	app.registerCommands()
	return app, nil
}
