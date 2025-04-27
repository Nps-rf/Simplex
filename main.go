package main

import (
	"fmt"
	"os"
	"os/exec"
)

// main запускает файловый менеджер
func main() {
	// Запуск через cmd/filemanager/main.go
	cmd := exec.Command("go", "run", "cmd/filemanager/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при запуске файлового менеджера: %v\n", err)
		os.Exit(1)
	}
}

// AddBookmark создает новую закладку.
func AddBookmark(alias, path string) error {
	// Проверяем, нет ли уже такого alias
	// Добавляем в структуру bookmarks и сохраняем в JSON
	return nil
}

// GoBookmark осуществляет переход к закладке.
func GoBookmark(alias string) (string, error) {
	// Находит соответствующий path в bookmarks,
	// возвращает его, чтобы затем Nav мог выполнить cd.
	path := "/path/to/bookmark" // Заглушка для предотвращения ошибки компиляции
	return path, nil
}
