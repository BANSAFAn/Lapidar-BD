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

	fmt.Printf("%s успешно запущен с префиксом '%s'. Нажмите CTRL+C для выхода.\n", cfg.BotName, cfg.Prefix)

	// Ожидание сигнала для завершения
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Закрытие сессии
	s.Close()
}
