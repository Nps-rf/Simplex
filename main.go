package main

import (
	"fmt"
	"os"
	"os/exec"
)

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
