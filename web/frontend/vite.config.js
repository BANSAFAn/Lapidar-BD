import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

// Получаем порт API из переменной окружения или используем порт по умолчанию
const apiPort = process.env.API_PORT || 8080;

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
          proxy.on('error', (err, _req, _res) => {
            console.log('Ошибка прокси:', err);
            console.log('Попробуйте использовать другой порт API через переменную окружения API_PORT');
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