package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"discord-bot/config"
)

// APIServer представляет API сервер для взаимодействия с фронтендом
type APIServer struct {
	config *config.Config
	addr   string
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
	return &APIServer{
		config: cfg,
		addr:   fmt.Sprintf("%s:%d", cfg.WebInterface.Host, cfg.WebInterface.Port),
	}
}

// Start запускает API сервер
func (api *APIServer) Start() error {
	// Настраиваем CORS для разработки фронтенда
	http.HandleFunc("/api/config", api.handleCORS(api.handleGetConfig))
	http.HandleFunc("/api/save-config", api.handleCORS(api.handleSaveConfig))
	http.HandleFunc("/api/stats", api.handleCORS(api.handleGetStats))
	http.HandleFunc("/api/commands", api.handleCORS(api.handleGetCommands))
	http.HandleFunc("/api/update-command", api.handleCORS(api.handleUpdateCommand))

	// Обслуживаем статические файлы React приложения
	fs := http.FileServer(http.Dir("web/frontend/build"))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Для API запросов не используем файловый сервер
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			return
		}

		// Для всех остальных запросов отдаем index.html
		if r.URL.Path != "/" && !FileExists("web/frontend/build"+r.URL.Path) {
			http.ServeFile(w, r, "web/frontend/build/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	}))

	// Выводим сообщение о запуске API сервера
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                        ║")
	fmt.Println("║  🌐 Веб-панель управления Discord ботом запущена!     ║")
	fmt.Printf("║  📌 Адрес: http://%s                           ║\n", api.addr)
	fmt.Println("║  ⚙️  Откройте эту ссылку в браузере для настройки бота ║")
	fmt.Println("║                                                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

	return http.ListenAndServe(api.addr, nil)
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
	if err := config.Save(api.config); err != nil {
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
