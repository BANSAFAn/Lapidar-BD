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

	// Создаем и запускаем API сервер для новой веб-панели
	apiServer := web.NewAPIServer(cfg)
	if err := apiServer.Start(); err != nil {
		fmt.Println("Ошибка запуска веб-интерфейса:", err)
	}
}
