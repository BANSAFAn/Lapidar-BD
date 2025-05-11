package main

import (
	"discord-bot/config"
	"discord-bot/db"
	"discord-bot/web"
	"fmt"
	"time"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	// Загрузка конфигурации администратора
	_, err = config.LoadAdminConfig()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации администратора:", err)
		return
	}

	// Инициализация базы данных
	err = db.InitializeDB()
	if err != nil {
		fmt.Println("Ошибка инициализации базы данных:", err)
		return
	}

	// Инициализация таблиц для аутентификации
	err = db.InitAuthTables()
	if err != nil {
		fmt.Println("Ошибка инициализации таблиц аутентификации:", err)
		return
	}

	// Запускаем периодическую очистку просроченных сессий
	go func() {
		for {
			time.Sleep(time.Hour)
			if err := db.CleanExpiredSessions(); err != nil {
				fmt.Println("Ошибка очистки просроченных сессий:", err)
			}
		}
	}()

	// Убедимся, что веб-интерфейс включен
	cfg.WebInterface.Enabled = true

	// Создаем и запускаем API сервер для новой веб-панели
	apiServer := web.NewAPIServer(cfg)

	// Обновляем API сервер для поддержки аутентификации
	apiServer.UpdateAPIServerForAuth()

	// Запускаем API сервер
	if err := apiServer.Start(); err != nil {
		fmt.Println("Ошибка запуска веб-интерфейса:", err)
	}
}
