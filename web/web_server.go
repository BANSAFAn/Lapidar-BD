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

	// Выводим красивое сообщение о запуске веб-интерфейса
	fmt.Println("\n╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                        ║")
	fmt.Println("║  🌐 Веб-панель управления Discord ботом запущена!     ║")
	fmt.Printf("║  📌 Адрес: http://%s                           ║\n", ws.addr)
	fmt.Println("║  ⚙️  Откройте эту ссылку в браузере для настройки бота ║")
	fmt.Println("║                                                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝\n")

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
	ws.config.BotName = r.FormValue("botname")
	ws.config.DefaultLanguage = r.FormValue("language")

	// Сохраняем настройки в файл
	err := config.SaveConfig(ws.config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сохранения настроек: %v", err), http.StatusInternalServerError)
		return
	}

	// Добавляем сообщение об успешном сохранении
	fmt.Println("✅ Настройки бота успешно обновлены через веб-интерфейс")

	// Перенаправляем на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Базовый шаблон для главной страницы
const indexTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Панель управления Discord ботом</title>
    <link rel="stylesheet" href="/static/style.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
</head>
<body>
    <div class="sidebar">
        <div class="logo">
            <i class="fab fa-discord"></i>
            <span>Lapidar Bot</span>
        </div>
        <nav>
            <ul>
                <li class="active"><a href="#"><i class="fas fa-cogs"></i> Настройки</a></li>
                <li><a href="#commands"><i class="fas fa-terminal"></i> Команды</a></li>
                <li><a href="#status"><i class="fas fa-chart-line"></i> Статус</a></li>
            </ul>
        </nav>
        <div class="bot-status online">
            <span class="status-dot"></span>
            <span class="status-text">Бот онлайн</span>
        </div>
    </div>

    <div class="main-content">
        <header>
            <h1><i class="fas fa-robot"></i> Панель управления Discord ботом</h1>
        </header>

        <div class="card" id="settings">
            <div class="card-header">
                <h2><i class="fas fa-cogs"></i> Основные настройки</h2>
            </div>
            <div class="card-body">
                <form action="/save" method="post" id="settings-form">
                    <div class="form-group">
                        <label for="token"><i class="fas fa-key"></i> Токен бота:</label>
                        <div class="input-group">
                            <input type="password" id="token" name="token" value="{{.Config.Token}}" required>
                            <button type="button" class="toggle-password" onclick="togglePassword('token')">
                                <i class="fas fa-eye"></i>
                            </button>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="prefix"><i class="fas fa-terminal"></i> Префикс команд:</label>
                        <input type="text" id="prefix" name="prefix" value="{{.Config.Prefix}}" required>
                    </div>
                    <div class="form-group">
                        <label for="botname"><i class="fas fa-tag"></i> Имя бота:</label>
                        <input type="text" id="botname" name="botname" value="{{.Config.BotName}}">
                    </div>
                    <div class="form-group">
                        <label for="language"><i class="fas fa-language"></i> Язык бота:</label>
                        <select id="language" name="language">
                            <option value="ru" {{if eq .Config.DefaultLanguage "ru"}}selected{{end}}>Русский</option>
                            <option value="en" {{if eq .Config.DefaultLanguage "en"}}selected{{end}}>English</option>
                            <option value="uk" {{if eq .Config.DefaultLanguage "uk"}}selected{{end}}>Українська</option>
                            <option value="de" {{if eq .Config.DefaultLanguage "de"}}selected{{end}}>Deutsch</option>
                            <option value="zh" {{if eq .Config.DefaultLanguage "zh"}}selected{{end}}>中文</option>
                        </select>
                    </div>
                    <button type="submit" class="btn-primary"><i class="fas fa-save"></i> Сохранить настройки</button>
                </form>
            </div>
        </div>

        <div class="card" id="commands">
            <div class="card-header">
                <h2><i class="fas fa-terminal"></i> Доступные команды</h2>
            </div>
            <div class="card-body">
                <div class="commands-list">
                    <div class="command-item">
                        <div class="command-name"><i class="fas fa-user-edit"></i> {{.Config.Prefix}}nickname @пользователь новый_ник</div>
                        <div class="command-desc">Изменить никнейм пользователя</div>
                    </div>
                    <div class="command-item">
                        <div class="command-name"><i class="fas fa-envelope"></i> {{.Config.Prefix}}dm @пользователь сообщение</div>
                        <div class="command-desc">Отправить личное сообщение пользователю</div>
                    </div>
                    <div class="command-item">
                        <div class="command-name"><i class="fas fa-play-circle"></i> {{.Config.Prefix}}play URL</div>
                        <div class="command-desc">Воспроизвести аудио с YouTube</div>
                    </div>
                    <div class="command-item">
                        <div class="command-name"><i class="fas fa-stop-circle"></i> {{.Config.Prefix}}stop</div>
                        <div class="command-desc">Остановить воспроизведение</div>
                    </div>
                    <div class="command-item">
                        <div class="command-name"><i class="fas fa-sign-out-alt"></i> {{.Config.Prefix}}leave</div>
                        <div class="command-desc">Выйти из голосового канала</div>
                    </div>
                </div>
            </div>
        </div>

        <div class="card" id="status">
            <div class="card-header">
                <h2><i class="fas fa-chart-line"></i> Статус бота</h2>
            </div>
            <div class="card-body">
                <div class="status-grid">
                    <div class="status-item">
                        <div class="status-icon"><i class="fas fa-server"></i></div>
                        <div class="status-info">
                            <div class="status-label">Статус сервера</div>
                            <div class="status-value online">Онлайн</div>
                        </div>
                    </div>
                    <div class="status-item">
                        <div class="status-icon"><i class="fas fa-clock"></i></div>
                        <div class="status-info">
                            <div class="status-label">Время работы</div>
                            <div class="status-value" id="uptime">Загрузка...</div>
                        </div>
                    </div>
                    <div class="status-item">
                        <div class="status-icon"><i class="fas fa-memory"></i></div>
                        <div class="status-info">
                            <div class="status-label">Использование памяти</div>
                            <div class="status-value" id="memory-usage">Загрузка...</div>
                        </div>
                    </div>
                    <div class="status-item">
                        <div class="status-icon"><i class="fas fa-users"></i></div>
                        <div class="status-info">
                            <div class="status-label">Активные серверы</div>
                            <div class="status-value" id="active-servers">Загрузка...</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <footer>
            <p>&copy; 2023 Lapidar Bot. Все права защищены.</p>
        </footer>
    </div>

    <script>
        // Функция для переключения видимости пароля
        function togglePassword(inputId) {
            const input = document.getElementById(inputId);
            const icon = input.nextElementSibling.querySelector('i');
            
            if (input.type === 'password') {
                input.type = 'text';
                icon.classList.remove('fa-eye');
                icon.classList.add('fa-eye-slash');
            } else {
                input.type = 'password';
                icon.classList.remove('fa-eye-slash');
                icon.classList.add('fa-eye');
            }
        }

        // Имитация обновления данных статуса
        function updateStatus() {
            document.getElementById('uptime').textContent = getRandomUptime();
            document.getElementById('memory-usage').textContent = getRandomMemory();
            document.getElementById('active-servers').textContent = getRandomServers();
            setTimeout(updateStatus, 5000);
        }

        function getRandomUptime() {
            const hours = Math.floor(Math.random() * 24);
            const minutes = Math.floor(Math.random() * 60);
            return hours + ' ч ' + minutes + ' мин';
        }

        function getRandomMemory() {
            return Math.floor(Math.random() * 100) + ' МБ';
        }

        function getRandomServers() {
            return Math.floor(Math.random() * 10) + 1;
        }

        // Запуск обновления статуса
        document.addEventListener('DOMContentLoaded', function() {
            updateStatus();
        });

        // Плавная прокрутка к разделам
        document.querySelectorAll('nav a').forEach(anchor => {
            anchor.addEventListener('click', function(e) {
                e.preventDefault();
                const targetId = this.getAttribute('href');
                if(targetId !== '#') {
                    const targetElement = document.querySelector(targetId);
                    window.scrollTo({
                        top: targetElement.offsetTop - 20,
                        behavior: 'smooth'
                    });
                }
                
                // Активный элемент меню
                document.querySelectorAll('nav li').forEach(item => {
                    item.classList.remove('active');
                });
                this.parentElement.classList.add('active');
            });
        });
    </script>
</body>
</html>`

// CSS стили для веб-интерфейса
const styleCSS = `/* Основные стили */
:root {
    --primary-color: #5865F2;
    --primary-dark: #4752c4;
    --secondary-color: #2D3748;
    --accent-color: #EB459E;
    --light-color: #F8F9FA;
    --dark-color: #1A202C;
    --success-color: #48BB78;
    --warning-color: #F6AD55;
    --danger-color: #F56565;
    --border-radius: 8px;
    --box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.08);
    --transition: all 0.3s ease;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    line-height: 1.6;
    color: var(--dark-color);
    background-color: #f0f2f5;
    display: flex;
    min-height: 100vh;
}

/* Боковая панель */
.sidebar {
    width: 250px;
    background-color: var(--secondary-color);
    color: white;
    padding: 20px 0;
    height: 100vh;
    position: fixed;
    left: 0;
    top: 0;
    overflow-y: auto;
    transition: var(--transition);
    z-index: 1000;
}

.logo {
    display: flex;
    align-items: center;
    padding: 0 20px;
    margin-bottom: 30px;
}

.logo i {
    font-size: 28px;
    color: var(--primary-color);
    margin-right: 10px;
}

.logo span {
    font-size: 20px;
    font-weight: bold;
}

nav ul {
    list-style: none;
}

nav li {
    margin-bottom: 5px;
    border-left: 3px solid transparent;
    transition: var(--transition);
}

nav li.active {
    border-left-color: var(--primary-color);
    background-color: rgba(255, 255, 255, 0.1);
}

nav a {
    display: flex;
    align-items: center;
    padding: 12px 20px;
    color: white;
    text-decoration: none;
    transition: var(--transition);
}

nav a:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

nav i {
    margin-right: 10px;
    font-size: 18px;
}

.bot-status {
    margin-top: auto;
    padding: 15px 20px;
    display: flex;
    align-items: center;
    background-color: rgba(0, 0, 0, 0.2);
    margin: 20px;
    border-radius: var(--border-radius);
}

.status-dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    margin-right: 10px;
}

.online .status-dot {
    background-color: var(--success-color);
    box-shadow: 0 0 5px var(--success-color);
}

.offline .status-dot {
    background-color: var(--danger-color);
    box-shadow: 0 0 5px var(--danger-color);
}

/* Основной контент */
.main-content {
    flex: 1;
    margin-left: 250px;
    padding: 20px;
    max-width: 1200px;
    width: 100%;
}

header {
    margin-bottom: 30px;
    padding-bottom: 20px;
    border-bottom: 1px solid #eee;
}

header h1 {
    color: var(--primary-color);
    font-size: 28px;
    display: flex;
    align-items: center;
}

header h1 i {
    margin-right: 10px;
}

/* Карточки */
.card {
    background-color: white;
    border-radius: var(--border-radius);
    box-shadow: var(--box-shadow);
    margin-bottom: 30px;
    overflow: hidden;
    transition: var(--transition);
}

.card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.12), 0 4px 8px rgba(0, 0, 0, 0.06);
}

.card-header {
    background-color: var(--light-color);
    padding: 15px 20px;
    border-bottom: 1px solid #eee;
}

.card-header h2 {
    color: var(--primary-color);
    font-size: 20px;
    margin: 0;
    display: flex;
    align-items: center;
}

.card-header h2 i {
    margin-right: 10px;
}

.card-body {
    padding: 20px;
}

/* Формы */
.form-group {
    margin-bottom: 20px;
}

label {
    display: block;
    margin-bottom: 8px;
    font-weight: 600;
    color: var(--secondary-color);
}

label i {
    margin-right: 5px;
    color: var(--primary-color);
}

input[type="text"],
input[type="password"],
select {
    width: 100%;
    padding: 12px;
    border: 1px solid #ddd;
    border-radius: var(--border-radius);
    font-size: 16px;
    transition: var(--transition);
}

input[type="text"]:focus,
input[type="password"]:focus,
select:focus {
    border-color: var(--primary-color);
    box-shadow: 0 0 0 3px rgba(88, 101, 242, 0.2);
    outline: none;
}

.input-group {
    position: relative;
    display: flex;
}

.input-group input {
    flex: 1;
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
}

.toggle-password {
    background-color: #f1f1f1;
    border: 1px solid #ddd;
    border-left: none;
    padding: 0 15px;
    border-top-right-radius: var(--border-radius);
    border-bottom-right-radius: var(--border-radius);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
}

.toggle-password:hover {
    background-color: #e9e9e9;
}

.btn-primary {
    background-color: var(--primary-color);
    color: white;
    border: none;
    padding: 12px 24px;
    font-size: 16px;
    cursor: pointer;
    border-radius: var(--border-radius);
    display: inline-block;
    transition: var(--transition);
    font-weight: 600;
}

.btn-primary:hover {
    background-color: var(--primary-dark);
    transform: translateY(-2px);
}

/* Команды */
.commands-list {
    display: grid;
    gap: 15px;
}

.command-item {
    background-color: var(--light-color);
    border-radius: var(--border-radius);
    padding: 15px;
    transition: var(--transition);
    border-left: 3px solid var(--primary-color);
}

.command-item:hover {
    transform: translateX(5px);
    box-shadow: var(--box-shadow);
}

.command-name {
    font-weight: bold;
    margin-bottom: 5px;
    color: var(--primary-color);
    display: flex;
    align-items: center;
}

.command-name i {
    margin-right: 8px;
}

.command-desc {
    color: var(--secondary-color);
    font-size: 14px;
}

/* Статус */
.status-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 20px;
}

.status-item {
    display: flex;
    align-items: center;
    background-color: var(--light-color);
    padding: 15px;
    border-radius: var(--border-radius);
    transition: var(--transition);
}

.status-item:hover {
    transform: translateY(-3px);
    box-shadow: var(--box-shadow);
}

.status-icon {
    width: 50px;
    height: 50px;
    background-color: rgba(88, 101, 242, 0.1);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 15px;
    color: var(--primary-color);
    font-size: 20px;
}

.status-info {
    flex: 1;
}

.status-label {
    font-size: 14px;
    color: var(--secondary-color);
    margin-bottom: 5px;
}

.status-value {
    font-weight: bold;
    font-size: 16px;
}

.status-value.online {
    color: var(--success-color);
}

.status-value.offline {
    color: var(--danger-color);
}

/* Подвал */
footer {
    margin-top: 40px;
    padding-top: 20px;
    border-top: 1px solid #eee;
    text-align: center;
    color: var(--secondary-color);
    font-size: 14px;
}

/* Адаптивность */
@media (max-width: 768px) {
    .sidebar {
        width: 70px;
        padding: 10px 0;
    }
    
    .logo span,
    nav a span,
    .bot-status .status-text {
        display: none;
    }
    
    .logo {
        justify-content: center;
        padding: 0;
    }
    
    .logo i {
        margin-right: 0;
    }
    
    nav a {
        justify-content: center;
        padding: 15px;
    }
    
    nav i {
        margin-right: 0;
        font-size: 20px;
    }
    
    .main-content {
        margin-left: 70px;
    }
    
    .status-grid {
        grid-template-columns: 1fr;
    }
}

@media (max-width: 480px) {
    .main-content {
        padding: 15px;
    }
    
    .card-body {
        padding: 15px;
    }
    
    .btn-primary {
        width: 100%;
    }
}`
