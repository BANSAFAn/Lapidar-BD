package web

import (
	"discord-bot/config"
	"discord-bot/db"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TOTPRequest представляет запрос с TOTP кодом
type TOTPRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
	Code  string `json:"code"`
}

// LoginResponse представляет ответ на запрос входа
type LoginResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Token      string `json:"token,omitempty"`
	Require2FA bool   `json:"require_2fa,omitempty"`
}

// AuthMiddleware проверяет JWT токен
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Проверяем формат токена
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}

		// Проверяем токен
		claims, err := db.VerifyJWT(parts[1])
		if err != nil {
			http.Error(w, "Невалидный токен: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Проверяем сессию
		sessionID, ok := claims["session_id"].(string)
		if !ok {
			http.Error(w, "Невалидный токен: отсутствует ID сессии", http.StatusUnauthorized)
			return
		}

		// Получаем сессию из базы данных
		session, err := db.GetSession(sessionID)
		if err != nil {
			http.Error(w, "Сессия не найдена", http.StatusUnauthorized)
			return
		}

		// Проверяем срок действия сессии
		if session.ExpiresAt.Before(time.Now()) {
			// Удаляем просроченную сессию
			db.DeleteSession(sessionID)
			http.Error(w, "Сессия истекла", http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в контекст запроса
		r.Header.Set("X-User-Email", session.Email)

		// Вызываем следующий обработчик
		next(w, r)
	}
}

// handleLogin обрабатывает запрос на вход
func (api *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем запрос
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, не заблокирован ли IP
	blocked, err := db.IsLoginBlocked(r.RemoteAddr, req.Email)
	if err != nil {
		http.Error(w, "Ошибка проверки блокировки: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if blocked {
		// Если IP заблокирован, отправляем сообщение о блокировке
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Слишком много неудачных попыток входа. Попробуйте позже.",
		})
		return
	}

	// Загружаем конфигурацию администратора
	adminConfig, err := config.LoadAdminConfig()
	if err != nil {
		http.Error(w, "Ошибка загрузки конфигурации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем email
	if req.Email != adminConfig.Email {
		// Логируем неудачную попытку входа
		db.LogLogin(req.Email, r.RemoteAddr, r.UserAgent(), false, "Неверный email")

		// Добавляем неудачную попытку входа
		blocked, err := db.AddLoginAttempt(r.RemoteAddr, req.Email)
		if err != nil {
			http.Error(w, "Ошибка обновления счетчика попыток: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправляем ответ с задержкой для предотвращения атак перебором
		time.Sleep(time.Second)
		w.Header().Set("Content-Type", "application/json")

		if blocked {
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Слишком много неудачных попыток входа. Попробуйте позже.",
			})
		} else {
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Неверный email или пароль",
			})
		}
		return
	}

	// Проверяем пароль
	if !db.VerifyPassword(adminConfig.Password, req.Password) {
		// Логируем неудачную попытку входа
		db.LogLogin(req.Email, r.RemoteAddr, r.UserAgent(), false, "Неверный пароль")

		// Добавляем неудачную попытку входа
		blocked, err := db.AddLoginAttempt(r.RemoteAddr, req.Email)
		if err != nil {
			http.Error(w, "Ошибка обновления счетчика попыток: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправляем ответ с задержкой для предотвращения атак перебором
		time.Sleep(time.Second)
		w.Header().Set("Content-Type", "application/json")

		if blocked {
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Слишком много неудачных попыток входа. Попробуйте позже.",
			})
		} else {
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Неверный email или пароль",
			})
		}
		return
	}

	// Сбрасываем счетчик попыток входа при успешной аутентификации
	err = db.ResetLoginAttempts(r.RemoteAddr, req.Email)
	if err != nil {
		http.Error(w, "Ошибка сброса счетчика попыток: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем временный токен для второго фактора
	tempToken, err := db.GenerateJWT(req.Email, "temp")
	if err != nil {
		http.Error(w, "Ошибка генерации токена: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Логируем успешную первую стадию входа
	db.LogLogin(req.Email, r.RemoteAddr, r.UserAgent(), true, "Успешная первая стадия входа")

	// Отправляем ответ с требованием второго фактора
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success:    true,
		Message:    "Требуется второй фактор аутентификации",
		Token:      tempToken,
		Require2FA: true,
	})
}

// handleVerifyTOTP обрабатывает запрос на проверку TOTP кода
func (api *APIServer) handleVerifyTOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем запрос
	var req TOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем, не заблокирован ли IP
	blocked, err := db.IsLoginBlocked(r.RemoteAddr, req.Email)
	if err != nil {
		http.Error(w, "Ошибка проверки блокировки: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if blocked {
		// Если IP заблокирован, отправляем сообщение о блокировке
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Success: false,
			Message: "Слишком много неудачных попыток входа. Попробуйте позже.",
		})
		return
	}

	// Проверяем временный токен
	claims, err := db.VerifyJWT(req.Token)
	if err != nil {
		http.Error(w, "Невалидный токен: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Проверяем email в токене
	email, ok := claims["email"].(string)
	if !ok || email != req.Email {
		http.Error(w, "Невалидный токен: неверный email", http.StatusUnauthorized)
		return
	}

	// Загружаем конфигурацию администратора
	adminConfig, err := config.LoadAdminConfig()
	if err != nil {
		http.Error(w, "Ошибка загрузки конфигурации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем TOTP код
	valid := totp.Validate(req.Code, adminConfig.TOTPSecret)
	if !valid {
		// Логируем неудачную попытку входа
		db.LogLogin(req.Email, r.RemoteAddr, r.UserAgent(), false, "Неверный TOTP код")

		// Добавляем неудачную попытку входа
		blocked, err := db.AddLoginAttempt(r.RemoteAddr, req.Email)
		if err != nil {
			http.Error(w, "Ошибка обновления счетчика попыток: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправляем ответ с задержкой для предотвращения атак перебором
		time.Sleep(time.Second)
		w.Header().Set("Content-Type", "application/json")

		if blocked {
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Слишком много неудачных попыток входа. Попробуйте позже.",
			})
		} else {
			json.NewEncoder(w).Encode(LoginResponse{
				Success: false,
				Message: "Неверный код аутентификации",
			})
		}
		return
	}

	// Сбрасываем счетчик попыток входа при успешной аутентификации
	err = db.ResetLoginAttempts(r.RemoteAddr, req.Email)
	if err != nil {
		http.Error(w, "Ошибка сброса счетчика попыток: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем сессию
	session, err := db.CreateSession(req.Email, r.RemoteAddr, r.UserAgent(), time.Hour*24)
	if err != nil {
		http.Error(w, "Ошибка создания сессии: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерируем JWT токен
	token, err := db.GenerateJWT(req.Email, session.ID)
	if err != nil {
		http.Error(w, "Ошибка генерации токена: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Логируем успешный вход
	db.LogLogin(req.Email, r.RemoteAddr, r.UserAgent(), true, "Успешный вход")

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Success: true,
		Message: "Успешный вход",
		Token:   token,
	})
}

// handleLogout обрабатывает запрос на выход
func (api *APIServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем токен из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	// Проверяем формат токена
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
		return
	}

	// Проверяем токен
	claims, err := db.VerifyJWT(parts[1])
	if err != nil {
		http.Error(w, "Невалидный токен: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Получаем ID сессии из токена
	sessionID, ok := claims["session_id"].(string)
	if !ok {
		http.Error(w, "Невалидный токен: отсутствует ID сессии", http.StatusUnauthorized)
		return
	}

	// Удаляем сессию
	err = db.DeleteSession(sessionID)
	if err != nil {
		http.Error(w, "Ошибка удаления сессии: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Успешный выход",
	})
}

// handleGetLoginLogs обрабатывает запрос на получение логов входа
func (api *APIServer) handleGetLoginLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем логи входа
	logs, err := db.GetLoginLogs(100)
	if err != nil {
		http.Error(w, "Ошибка получения логов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GenerateTOTPQRCode генерирует QR-код для TOTP
func GenerateTOTPQRCode(email, secret string) (string, error) {
	// Создаем ключ TOTP
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Lapidar Bot Admin Panel",
		AccountName: email,
		Secret:      []byte(secret),
	})
	if err != nil {
		return "", err
	}

	return key.URL(), nil
}

// handleSetupTOTP обрабатывает запрос на настройку TOTP
func (api *APIServer) handleSetupTOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Загружаем конфигурацию администратора
	adminConfig, err := config.LoadAdminConfig()
	if err != nil {
		http.Error(w, "Ошибка загрузки конфигурации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Генерируем QR-код
	qrURL, err := GenerateTOTPQRCode(adminConfig.Email, adminConfig.TOTPSecret)
	if err != nil {
		http.Error(w, "Ошибка генерации QR-кода: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"qr_url":  qrURL,
		"secret":  adminConfig.TOTPSecret,
	})
}

// UpdateAPIServerForAuth обновляет API сервер для поддержки аутентификации
// UpdateAPIServerForAuth обновляет API сервер для поддержки аутентификации
// Этот метод вызывается после создания API сервера, но до его запуска
func (api *APIServer) UpdateAPIServerForAuth() {
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

	// Регистрируем обработчики аутентификации (без защиты)
	r.HandleFunc("/api/login", api.handleLogin).Methods("POST")
	r.HandleFunc("/api/verify-totp", api.handleVerifyTOTP).Methods("POST")

	// Регистрируем защищенные обработчики
	r.HandleFunc("/api/logout", AuthMiddleware(api.handleLogout)).Methods("POST")
	r.HandleFunc("/api/login-logs", AuthMiddleware(api.handleGetLoginLogs)).Methods("GET")
	r.HandleFunc("/api/setup-totp", AuthMiddleware(api.handleSetupTOTP)).Methods("GET")

	// Защищаем API конфигурации
	r.HandleFunc("/api/config", AuthMiddleware(api.handleGetConfig)).Methods("GET")
	r.HandleFunc("/api/config", AuthMiddleware(api.handleSaveConfig)).Methods("POST")
	r.HandleFunc("/api/stats", AuthMiddleware(api.handleGetStats)).Methods("GET")
	r.HandleFunc("/api/commands", AuthMiddleware(api.handleGetCommands)).Methods("GET")
	r.HandleFunc("/api/commands", AuthMiddleware(api.handleUpdateCommands)).Methods("POST")

	// Обслуживаем фронтенд
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("web/frontend/build")))

	// Выводим сообщение о включении аутентификации
	fmt.Println("✅ Двухфакторная аутентификация для веб-панели включена")

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
}
