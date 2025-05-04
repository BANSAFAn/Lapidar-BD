package main

import (
	"discord-bot/config"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	// Тест загрузки конфигурации по умолчанию
	cfg, err := config.Load()
	if err != nil {
		t.Errorf("Ошибка загрузки конфигурации: %v", err)
	}

	// Проверка значений по умолчанию
	if cfg.Prefix == "" {
		t.Error("Префикс команд не должен быть пустым")
	}

	if cfg.DefaultLanguage == "" {
		t.Error("Язык по умолчанию не должен быть пустым")
	}
}

func TestMain(m *testing.M) {
	// Здесь можно добавить код для настройки тестового окружения
	// перед запуском тестов
	m.Run()
}
