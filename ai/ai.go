package ai

import (
	"errors"
	"fmt"

	"discord-bot/config"
)

// AIProvider представляет интерфейс для работы с различными AI моделями
type AIProvider interface {
	// Initialize инициализирует провайдер AI
	Initialize() error
	// GenerateResponse генерирует ответ на запрос пользователя
	GenerateResponse(prompt string) (string, error)
	// GetName возвращает название AI модели
	GetName() string
}

// AvailableProviders содержит список доступных AI провайдеров
var AvailableProviders = map[string]AIProvider{
	"gemini":  &GeminiProvider{},
	"grok":    &GrokProvider{},
	"chatgpt": &ChatGPTProvider{},
	"qwen":    &QwenProvider{},
	"claude":  &ClaudeProvider{},
}

// DefaultProvider содержит провайдер AI по умолчанию
var DefaultProvider AIProvider

// Initialize инициализирует все доступные AI провайдеры
func Initialize() error {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Инициализируем провайдеры, для которых есть API ключи
	for name, provider := range AvailableProviders {
		err := provider.Initialize()
		if err == nil {
			// Если это первый успешно инициализированный провайдер, устанавливаем его как провайдер по умолчанию
			if DefaultProvider == nil {
				DefaultProvider = provider
				fmt.Printf("Установлен провайдер AI по умолчанию: %s\n", provider.GetName())
			}
		} else {
			fmt.Printf("Ошибка инициализации провайдера %s: %v\n", name, err)
		}
	}

	if DefaultProvider == nil {
		return errors.New("не удалось инициализировать ни один провайдер AI")
	}

	return nil
}

// GetProvider возвращает провайдер AI по имени
func GetProvider(name string) (AIProvider, error) {
	provider, exists := AvailableProviders[name]
	if !exists {
		return nil, fmt.Errorf("провайдер AI '%s' не найден", name)
	}
	return provider, nil
}

// GenerateResponse генерирует ответ на запрос пользователя, используя провайдер по умолчанию
func GenerateResponse(prompt string) (string, error) {
	if DefaultProvider == nil {
		return "", errors.New("провайдер AI не инициализирован")
	}
	return DefaultProvider.GenerateResponse(prompt)
}
