# Скрипт для установки file-manager (fm.exe) в системный PATH
# Для запуска: PowerShell -ExecutionPolicy Bypass -File install.ps1

# Функция для проверки наличия прав администратора
function Test-Admin {
    $currentUser = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    $currentUser.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Проверяем наличие прав администратора
if (-not (Test-Admin)) {
    Write-Host "Этот скрипт требует прав администратора для установки в системную директорию."
    Write-Host "Пожалуйста, запустите PowerShell от имени администратора и попробуйте снова."
    exit 1
}

# Создаем директорию для приложения, если её нет
$installDir = "C:\Program Files\FileManager"
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
    Write-Host "Создана директория $installDir"
}

# Путь к исполняемому файлу
$exePath = Join-Path $PSScriptRoot "fm.exe"

# Проверяем, существует ли исполняемый файл
if (-not (Test-Path $exePath)) {
    Write-Host "Файл fm.exe не найден в текущей директории."
    Write-Host "Пожалуйста, убедитесь, что вы запускаете скрипт из директории, содержащей скомпилированный файл fm.exe."
    exit 1
}

# Копируем исполняемый файл в директорию установки
Copy-Item -Path $exePath -Destination $installDir -Force
Write-Host "Файл fm.exe скопирован в $installDir"

# Проверяем, находится ли директория уже в PATH
$envPath = [Environment]::GetEnvironmentVariable("PATH", [EnvironmentVariableTarget]::Machine)
if ($envPath -notlike "*$installDir*") {
    # Добавляем директорию в PATH
    $newPath = "$envPath;$installDir"
    [Environment]::SetEnvironmentVariable("PATH", $newPath, [EnvironmentVariableTarget]::Machine)
    Write-Host "Директория $installDir добавлена в системный PATH"
    
    # Обновляем PATH в текущей сессии
    $env:PATH = "$env:PATH;$installDir"
}
else {
    Write-Host "Директория $installDir уже находится в системном PATH"
}

# Создаем командный алиас (опционально)
$aliasPath = Join-Path $installDir "fm.cmd"
@"
@echo off
"%~dp0fm.exe" %*
"@ | Out-File -FilePath $aliasPath -Encoding ASCII -Force
Write-Host "Создан командный алиас fm.cmd"

Write-Host ""
Write-Host "Установка завершена! Теперь вы можете использовать команду 'fm' из любой директории."
Write-Host "Примеры использования:"
Write-Host "fm ls                    # Показать содержимое текущей директории"
Write-Host "fm help                  # Показать список доступных команд"
Write-Host "fm mkdir test_dir        # Создать директорию"
Write-Host "fm touch test_file.txt   # Создать файл"
Write-Host ""
Write-Host "Примечание: Может потребоваться перезапуск командной строки или PowerShell для применения изменений PATH." 