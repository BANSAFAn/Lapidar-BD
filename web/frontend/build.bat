@echo off
echo Установка зависимостей для фронтенда...
npm install

echo Сборка React приложения...
npm run build

echo Копирование собранных файлов в директорию web/frontend/build...
if not exist ..\frontend\build mkdir ..\frontend\build
xcopy /E /Y build\* ..\frontend\build\

echo Сборка завершена! Веб-панель готова к использованию.