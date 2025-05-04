package ai

import (
	"context"
	"fmt"

	"discord-bot/config"

	geminiai "google.golang.org/api/generativelanguage/v1"
	"google.golang.org/api/option"
)

// GeminiProvider реализует интерфейс AIProvider для Gemini AI
type GeminiProvider struct {
	service *geminiai.Service
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

	// Инициализируем сервис
	ctx := context.Background()
	p.service, err = geminiai.NewService(ctx, option.WithAPIKey(cfg.GeminiAPIKey))
	if err != nil {
		return fmt.Errorf("ошибка создания сервиса Gemini: %w", err)
	}

	return nil
}

// GenerateResponse генерирует ответ на запрос пользователя
func (p *GeminiProvider) GenerateResponse(prompt string) (string, error) {
	if p.service == nil {
		return "", fmt.Errorf("сервис Gemini не инициализирован")
	}

	// Создаем запрос
	request := &geminiai.GenerateContentRequest{
		Contents: []*geminiai.Content{
			{
				Parts: []*geminiai.Part{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	// Отправляем запрос
	resp, err := p.service.Models.GenerateContent("models/gemini-pro", request).Do()
	if err != nil {
		return "", fmt.Errorf("ошибка при запросе к Gemini API: %w", err)
	}

	// Обрабатываем ответ
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("пустой ответ от Gemini API")
	}

	return resp.Candidates[0].Content.Parts[0].Text, nil
}

// GetName возвращает название AI модели
func (p *GeminiProvider) GetName() string {
	return "Gemini"
}
