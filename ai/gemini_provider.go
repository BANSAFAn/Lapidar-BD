package ai

import (
	"fmt"

	"discord-bot/config"
)

type GeminiProvider struct{}

func (p *GeminiProvider) Initialize() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	if cfg.GeminiAPIKey == "" {
		return fmt.Errorf("API ключ Gemini не указан в конфигурации")
	}

	return nil
}

func (p *GeminiProvider) GenerateResponse(prompt string) (string, error) {
	return "Функция временно недоступна", nil
}

func (p *GeminiProvider) GetName() string {
	return "Gemini"
}
