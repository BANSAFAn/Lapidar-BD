package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"discord-bot/config"
)

// APIServer представляет API сервер для управления ботом
type APIServer struct {
	config   *config.Config
	mainAddr string
	altAddrs []string
}

// BotStats представляет статистику бота
type BotStats struct {
	Servers     int    `json:"servers"`
	Users       int    `json:"users"`
	Channels    int    `json:"channels"`
	Commands    int    `json:"commands"`
	Uptime      string `json:"uptime"`
	MemoryUsage string `json:"memoryUsage"`
}

// Command представляет команду бота
type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
	Category    string `json:"category"`
	Enabled     bool   `json:"enabled"`
}

// NewAPIServer создает новый экземпляр API сервера
func NewAPIServer(cfg *config.Config) *APIServer {
	// Создаем основной адрес
	mainAddr := fmt.Sprintf("%s:%d", cfg.WebInterface.Host, cfg.WebInterface.Port)

	// Создаем альтернативные адреса
	altAddrs := make([]string, len(cfg.WebInterface.AltPorts))
	for i, port := range cfg.WebInterface.AltPorts {
		altAddrs[i] = fmt.Sprintf("%s:%d", cfg.WebInterface.Host, port)
	}

	return &APIServer{
		config:   cfg,
		mainAddr: mainAddr,
		altAddrs: altAddrs,
	}
}

// Start запускает API сервер на нескольких портах
func (api *APIServer) Start() error {
	// Создаем роутер
	r := mux.NewRouter()

	// Настраиваем CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Применяем CORS middleware
	handler := c.Handler(r)

	// Регистрируем обработчики API
	r.HandleFunc("/api/config", api.handleGetConfig).Methods("GET")
	r.HandleFunc("/api/config", api.handleSaveConfig).Methods("POST")
	r.HandleFunc("/api/stats", api.handleGetStats).Methods("GET")
	r.HandleFunc("/api/commands", api.handleGetCommands).Methods("GET")
	r.HandleFunc("/api/commands", api.handleUpdateCommands).Methods("POST")

	// Регистрируем обработчики аутентификации
	r.HandleFunc("/api/login", handleLogin).Methods("POST")
	r.HandleFunc("/api/verify-totp", handleVerifyTOTP).Methods("POST")

	// Обслуживаем фронтенд
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("web/frontend/build")))

	// Выводим сообщение о запуске API сервера
	fmt.Printf("API сервер запущен на http://%s (основной)\n", api.mainAddr)
	for i, addr := range api.altAddrs {
		fmt.Printf("API сервер запущен на http://%s (альтернативный %d)\n", addr, i+1)
	}

	// Запускаем основной сервер в отдельной горутине
	go func() {
		err := http.ListenAndServe(api.mainAddr, handler)
		if err != nil {
			fmt.Printf("Ошибка запуска основного API сервера на %s: %v\n", api.mainAddr, err)
		}
	}()

	// Запускаем альтернативные серверы в отдельных горутинах
	for _, addr := range api.altAddrs {
		go func(address string) {
			err := http.ListenAndServe(address, handler)
			if err != nil {
				fmt.Printf("Ошибка запуска альтернативного API сервера на %s: %v\n", address, err)
			}
		}(addr)
	}

	// Ждем бесконечно, чтобы горутины могли работать
	select {}
}

// handleCORS добавляет CORS заголовки к ответам
func (api *APIServer) handleCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

// handleGetConfig возвращает текущую конфигурацию бота
func (api *APIServer) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api.config)
}

// handleSaveConfig сохраняет новую конфигурацию бота
func (api *APIServer) handleSaveConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var newConfig config.Config
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем конфигурацию
	api.config.Token = newConfig.Token
	api.config.Prefix = newConfig.Prefix
	api.config.BotName = newConfig.BotName
	api.config.DefaultLanguage = newConfig.DefaultLanguage
	api.config.WebInterface = newConfig.WebInterface

	// Сохраняем конфигурацию в файл
	if err := config.SaveConfig(api.config); err != nil {
		http.Error(w, "Ошибка сохранения конфигурации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handleGetStats возвращает статистику бота
func (api *APIServer) handleGetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// В реальном приложении здесь будет получение статистики из бота
	// Пока используем тестовые данные
	stats := BotStats{
		Servers:     15,
		Users:       1250,
		Channels:    87,
		Commands:    42,
		Uptime:      "3 дня 7 часов",
		MemoryUsage: "128 MB",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// handleGetCommands возвращает список команд бота
func (api *APIServer) handleGetCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// В реальном приложении здесь будет получение списка команд из бота
	// Пока используем тестовые данные
	commands := []Command{
		{
			Name:        "help",
			Description: "Показывает список доступных команд",
			Usage:       "!help [команда]",
			Category:    "Основные",
			Enabled:     true,
		},
		{
			Name:        "ping",
			Description: "Проверяет задержку бота",
			Usage:       "!ping",
			Category:    "Утилиты",
			Enabled:     true,
		},
		{
			Name:        "ban",
			Description: "Банит пользователя на сервере",
			Usage:       "!ban @пользователь [причина]",
			Category:    "Модерация",
			Enabled:     true,
		},
		{
			Name:        "play",
			Description: "Воспроизводит музыку в голосовом канале",
			Usage:       "!play [ссылка или название]",
			Category:    "Музыка",
			Enabled:     true,
		},
		{
			Name:        "stats",
			Description: "Показывает статистику бота",
			Usage:       "!stats",
			Category:    "Информация",
			Enabled:     true,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commands)
}

// handleUpdateCommand обновляет статус команды
func (api *APIServer) handleUpdateCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var command Command
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// В реальном приложении здесь будет обновление статуса команды в боте
	// Пока просто возвращаем успешный ответ

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// FileExists проверяет существование файла
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
