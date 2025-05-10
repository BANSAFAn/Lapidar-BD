import axios from 'axios';

// Базовый URL для API запросов
const API_URL = '/api';

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
      const response = await axios.get(`${API_URL}/config`);
      return response.data;
    } catch (error) {
      console.error('Ошибка при получении конфигурации:', error);
      throw error;
    }
  },

  // Сохранение конфигурации бота
  saveConfig: async (config: BotConfig): Promise<{ status: string }> => {
    try {
      const response = await axios.post(`${API_URL}/save-config`, config);
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
      // const response = await axios.get(`${API_URL}/stats`);
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
      // const response = await axios.get(`${API_URL}/commands`);
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
  updateCommand: async (command: Command): Promise<{ status: string }> => {
    try {
      // В будущем здесь будет реальный запрос к API
      // const response = await axios.post(`${API_URL}/update-command`, command);
      // return response.data;
      
      // Временная заглушка для демонстрации
      return { status: 'success' };
    } catch (error) {
      console.error('Ошибка при обновлении команды:', error);
      throw error;
    }
  }
};

export default apiService;