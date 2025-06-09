package config

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"time"
)

// AdminConfig содержит настройки администратора для веб-панели
type AdminConfig struct {
	Email     string `json:"email"`      // Email администратора
	Password  string `json:"password"`  // Хеш пароля администратора
	TOTPSecret string `json:"totp_secret"` // Секретный ключ для генерации TOTP кодов
	JWTSecret  string `json:"jwt_secret"`  // Секретный ключ для генерации JWT токенов
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
				Email:      "admin@example.com",
				Password:   "$2a$10$XgXLGQAJAYv8CKJE2aJzSO0CT6.uOmOEy0Oj.1iP.hO2JJw2aN12O", // хеш для 'admin'
				TOTPSecret: "JBSWY3DPEHPK3PXP",                                             // Пример секрета для TOTP
				JWTSecret:  GenerateSecureRandomString(32),                                   // Генерируем случайный JWT секрет
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

	// Временная структура для поддержки старого формата
	type OldAdminConfig struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Secret   string `json:"secret"`
	}

	// Сначала пробуем загрузить в старом формате
	var oldConfig OldAdminConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&oldConfig)
	if err != nil {
		// Если не удалось декодировать в старом формате, пробуем в новом
		file.Seek(0, 0) // Возвращаемся в начало файла
		config := &AdminConfig{}
		decoder = json.NewDecoder(file)
		err = decoder.Decode(config)
		return config, err
	}

	// Если загрузили в старом формате, мигрируем в новый
	config := &AdminConfig{
		Email:      oldConfig.Email,
		Password:   oldConfig.Password,
		TOTPSecret: oldConfig.Secret,
		JWTSecret:  GenerateSecureRandomString(32),
	}

	// Сохраняем в новом формате
	if err := SaveAdminConfig(config); err != nil {
		return nil, err
	}

	return config, nil
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

// GenerateSecureRandomString генерирует криптографически стойкую случайную строку заданной длины
func GenerateSecureRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// В случае ошибки возвращаем фиксированную строку, но это крайне маловероятно
		return "SECURE_JWT_SECRET_FALLBACK_DO_NOT_USE_IN_PRODUCTION"
	}
	
	// Кодируем в base64 и обрезаем до нужной длины
	encoded := base64.StdEncoding.EncodeToString(b)
	if len(encoded) > length {
		return encoded[:length]
	}
	return encoded
}
