package ai

import (
	"errors"
	"fmt"

	"discord-bot/config"
)

type AIProvider interface {
	Initialize() error
	GenerateResponse(prompt string) (string, error)
	GetName() string
}

var AvailableProviders = map[string]AIProvider{
	"gemini":  &GeminiProvider{},
	"grok":    &GrokProvider{},
	"chatgpt": &ChatGPTProvider{},
	"qwen":    &QwenProvider{},
	"claude":  &ClaudeProvider{},
}

var DefaultProvider AIProvider

func Initialize() error {
	if _, err := config.Load(); err != nil {
		return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	for name, provider := range AvailableProviders {
		err := provider.Initialize()
		if err == nil {
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

func GetProvider(name string) (AIProvider, error) {
	provider, exists := AvailableProviders[name]
	if !exists {
		return nil, fmt.Errorf("провайдер AI '%s' не найден", name)
	}
	return provider, nil
}

func GenerateResponse(prompt string) (string, error) {
	if DefaultProvider == nil {
		return "", errors.New("провайдер AI не инициализирован")
	}
	return DefaultProvider.GenerateResponse(prompt)
}
