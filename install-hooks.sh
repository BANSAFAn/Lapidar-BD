#!/bin/bash

# Скрипт для установки git hooks

echo "Установка git hooks..."

# Создаем директорию .git/hooks, если она не существует
mkdir -p .git/hooks

# Копируем pre-commit хук
cp .github/hooks/pre-commit .git/hooks/

# Делаем хук исполняемым
chmod +x .git/hooks/pre-commit

echo "Git hooks успешно установлены!"
echo "Pre-commit хук будет проверять код на уязвимости перед каждым коммитом."

# Установка необходимых инструментов для проверки безопасности
echo "Установка инструментов для проверки безопасности..."

# Установка gosec для Go
if command -v go >/dev/null 2>&1; then
  echo "Установка gosec..."
  go install github.com/securego/gosec/v2/cmd/gosec@latest
else
  echo "Go не установлен. Пропуск установки gosec."
fi

# Установка зависимостей для фронтенда
if [ -d "./web/frontend" ]; then
  cd ./web/frontend
  
  if [ -f "package.json" ]; then
    echo "Установка зависимостей для фронтенда..."
    npm install
  fi
  
  cd - > /dev/null
fi

echo "Установка завершена!"