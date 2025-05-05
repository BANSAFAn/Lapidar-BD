package main

import (
	"discord-bot/config"
	"os"
	"strings"
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

	// Удаляем неподдерживаемый флаг testlogfile, который вызывает ошибку
	args := []string{}
	for _, arg := range os.Args {
		if !strings.HasPrefix(arg, "-test.testlogfile") {
			args = append(args, arg)
		}
	}
	os.Args = args

	m.Run()
}
