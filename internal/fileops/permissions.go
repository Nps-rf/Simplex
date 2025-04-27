package fileops

import (
	"fmt"
	"os"
	"strconv"
)

// PermissionsManager предоставляет функции для управления правами доступа к файлам
type PermissionsManager struct{}

// NewPermissionsManager создает новый экземпляр PermissionsManager
func NewPermissionsManager() *PermissionsManager {
	return &PermissionsManager{}
}

// ChangePermissions изменяет права доступа к файлу или директории
func (p *PermissionsManager) ChangePermissions(path string, permissions string) error {
	// Преобразование строки с восьмеричным числом в uint32
	mode, err := strconv.ParseUint(permissions, 8, 32)
	if err != nil {
		return fmt.Errorf("некорректный формат прав доступа: %w", err)
	}

	// Применение новых прав доступа
	err = os.Chmod(path, os.FileMode(mode))
	if err != nil {
		return fmt.Errorf("не удалось изменить права доступа для %s: %w", path, err)
	}

	return nil
}

// GetPermissions возвращает текущие права доступа к файлу или директории
func (p *PermissionsManager) GetPermissions(path string) (string, error) {
	// Получение информации о файле
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("не удалось получить информацию о %s: %w", path, err)
	}

	// Преобразование прав доступа в строку восьмеричного числа
	mode := fileInfo.Mode().Perm()
	return fmt.Sprintf("%o", mode), nil
}

// ChangeOwner изменяет владельца файла или директории (не работает в Windows)
func (p *PermissionsManager) ChangeOwner(path string, uid, gid int) error {
	err := os.Chown(path, uid, gid)
	if err != nil {
		return fmt.Errorf("не удалось изменить владельца для %s: %w", path, err)
	}

	return nil
}

// FormatPermissions форматирует права доступа в удобочитаемом виде (как в ls -l)
func (p *PermissionsManager) FormatPermissions(mode os.FileMode) string {
	result := ""

	// Тип файла
	if mode.IsDir() {
		result += "d"
	} else if mode&os.ModeSymlink != 0 {
		result += "l"
	} else if mode&os.ModeNamedPipe != 0 {
		result += "p"
	} else if mode&os.ModeSocket != 0 {
		result += "s"
	} else if mode&os.ModeDevice != 0 {
		result += "b"
	} else if mode&os.ModeCharDevice != 0 {
		result += "c"
	} else {
		result += "-"
	}

	// Права для владельца
	result += p.formatPermissionBits(mode, 6)

	// Права для группы
	result += p.formatPermissionBits(mode, 3)

	// Права для остальных
	result += p.formatPermissionBits(mode, 0)

	return result
}

// formatPermissionBits форматирует биты прав доступа
func (p *PermissionsManager) formatPermissionBits(mode os.FileMode, shift uint) string {
	var result string

	// Чтение
	if mode&(1<<(shift+2)) != 0 {
		result += "r"
	} else {
		result += "-"
	}

	// Запись
	if mode&(1<<(shift+1)) != 0 {
		result += "w"
	} else {
		result += "-"
	}

	// Исполнение
	if mode&(1<<shift) != 0 {
		result += "x"
	} else {
		result += "-"
	}

	return result
}
