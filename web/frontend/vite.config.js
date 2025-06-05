import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import fs from 'fs';
import path from 'path';

// Функция для чтения конфигурации бота (для получения альтернативных портов)
function getBotConfig() {
  try {
    // Путь к файлу конфигурации может отличаться в зависимости от структуры проекта
    // Здесь предполагается, что он находится в корне проекта
    const configPath = path.resolve(__dirname, '../../config.json');
    if (fs.existsSync(configPath)) {
      const configData = fs.readFileSync(configPath, 'utf8');
      return JSON.parse(configData);
    }
  } catch (error) {
    console.warn('Не удалось прочитать файл конфигурации:', error.message);
  }
  return { WebInterface: { Port: 8080, AltPorts: [3000, 8000] } }; // Значения по умолчанию
}

// Получаем порт API из переменной окружения или используем порт по умолчанию
const apiPort = process.env.API_PORT || getBotConfig().WebInterface.Port || 8080;

// Получаем альтернативные порты из конфигурации
const altPorts = process.env.ALT_PORTS ? 
  process.env.ALT_PORTS.split(',').map(port => parseInt(port.trim())) : 
  getBotConfig().WebInterface.AltPorts || [3000, 8000];

console.log(`Основной порт API: ${apiPort}`);
console.log(`Альтернативные порты: ${altPorts.join(', ')}`);

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: `http://localhost:${apiPort}`,
        changeOrigin: true,
        // Добавляем обработку ошибок для автоматического переключения на альтернативные порты
        configure: (proxy, _options) => {
          let currentPortIndex = -1;
          let isUsingAltPort = false;
          
          proxy.on('error', (err, req, res) => {
            console.log('Ошибка прокси:', err.message);
            
            // Если уже используется альтернативный порт и он не работает, пробуем следующий
            if (isUsingAltPort) {
              currentPortIndex = (currentPortIndex + 1) % altPorts.length;
            } else {
              // Первая ошибка - переключаемся на первый альтернативный порт
              currentPortIndex = 0;
              isUsingAltPort = true;
            }
            
            const newPort = altPorts[currentPortIndex];
            console.log(`Переключение на альтернативный порт: ${newPort}`);
            
            // Обновляем целевой порт для прокси
            proxy.options.target = `http://localhost:${newPort}`;
            
            // Повторяем запрос с новым портом
            proxy.web(req, res, proxy.options);
          });
        },
      },
    },
  },
  build: {
    outDir: 'build',
    emptyOutDir: true,
    sourcemap: false,
  },
});