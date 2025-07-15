// Package app реализует основную логику командного интерфейса файлового менеджера.
package app

// registerCommands регистрирует все доступные команды
func (a *App) registerCommands() {
	a.commands = map[string]Command{
		"help": {
			Name:        "help",
			Description: "Показать список доступных команд",
			Execute:     a.cmdHelp,
		},
		"ls": {
			Name:        "ls",
			Description: "Показать содержимое текущей директории",
			Execute:     a.cmdListDir,
		},
		"cd": {
			Name:        "cd",
			Description: "Изменить текущую директорию: cd <путь>",
			Execute:     a.cmdChangeDir,
		},
		"pwd": {
			Name:        "pwd",
			Description: "Показать текущую директорию",
			Execute:     a.cmdPrintWorkingDir,
		},
		"mkdir": {
			Name:        "mkdir",
			Description: "Создать новую директорию: mkdir <имя>",
			Execute:     a.cmdMakeDir,
		},
		"touch": {
			Name:        "touch",
			Description: "Создать новый файл: touch <имя>",
			Execute:     a.cmdCreateFile,
		},
		"rm": {
			Name:        "rm",
			Description: "Удалить файл: rm <имя>",
			Execute:     a.cmdRemoveFile,
		},
		"rmdir": {
			Name:        "rmdir",
			Description: "Удалить директорию: rmdir <имя>",
			Execute:     a.cmdRemoveDir,
		},
		"cp": {
			Name:        "cp",
			Description: "Копировать файл/директорию: cp <источник> <назначение>",
			Execute:     a.cmdCopy,
		},
		"mv": {
			Name:        "mv",
			Description: "Переместить/переименовать файл/директорию: mv <источник> <назначение>",
			Execute:     a.cmdMove,
		},
		"find": {
			Name:        "find",
			Description: "Найти файлы по имени: find <шаблон>",
			Execute:     a.cmdFindByName,
		},
		"grep": {
			Name:        "grep",
			Description: "Найти файлы по содержимому: grep <текст>",
			Execute:     a.cmdFindByContent,
		},
		"info": {
			Name:        "info",
			Description: "Показать информацию о файле/директории: info <имя>",
			Execute:     a.cmdFileInfo,
		},
		"exit": {
			Name:        "exit",
			Description: "Выйти из программы",
			Execute:     a.cmdExit,
		},
		"cat": {
			Name:        "cat",
			Description: "Просмотр содержимого текстового файла: cat <имя> [начальная_строка] [количество_строк]",
			Execute:     a.cmdViewFile,
		},
		"chmod": {
			Name:        "chmod",
			Description: "Изменить права доступа к файлу: chmod <режим> <имя>",
			Execute:     a.cmdChangePermissions,
		},
		"archive": {
			Name:        "archive",
			Description: "Создать архив: archive <имя_архива> <формат> <файл1> [файл2...]",
			Execute:     a.cmdCreateArchive,
		},
		"extract": {
			Name:        "extract",
			Description: "Распаковать архив: extract <архив> <директория>",
			Execute:     a.cmdExtractArchive,
		},
		"list-archive": {
			Name:        "list-archive",
			Description: "Показать содержимое архива: list-archive <архив>",
			Execute:     a.cmdListArchive,
		},
		"bookmark": {
			Name:        "bookmark",
			Description: "Управление закладками: bookmark add <имя> [путь] | list | remove <имя> | go <имя>",
			Execute:     a.cmdManageBookmarks,
		},
		"filter": {
			Name:        "filter",
			Description: "Фильтрация файлов: filter [--ext=<расширение>] [--name=<шаблон>] [--size=<мин>-<макс>] [--date=<начало>-<конец>] [--type=<f|d|h>]",
			Execute:     a.cmdFilter,
		},
		"log": {
			Name:        "log",
			Description: "Просмотр журнала операций: log [количество]",
			Execute:     a.cmdViewLog,
		},
		"colors": {
			Name:        "colors",
			Description: "Включить/отключить цветной вывод",
			Execute:     a.cmdToggleColors,
		},
		"empty-trash": {
			Name:        "empty-trash",
			Description: "Очистить корзину (удалить все файлы)",
			Execute:     a.cmdEmptyTrash,
		},
		"trash-list": {
			Name:        "trash-list",
			Description: "Показать содержимое корзины",
			Execute:     a.cmdTrashList,
		},
		"restore": {
			Name:        "restore",
			Description: "Восстановить файл из корзины (Linux)",
			Execute:     a.cmdRestoreFromTrash,
		},
	}
}