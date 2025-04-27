#!/bin/bash
# Скрипт для установки file-manager (fm) в системный PATH для Linux/macOS

# Проверка наличия прав администратора
if [ "$(id -u)" != "0" ]; then
    echo "Этот скрипт требует прав администратора для установки в системную директорию."
    echo "Пожалуйста, запустите скрипт с sudo: sudo ./install.sh"
    exit 1
fi

# Проверка операционной системы
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    INSTALL_DIR="/usr/local/bin"
    echo "Обнаружена macOS, установка в $INSTALL_DIR"
else
    # Linux и другие системы
    INSTALL_DIR="/usr/local/bin"
    echo "Обнаружена Linux или другая UNIX-подобная система, установка в $INSTALL_DIR"
fi

# Проверка наличия исполняемого файла
if [[ "$OSTYPE" == "darwin"* ]]; then
    # Для macOS компилируем отдельно
    echo "Компиляция для macOS..."
    go build -o fm cmd/filemanager/main.go
    if [ ! -f "fm" ]; then
        echo "Ошибка компиляции. Пожалуйста, убедитесь, что у вас установлен Go и проект настроен правильно."
        exit 1
    fi
else
    # Для Linux компилируем отдельно
    echo "Компиляция для Linux..."
    go build -o fm cmd/filemanager/main.go
    if [ ! -f "fm" ]; then
        echo "Ошибка компиляции. Пожалуйста, убедитесь, что у вас установлен Go и проект настроен правильно."
        exit 1
    fi
fi

# Копирование исполняемого файла
cp -f fm "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/fm"
echo "Файл fm скопирован в $INSTALL_DIR и сделан исполняемым"

echo ""
echo "Установка завершена! Теперь вы можете использовать команду 'fm' из любой директории."
echo "Примеры использования:"
echo "fm ls                    # Показать содержимое текущей директории"
echo "fm help                  # Показать список доступных команд"
echo "fm mkdir test_dir        # Создать директорию"
echo "fm touch test_file.txt   # Создать файл"
echo "" 