package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"discord-bot/config"
)

//go:embed templates
var templatesFS embed.FS

//go:embed static
var staticFS embed.FS

// WebServer представляет веб-сервер для управления ботом
type WebServer struct {
	templates *template.Template
	config    *config.Config
	addr      string
}

// NewWebServer создает новый экземпляр веб-сервера
func NewWebServer(cfg *config.Config) *WebServer {
	// Создаем директории для шаблонов и статических файлов, если они не существуют
	os.MkdirAll("web/templates", 0755)
	os.MkdirAll("web/static", 0755)

	// Создаем базовый шаблон, если он не существует
	if _, err := os.Stat("web/templates/index.html"); os.IsNotExist(err) {
		os.WriteFile("web/templates/index.html", []byte(indexTemplate), 0644)
	}

	// Создаем CSS файл, если он не существует
	if _, err := os.Stat("web/static/style.css"); os.IsNotExist(err) {
		os.WriteFile("web/static/style.css", []byte(styleCSS), 0644)
	}

	// Загружаем шаблоны
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))

	return &WebServer{
		templates: tmpl,
		config:    cfg,
		addr:      fmt.Sprintf("%s:%d", cfg.WebInterface.Host, cfg.WebInterface.Port),
	}
}

// Start запускает веб-сервер
func (ws *WebServer) Start() error {
	// Обработчик для главной страницы
	http.HandleFunc("/", ws.handleIndex)

	// Обработчик для статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Обработчик для сохранения настроек
	http.HandleFunc("/save", ws.handleSaveSettings)

	fmt.Printf("Веб-интерфейс запущен на http://%s\n", ws.addr)
	return http.ListenAndServe(ws.addr, nil)
}

// handleIndex обрабатывает запрос на главную страницу
func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Config *config.Config
	}{
		Config: ws.config,
	}

	ws.templates.ExecuteTemplate(w, "index.html", data)
}

// handleSaveSettings обрабатывает запрос на сохранение настроек
func (ws *WebServer) handleSaveSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из формы
	r.ParseForm()

	// Обновляем настройки бота
	ws.config.Token = r.FormValue("token")
	ws.config.Prefix = r.FormValue("prefix")

	// Сохраняем настройки в файл
	err := ws.config.Save("config.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сохранения настроек: %v", err), http.StatusInternalServerError)
		return
	}

	// Перенаправляем на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Базовый шаблон для главной страницы
const indexTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Настройка Discord бота</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <div class="container">
        <h1>Настройка Discord бота</h1>
        <form action="/save" method="post">
            <div class="form-group">
                <label for="token">Токен бота:</label>
                <input type="text" id="token" name="token" value="{{.Config.Token}}" required>
            </div>
            <div class="form-group">
                <label for="prefix">Префикс команд:</label>
                <input type="text" id="prefix" name="prefix" value="{{.Config.Prefix}}" required>
            </div>
            <button type="submit">Сохранить</button>
        </form>

        <div class="commands-section">
            <h2>Доступные команды</h2>
            <ul>
                <li><strong>{{.Config.Prefix}}nickname @пользователь новый_ник</strong> - изменить никнейм пользователя</li>
                <li><strong>{{.Config.Prefix}}dm @пользователь сообщение</strong> - отправить личное сообщение пользователю</li>
                <li><strong>{{.Config.Prefix}}play URL</strong> - воспроизвести аудио с YouTube</li>
                <li><strong>{{.Config.Prefix}}stop</strong> - остановить воспроизведение</li>
                <li><strong>{{.Config.Prefix}}leave</strong> - выйти из голосового канала</li>
            </ul>
        </div>
    </div>
</body>
</html>`

// CSS стили для веб-интерфейса
const styleCSS = `body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    line-height: 1.6;
    color: #333;
    margin: 0;
    padding: 0;
    background-color: #f5f5f5;
}

.container {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    background-color: white;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    margin-top: 20px;
}

h1 {
    color: #5865F2;
    text-align: center;
    margin-bottom: 30px;
}

h2 {
    color: #5865F2;
    margin-top: 30px;
}

.form-group {
    margin-bottom: 20px;
}

label {
    display: block;
    margin-bottom: 5px;
    font-weight: bold;
}

input[type="text"] {
    width: 100%;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 16px;
}

button {
    background-color: #5865F2;
    color: white;
    border: none;
    padding: 10px 20px;
    font-size: 16px;
    cursor: pointer;
    border-radius: 4px;
    display: block;
    margin: 0 auto;
}

button:hover {
    background-color: #4752c4;
}

.commands-section {
    margin-top: 40px;
    border-top: 1px solid #eee;
    padding-top: 20px;
}

ul {
    list-style-type: none;
    padding: 0;
}

li {
    padding: 10px 0;
    border-bottom: 1px solid #eee;
}

li strong {
    color: #5865F2;
}`
