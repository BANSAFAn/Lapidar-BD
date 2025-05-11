package config

import (
	"encoding/json"
	"os"
	"time"
)

// AdminConfig содержит настройки администратора для веб-панели
type AdminConfig struct {
	Email    string `json:"email"`    // Email администратора
	Password string `json:"password"` // Хеш пароля администратора
	Secret   string `json:"secret"`   // Секретный ключ для генерации TOTP кодов
}

// LoginLog содержит информацию о входе в систему
type LoginLog struct {
	Email     string    `json:"email"`      // Email пользователя
	IP        string    `json:"ip"`         // IP-адрес
	UserAgent string    `json:"user_agent"` // User-Agent браузера
	Timestamp time.Time `json:"timestamp"`  // Время входа
	Success   bool      `json:"success"`    // Успешный вход или нет
}

// LoadAdminConfig загружает конфигурацию администратора из файла
func LoadAdminConfig() (*AdminConfig, error) {
	file, err := os.Open("config/admin.json")
	if err != nil {
		// Если файл не существует, создаем конфигурацию по умолчанию
		if os.IsNotExist(err) {
			// Создаем директорию config, если она не существует
			if _, err := os.Stat("config"); os.IsNotExist(err) {
				if err := os.Mkdir("config", 0755); err != nil {
					return nil, err
				}
			}

			// Создаем конфигурацию по умолчанию
			// Пароль по умолчанию: admin
			defaultConfig := &AdminConfig{
				Email:    "admin@example.com",
				Password: "$2a$10$XgXLGQAJAYv8CKJE2aJzSO0CT6.uOmOEy0Oj.1iP.hO2JJw2aN12O", // хеш для 'admin'
				Secret:   "JBSWY3DPEHPK3PXP",                                             // Пример секрета для TOTP
			}

			// Создаем файл с конфигурацией по умолчанию
			if err := SaveAdminConfig(defaultConfig); err != nil {
				return nil, err
			}

			return defaultConfig, nil
		}
		return nil, err
	}
	defer file.Close()

	config := &AdminConfig{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	return config, err
}

// SaveAdminConfig сохраняет конфигурацию администратора в файл
func SaveAdminConfig(config *AdminConfig) error {
	file, err := os.Create("config/admin.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}
