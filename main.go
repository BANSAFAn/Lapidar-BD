package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"discord-bot/ai"
	"discord-bot/config"
	"discord-bot/db"
	"discord-bot/handlers"
	"discord-bot/localization"
	"discord-bot/web"

	"github.com/bwmarrin/discordgo"
)

var s *discordgo.Session

func init() {
	flag.Parse()
}

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	// Запуск веб-интерфейса, если он включен
	if cfg.WebInterface.Enabled {
		webServer := web.NewWebServer(cfg)
		go func() {
			if err := webServer.Start(); err != nil {
				fmt.Println("Ошибка запуска веб-интерфейса:", err)
			}
		}()
	}

	// Инициализация базы данных
	dbConfig := db.DatabaseConfig{
		Type:     "sqlite",
		Database: "data/bot.db",
	}
	if err := db.Initialize(dbConfig); err != nil {
		fmt.Println("Ошибка инициализации базы данных:", err)
		return
	}

	// Инициализация AI провайдеров
	if err := ai.Initialize(); err != nil {
		fmt.Println("Ошибка инициализации AI провайдеров:", err)
		// Продолжаем работу даже при ошибке AI
		fmt.Println("Бот будет работать с ограниченной функциональностью AI")
	}

	// Инициализация системы локализации
	if err := localization.Initialize(); err != nil {
		fmt.Println("Ошибка инициализации системы локализации:", err)
		return
	}

	// Создание новой сессии Discord
	s, err = discordgo.New("Bot " + cfg.Token)
	if err != nil {
		fmt.Println("Ошибка создания сессии Discord:", err)
		return
	}

	// Регистрация обработчиков событий
	s.AddHandler(handlers.MessageCreate)
	s.AddHandler(handlers.ReactionAdd)

	// Добавляем интенты для получения информации о пользователях
	s.Identify.Intents |= discordgo.IntentsGuildMembers

	// Инициализация слеш-команд AI
	if err := handlers.InitAICommands(s); err != nil {
		fmt.Println("Ошибка инициализации слеш-команд AI:", err)
		// Продолжаем работу даже при ошибке инициализации команд
		fmt.Println("Бот будет работать без слеш-команд AI")
	}

	// Добавляем интенты для голосовых каналов
	s.Identify.Intents |= discordgo.IntentsGuildVoiceStates

	// Открытие соединения с Discord
	err = s.Open()
	if err != nil {
		fmt.Println("Ошибка открытия соединения:", err)
		return
	}

	// Выводим красивое сообщение о запуске бота
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                        ║")
	fmt.Printf("║  🤖 %s успешно запущен!                        ║\n", cfg.BotName)
	fmt.Printf("║  🔧 Префикс команд: '%s'                                ║\n", cfg.Prefix)
	fmt.Println("║  ✅ Бот полностью функционирует и готов к работе       ║")
	if cfg.WebInterface.Enabled {
		fmt.Printf("║  🌐 Панель управления: http://%s:%d                ║\n", cfg.WebInterface.Host, cfg.WebInterface.Port)
	}
	fmt.Println("║  ⚠️  Нажмите CTRL+C для завершения работы              ║")
	fmt.Println("║                                                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

	// Ожидание сигнала для завершения
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Закрытие сессии
	s.Close()
}
