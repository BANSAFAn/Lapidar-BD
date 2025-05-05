package ai

import (
	"fmt"

	"discord-bot/config"
)

// GeminiProvider реализует интерфейс AIProvider для Gemini
type GeminiProvider struct {
	apiKey      string
	initialized bool
}

// Initialize инициализирует клиент Gemini API
func (p *GeminiProvider) Initialize() error {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Проверяем наличие API ключа
	if cfg.GeminiAPIKey == "" {
		return fmt.Errorf("API ключ Gemini не указан в конфигурации")
	}

	p.apiKey = cfg.GeminiAPIKey
	p.initialized = true
	return nil
}

// GenerateResponse генерирует ответ на запрос пользователя
func (p *GeminiProvider) GenerateResponse(prompt string) (string, error) {
	if !p.initialized {
		return "", fmt.Errorf("сервис Gemini не инициализирован")
	}

	// Здесь будет реализация запроса к API Gemini
	return "[Заглушка] Ответ от Gemini: " + prompt, nil
}

// GetName возвращает название AI модели
func (p *GeminiProvider) GetName() string {
	return "Gemini"
}
