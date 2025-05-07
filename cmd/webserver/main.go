package main

import (
	"discord-bot/config"
	"discord-bot/web"
	"fmt"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	// Убедимся, что веб-интерфейс включен
	cfg.WebInterface.Enabled = true

	// Создаем и запускаем веб-сервер
	webServer := web.NewWebServer(cfg)
	if err := webServer.Start(); err != nil {
		fmt.Println("Ошибка запуска веб-интерфейса:", err)
	}
}
