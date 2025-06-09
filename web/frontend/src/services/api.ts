import axios from 'axios';

// Базовый URL для API запросов
const API_URL = '/api';

// Создаем экземпляр axios с настроенными перехватчиками
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Перехватчик ответов для сохранения CSRF токена
api.interceptors.response.use(
  (response) => {
    // Если в ответе есть заголовок X-CSRF-Token, сохраняем его
    const csrfToken = response.headers['x-csrf-token'];
    if (csrfToken) {
      // Сохраняем токен для использования в следующих запросах
      localStorage.setItem('csrfToken', csrfToken);
    }
    return response;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Перехватчик запросов для добавления CSRF токена
api.interceptors.request.use(
  (config) => {
    // Получаем сохраненный CSRF токен
    const csrfToken = localStorage.getItem('csrfToken');
    
    // Если токен есть и запрос не GET, добавляем его в заголовки
    if (csrfToken && config.method !== 'get') {
      config.headers['X-CSRF-Token'] = csrfToken;
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Интерфейсы для типизации данных
export interface BotConfig {
  Token: string;
  Prefix: string;
  BotName: string;
  DefaultLanguage: string;
  WebInterface: {
    Enabled: boolean;
    Host: string;
    Port: number;
  };
}

export interface BotStats {
  servers: number;
  users: number;
  channels: number;
  commands: number;
  uptime: string;
  memoryUsage: string;
}

export interface Command {
  name: string;
  description: string;
  usage: string;
  category: string;
  enabled: boolean;
}

// API сервис для взаимодействия с бэкендом
const apiService = {
  // Получение конфигурации бота
  getConfig: async (): Promise<BotConfig> => {
    try {
      const response = await api.get('/config');
      return response.data;
    } catch (error) {
      console.error('Ошибка при получении конфигурации:', error);
      throw error;
    }
  },

  // Сохранение конфигурации бота
  saveConfig: async (config: BotConfig): Promise<{ status: string }> => {
    try {
      const response = await api.post('/save-config', config);
      return response.data;
    } catch (error) {
      console.error('Ошибка при сохранении конфигурации:', error);
      throw error;
    }
  },

  // Получение статистики бота (заглушка, будет реализована на бэкенде)
  getStats: async (): Promise<BotStats> => {
    try {
      // В будущем здесь будет реальный запрос к API
      // const response = await api.get('/stats');
      // return response.data;
      
      // Временная заглушка для демонстрации
      return {
        servers: 15,
        users: 1250,
        channels: 87,
        commands: 42,
        uptime: '3 дня 7 часов',
        memoryUsage: '128 MB'
      };
    } catch (error) {
      console.error('Ошибка при получении статистики:', error);
      throw error;
    }
  },

  // Получение списка команд (заглушка, будет реализована на бэкенде)
  getCommands: async (): Promise<Command[]> => {
    try {
      // В будущем здесь будет реальный запрос к API
      // const response = await api.get('/commands');
      // return response.data;
      
      // Временная заглушка для демонстрации
      return [
        {
          name: 'help',
          description: 'Показывает список доступных команд',
          usage: '!help [команда]',
          category: 'Основные',
          enabled: true
        },
        {
          name: 'ping',
          description: 'Проверяет задержку бота',
          usage: '!ping',
          category: 'Утилиты',
          enabled: true
        },
        {
          name: 'ban',
          description: 'Банит пользователя на сервере',
          usage: '!ban @пользователь [причина]',
          category: 'Модерация',
          enabled: true
        },
        {
          name: 'play',
          description: 'Воспроизводит музыку в голосовом канале',
          usage: '!play [ссылка или название]',
          category: 'Музыка',
          enabled: true
        },
        {
          name: 'stats',
          description: 'Показывает статистику бота',
          usage: '!stats',
          category: 'Информация',
          enabled: true
        }
      ];
    } catch (error) {
      console.error('Ошибка при получении списка команд:', error);
      throw error;
    }
  },

  // Обновление статуса команды (заглушка, будет реализована на бэкенде)
  updateCommand: async (name: string, enabled: boolean): Promise<{ status: string }> => {
    try {
      // В будущем здесь будет реальный запрос к API
      // const response = await api.post('/update-command', { name, enabled });
      // return response.data;
      
      // Временная заглушка для демонстрации
      console.log(`Обновление статуса команды ${name} на ${enabled ? 'включено' : 'выключено'}`);
      return { status: 'success' };
    } catch (error) {
      console.error('Ошибка при обновлении статуса команды:', error);
      throw error;
    }
  },

  // Авторизация пользователя
  login: async (email: string, password: string): Promise<{ token: string }> => {
    try {
      const response = await api.post('/login', { email, password });
      return response.data;
    } catch (error) {
      console.error('Ошибка при авторизации:', error);
      throw error;
    }
  },

  // Проверка TOTP кода
  verifyTOTP: async (token: string, code: string): Promise<{ token: string }> => {
    try {
      const response = await api.post('/verify-totp', { token, code });
      return response.data;
    } catch (error) {
      console.error('Ошибка при проверке TOTP кода:', error);
      throw error;
    }
  },

  // Выход из системы
  logout: async (): Promise<{ status: string }> => {
    try {
      const response = await api.post('/logout');
      return response.data;
    } catch (error) {
      console.error('Ошибка при выходе из системы:', error);
      throw error;
    }
  }
};

export default apiService;