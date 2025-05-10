@echo off
REM Скрипт для установки git hooks в Windows

echo Установка git hooks...

REM Создаем директорию .git/hooks, если она не существует
if not exist ".git\hooks" mkdir .git\hooks

REM Копируем pre-commit хук
copy /Y ".github\hooks\pre-commit" ".git\hooks\"

echo Git hooks успешно установлены!
echo Pre-commit хук будет проверять код на уязвимости перед каждым коммитом.

REM Установка необходимых инструментов для проверки безопасности
echo Установка инструментов для проверки безопасности...

REM Установка gosec для Go
where go >nul 2>nul
if %ERRORLEVEL% == 0 (
  echo Установка gosec...
  go install github.com/securego/gosec/v2/cmd/gosec@latest
) else (
  echo Go не установлен. Пропуск установки gosec.
)

REM Установка зависимостей для фронтенда
if exist "web\frontend" (
  cd web\frontend
  
  if exist "package.json" (
    echo Установка зависимостей для фронтенда...
    call npm install
  )
  
  cd ..\..
)

echo Установка завершена!
pause