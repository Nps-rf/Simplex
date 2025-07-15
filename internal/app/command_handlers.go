package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"file-manager/internal/i18n"
	"file-manager/internal/logger"
	"file-manager/internal/navigation"
)

// Команды файлового менеджера

func (a *App) cmdHelp(_ []string) error {
	fmt.Println(i18n.T("help"))

	categories := map[string][]Command{
		i18n.T("category_navigation"): {},
		i18n.T("category_fileops"):    {},
		i18n.T("category_search"):     {},
		i18n.T("category_archive"):    {},
		i18n.T("category_settings"):   {},
		i18n.T("category_other"):      {},
	}

	for _, cmd := range a.commands {
		switch cmd.Name {
		case "ls", "cd", "pwd", "bookmark":
			categories[i18n.T("category_navigation")] = append(categories[i18n.T("category_navigation")], cmd)
		case "mkdir", "touch", "rm", "rmdir", "cp", "mv", "chmod":
			categories[i18n.T("category_fileops")] = append(categories[i18n.T("category_fileops")], cmd)
		case "find", "grep", "info", "cat":
			categories[i18n.T("category_search")] = append(categories[i18n.T("category_search")], cmd)
		case "archive", "extract", "list-archive":
			categories[i18n.T("category_archive")] = append(categories[i18n.T("category_archive")], cmd)
		case "filter", "colors", "log":
			categories[i18n.T("category_settings")] = append(categories[i18n.T("category_settings")], cmd)
		default:
			categories[i18n.T("category_other")] = append(categories[i18n.T("category_other")], cmd)
		}
	}

	for category, cmds := range categories {
		if len(cmds) > 0 {
			fmt.Printf("\n%s:\n", category)
			for _, cmd := range cmds {
				fmt.Printf("  %-15s - %s\n", cmd.Name, i18n.T(cmd.Name))
			}
		}
	}

	return nil
}

func (a *App) cmdListDir(_ []string) error {
	entries, err := a.navigator.ListDirectory()
	if err != nil {
		return err
	}

	// Применяем фильтр, если он активен
	if a.filterOptions != nil {
		dir, dirErr := a.navigator.GetCurrentDirectory()
		if dirErr != nil {
			return dirErr
		}
		entries, err = navigation.Filter(entries, dir, a.filterOptions)
		if err != nil {
			return fmt.Errorf("ошибка при применении фильтра: %w", err)
		}
	}
	dir, dirErr := a.navigator.GetCurrentDirectory()
	if dirErr != nil {
		return dirErr
	}
	fmt.Printf("Содержимое директории: %s\n\n", dir)
	fmt.Println("ТИП  ИМЯ                           РАЗМЕР     ИЗМЕНЕН")
	fmt.Println("---------------------------------------------------")

	for _, entry := range entries {
		formattedEntry, err := a.display.FormatDirEntry(entry, dir)
		if err != nil {
			return err
		}
		fmt.Println(formattedEntry)
	}

	return nil
}

func (a *App) cmdChangeDir(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	return a.navigator.ChangeDirectory(args[0])
}

func (a *App) cmdPrintWorkingDir(_ []string) error {
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

func (a *App) cmdMakeDir(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, args[0])
	return a.fileOperator.CreateDirectory(path)
}

func (a *App) cmdCreateFile(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, args[0])
	return a.fileOperator.CreateFile(path)
}

func (a *App) cmdRemoveFile(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, args[0])
	return a.fileOperator.DeleteFile(path)
}

func (a *App) cmdRemoveDir(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, args[0])
	return a.fileOperator.DeleteDirectory(path)
}

func (a *App) cmdCopy(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf(i18n.T("args_expected_2"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	sourcePath := filepath.Join(dir, args[0])
	destPath := filepath.Join(dir, args[1])

	info, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return a.fileOperator.CopyDirectory(sourcePath, destPath)
	}

	return a.fileOperator.CopyFile(sourcePath, destPath)
}

func (a *App) cmdMove(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf(i18n.T("args_expected_2"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	sourcePath := filepath.Join(dir, args[0])
	destPath := filepath.Join(dir, args[1])

	return a.fileOperator.MoveFile(sourcePath, destPath)
}

func (a *App) cmdFindByName(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	results, err := a.searcher.SearchByName(dir, args[0])
	if err != nil {
		return err
	}
	fmt.Println(a.display.FormatSearchResults(results, args[0]))
	return nil
}

func (a *App) cmdFindByContent(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	results, err := a.searcher.SearchByContent(dir, args[0])
	if err != nil {
		return err
	}
	fmt.Println(a.display.FormatSearchResults(results, args[0]))
	return nil
}

func (a *App) cmdFileInfo(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, args[0])
	fileInfo, err := a.display.GetFileInfo(path)
	if err != nil {
		return err
	}
	fmt.Println(a.display.FormatFileInfo(fileInfo))
	return nil
}

func (a *App) cmdExit(_ []string) error {
	a.isRunning = false
	return nil
}

func (a *App) cmdViewFile(args []string) error {
	if len(args) < 1 || len(args) > 3 {
		return fmt.Errorf(i18n.T("args_expected_1_to_3"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, args[0])
	startLine := 0
	maxLines := 20
	if len(args) >= 2 {
		startLine, err = strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("некорректный номер начальной строки: %w", err)
		}
	}
	if len(args) >= 3 {
		maxLines, err = strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("некорректное количество строк: %w", err)
		}
	}
	lines, err := a.fileViewer.ViewTextFile(path, startLine, maxLines)
	if err != nil {
		return err
	}
	fmt.Println(a.fileViewer.FormatTextContent(lines, startLine))
	totalLines, err := a.fileViewer.GetTotalLines(path)
	if err == nil {
		fmt.Printf("\nПоказано строк: %d-%d из %d\n", startLine, startLine+len(lines)-1, totalLines)
	}
	return nil
}

func (a *App) cmdChangePermissions(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf(i18n.T("args_expected_2"), len(args))
	}
	mode := args[0]
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, args[1])
	return a.permissionsManager.ChangePermissions(path, mode)
}

func (a *App) cmdCreateArchive(args []string) error {
	if len(args) < 3 {
		return fmt.Errorf(i18n.T("args_expected_min_3"), len(args))
	}
	archiveName := args[0]
	format := args[1]
	sources := []string{}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	for _, src := range args[2:] {
		sources = append(sources, filepath.Join(dir, src))
	}
	destination := filepath.Join(dir, archiveName)
	return a.archiver.ArchiveFiles(sources, destination, format)
}

func (a *App) cmdExtractArchive(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf(i18n.T("args_expected_2"), len(args))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	source := filepath.Join(dir, args[0])
	destination := filepath.Join(dir, args[1])
	return a.archiver.ExtractArchive(source, destination)
}

func (a *App) cmdListArchive(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("error"), fmt.Sprintf(i18n.T("archive_args"), 1, len(args)))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	source := filepath.Join(dir, args[0])
	contents, err := a.archiver.ListArchiveContents(source)
	if err != nil {
		return err
	}
	fmt.Printf(i18n.T("archive_contents")+"\n\n", args[0])
	for i, item := range contents {
		fmt.Printf("%d. %s\n", i+1, item)
	}
	return nil
}

func (a *App) cmdManageBookmarks(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf(i18n.T("error"), i18n.T("bookmark_args"))
	}
	dir, err := a.navigator.GetCurrentDirectory()
	if err != nil {
		return err
	}
	switch args[0] {
	case "add":
		if len(args) < 2 {
			return fmt.Errorf(i18n.T("error"), i18n.T("bookmark_add_args"))
		}
		name := args[1]
		path := dir
		if len(args) >= 3 {
			path = args[2]
			if !filepath.IsAbs(path) {
				path = filepath.Join(dir, path)
			}
		}
		return a.bookmarkManager.AddBookmark(name, path)
	case "list":
		bookmarks := a.bookmarkManager.ListBookmarks()
		fmt.Println(i18n.T("bookmark_list"))
		for i, bookmark := range bookmarks {
			fmt.Printf("%d. %s -> %s\n", i+1, bookmark.Name, bookmark.Path)
		}
		return nil
	case "remove":
		if len(args) < 2 {
			return fmt.Errorf(i18n.T("error"), i18n.T("bookmark_remove_args"))
		}
		return a.bookmarkManager.RemoveBookmark(args[1])
	case "go":
		if len(args) < 2 {
			return fmt.Errorf(i18n.T("error"), i18n.T("bookmark_go_args"))
		}
		path, err := a.bookmarkManager.GetBookmarkPath(args[1])
		if err != nil {
			return err
		}
		return a.navigator.ChangeDirectory(path)
	default:
		return fmt.Errorf(i18n.T("error"), fmt.Sprintf(i18n.T("bookmark_unknown"), args[0]))
	}
}

func (a *App) cmdFilter(args []string) error {
	if len(args) == 0 {
		a.filterOptions = navigation.NewFilterOptions()
		fmt.Println(i18n.T("filter_reset"))
		return nil
	}
	newOptions := navigation.NewFilterOptions()
	for _, arg := range args {
		if strings.HasPrefix(arg, "--ext=") {
			ext := strings.TrimPrefix(arg, "--ext=")
			if ext != "" {
				extensions := strings.Split(ext, ",")
				newOptions.Extensions = extensions
			}
		} else if strings.HasPrefix(arg, "--name=") {
			pattern := strings.TrimPrefix(arg, "--name=")
			newOptions.NamePattern = pattern
		} else if strings.HasPrefix(arg, "--size=") {
			sizeRange := strings.TrimPrefix(arg, "--size=")
			parts := strings.Split(sizeRange, "-")

			if len(parts) == 2 {
				// Минимальный размер
				if parts[0] != "" {
					minSize, err := strconv.ParseInt(parts[0], 10, 64)
					if err != nil {
						return fmt.Errorf("некорректный минимальный размер: %w", err)
					}
					newOptions.MinSize = minSize
				}

				// Максимальный размер
				if parts[1] != "" {
					maxSize, err := strconv.ParseInt(parts[1], 10, 64)
					if err != nil {
						return fmt.Errorf("некорректный максимальный размер: %w", err)
					}
					newOptions.MaxSize = maxSize
				}
			}
		} else if strings.HasPrefix(arg, "--date=") {
			dateRange := strings.TrimPrefix(arg, "--date=")
			parts := strings.Split(dateRange, "-")

			if len(parts) == 2 {
				// Дата начала
				if parts[0] != "" {
					startDate, err := time.Parse("2006-01-02", parts[0])
					if err != nil {
						return fmt.Errorf("некорректная дата начала (формат YYYY-MM-DD): %w", err)
					}
					newOptions.ModifiedAfter = startDate
				}

				// Дата окончания
				if parts[1] != "" {
					endDate, err := time.Parse("2006-01-02", parts[1])
					if err != nil {
						return fmt.Errorf("некорректная дата окончания (формат YYYY-MM-DD): %w", err)
					}
					// Устанавливаем конец дня
					endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
					newOptions.ModifiedBefore = endDate
				}
			}
		} else if strings.HasPrefix(arg, "--type=") {
			types := strings.TrimPrefix(arg, "--type=")

			// По умолчанию ничего не показываем
			newOptions.ShowDirs = false
			newOptions.ShowFiles = false
			newOptions.ShowHidden = false

			for _, t := range types {
				switch t {
				case 'd':
					newOptions.ShowDirs = true
				case 'f':
					newOptions.ShowFiles = true
				case 'h':
					newOptions.ShowHidden = true
				}
			}
		}
	}
	a.filterOptions = newOptions
	fmt.Println(i18n.T("filter_applied"))
	return a.cmdListDir([]string{})
}

func (a *App) cmdViewLog(args []string) error {
	maxEntries := 10

	if len(args) > 0 {
		var err error
		maxEntries, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("некорректное количество записей: %w", err)
		}
	}

	entries := a.logger.GetEntries(maxEntries)

	fmt.Printf("Журнал операций (последние %d):\n\n", len(entries))
	for _, entry := range entries {
		fmt.Println(logger.FormatEntryForDisplay(entry))
	}

	return nil
}

func (a *App) cmdToggleColors(_ []string) error {
	a.display.ToggleColors()

	if a.display.UseColors {
		fmt.Println(i18n.T("colors_on"))
	} else {
		fmt.Println(i18n.T("colors_off"))
	}

	return nil
}

func (a *App) cmdEmptyTrash(_ []string) error {
	err := a.fileOperator.SoftDeleter.EmptyTrash()
	if err != nil {
		return fmt.Errorf(i18n.T("error"), err)
	}
	fmt.Println(i18n.T("trash_empty"))
	return nil
}

func (a *App) cmdTrashList(_ []string) error {
	files, err := a.fileOperator.SoftDeleter.ListTrash()
	if err != nil {
		return fmt.Errorf(i18n.T("error"), err)
	}
	if len(files) == 0 {
		fmt.Println(i18n.T("trash_empty_already"))
		return nil
	}
	fmt.Println(i18n.T("trash_contents"))
	for i, f := range files {
		fmt.Printf("%d. %s\n", i+1, f)
	}
	return nil
}

func (a *App) cmdRestoreFromTrash(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf(i18n.T("args_expected_1"), len(args))
	}
	err := a.fileOperator.SoftDeleter.RestoreFromTrash(args[0])
	if err != nil {
		return fmt.Errorf(i18n.T("error"), err)
	}
	fmt.Println(i18n.T("file_restored"))
	return nil
}