package db

import (
	"discord-bot/config"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Структура для хранения информации о сессии
type Session struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Структура для хранения информации о логе входа
type LoginLog struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
}

// Структура для хранения информации о попытках входа
type LoginAttempt struct {
	IP        string    `json:"ip"`
	Email     string    `json:"email"`
	Attempts  int       `json:"attempts"`
	LastTry   time.Time `json:"last_try"`
	Blocked   bool      `json:"blocked"`
	BlockedAt time.Time `json:"blocked_at"`
}

// InitAuthTables инициализирует таблицы для аутентификации
func InitAuthTables() error {
	// Таблица для хранения сессий
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			ip TEXT NOT NULL,
			user_agent TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			expires_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Таблица для хранения логов входа
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS login_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			ip TEXT NOT NULL,
			user_agent TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			success BOOLEAN NOT NULL,
			message TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Таблица для хранения попыток входа (для защиты от брутфорса)
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS login_attempts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT NOT NULL,
			email TEXT NOT NULL,
			attempts INTEGER NOT NULL DEFAULT 0,
			last_try DATETIME NOT NULL,
			blocked BOOLEAN NOT NULL DEFAULT 0,
			blocked_at DATETIME,
			UNIQUE(ip, email)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// CreateSession создает новую сессию
func CreateSession(email, ip, userAgent string, duration time.Duration) (*Session, error) {
	session := &Session{
		ID:        GenerateRandomString(32),
		Email:     email,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	_, err := DB.Exec(
		"INSERT INTO sessions (id, email, ip, user_agent, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?)",
		session.ID, session.Email, session.IP, session.UserAgent, session.CreatedAt, session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession получает сессию по ID
func GetSession(sessionID string) (*Session, error) {
	session := &Session{}
	err := DB.QueryRow(
		"SELECT id, email, ip, user_agent, created_at, expires_at FROM sessions WHERE id = ?",
		sessionID,
	).Scan(&session.ID, &session.Email, &session.IP, &session.UserAgent, &session.CreatedAt, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// DeleteSession удаляет сессию по ID
func DeleteSession(sessionID string) error {
	_, err := DB.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	return err
}

// CleanExpiredSessions удаляет просроченные сессии
func CleanExpiredSessions() error {
	_, err := DB.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
}

// LogLogin записывает информацию о попытке входа
func LogLogin(email, ip, userAgent string, success bool, message string) error {
	_, err := DB.Exec(
		"INSERT INTO login_logs (email, ip, user_agent, timestamp, success, message) VALUES (?, ?, ?, ?, ?, ?)",
		email, ip, userAgent, time.Now(), success, message,
	)
	return err
}

// GetLoginLogs получает логи входа
func GetLoginLogs(limit int) ([]LoginLog, error) {
	rows, err := DB.Query(
		"SELECT id, email, ip, user_agent, timestamp, success, message FROM login_logs ORDER BY timestamp DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LoginLog
	for rows.Next() {
		var log LoginLog
		err := rows.Scan(&log.ID, &log.Email, &log.IP, &log.UserAgent, &log.Timestamp, &log.Success, &log.Message)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// VerifyPassword проверяет пароль
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// HashPassword хеширует пароль с использованием bcrypt
// Используется более высокий уровень стоимости для повышения безопасности
func HashPassword(password string) (string, error) {
	// Используем более высокий уровень стоимости для повышения безопасности
	// DefaultCost = 10, MaxCost = 31, рекомендуется использовать 12-14 для хорошего баланса
	cost := 12
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// IsPasswordStrong проверяет сложность пароля
// Пароль должен содержать минимум 8 символов, включая цифры, строчные и заглавные буквы, и специальные символы
func IsPasswordStrong(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Пароль должен содержать минимум 8 символов"
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*' || char == '(' || char == ')' || char == '-' || char == '_' || char == '+' || char == '=' || char == '{' || char == '}' || char == '[' || char == ']' || char == '|' || char == '\\' || char == ':' || char == ';' || char == '"' || char == '\'' || char == '<' || char == '>' || char == ',' || char == '.' || char == '?' || char == '/':
			hasSpecial = true
		}
	}

	if !hasLower {
		return false, "Пароль должен содержать хотя бы одну строчную букву"
	}
	if !hasUpper {
		return false, "Пароль должен содержать хотя бы одну заглавную букву"
	}
	if !hasDigit {
		return false, "Пароль должен содержать хотя бы одну цифру"
	}
	if !hasSpecial {
		return false, "Пароль должен содержать хотя бы один специальный символ"
	}

	return true, ""
}

// GenerateJWT генерирует JWT токен с улучшенной безопасностью
func GenerateJWT(email string, sessionID string) (string, error) {
	// Загружаем конфигурацию администратора для получения секретного ключа
	adminConfig, err := config.LoadAdminConfig()
	if err != nil {
		return "", err
	}

	// Создаем новый токен
	token := jwt.New(jwt.SigningMethodHS256)

	// Получаем текущее время для расчета времени истечения
	now := time.Now()
	expTime := now.Add(time.Hour * 24) // Токен действителен 24 часа

	// Устанавливаем стандартные и дополнительные claims для повышения безопасности
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["session_id"] = sessionID
	claims["iat"] = now.Unix()                // Время выпуска токена
	claims["exp"] = expTime.Unix()           // Время истечения токена
	claims["nbf"] = now.Unix()               // Время, до которого токен недействителен
	claims["jti"] = GenerateRandomString(16) // Уникальный идентификатор токена

	// Добавляем информацию о пользователе и устройстве
	claims["type"] = "access"                // Тип токена

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(adminConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT проверяет JWT токен с расширенной проверкой безопасности
func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	// Загружаем конфигурацию администратора для получения секретного ключа
	adminConfig, err := config.LoadAdminConfig()
	if err != nil {
		return nil, err
	}

	// Парсим токен с проверкой всех стандартных claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}

		// Проверяем наличие необходимых claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("невалидные claims")
		}

		// Проверяем тип токена
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" && tokenType != "temp" {
			return nil, fmt.Errorf("неверный тип токена")
		}

		// Проверяем наличие email и session_id (кроме временных токенов)
		if tokenType == "access" {
			if _, ok := claims["email"].(string); !ok {
				return nil, fmt.Errorf("отсутствует email")
			}
			if _, ok := claims["session_id"].(string); !ok {
				return nil, fmt.Errorf("отсутствует session_id")
			}
		}

		return []byte(adminConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Дополнительная проверка времени действия токена
		now := time.Now().Unix()

		// Проверяем, не истек ли срок действия токена
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < now {
				return nil, fmt.Errorf("токен истек")
			}
		}

		// Проверяем, не используется ли токен раньше времени
		if nbf, ok := claims["nbf"].(float64); ok {
			if now < int64(nbf) {
				return nil, fmt.Errorf("токен еще не действителен")
			}
		}

		return claims, nil
	}

	return nil, fmt.Errorf("невалидный токен")
}

// GenerateRandomString генерирует случайную строку заданной длины
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[RandomInt(0, len(charset)-1)]
	}
	return string(b)
}

// RandomInt генерирует случайное число в заданном диапазоне
func RandomInt(min, max int) int {
	return min + time.Now().Nanosecond()%(max-min+1)
}

// Константы для настройки защиты от брутфорса
const (
	MAX_LOGIN_ATTEMPTS = 5         // Максимальное количество попыток входа
	BLOCK_DURATION     = time.Minute * 15 // Длительность блокировки
	ATTEMPT_RESET_TIME = time.Hour * 1   // Время сброса счетчика попыток
)

// AddLoginAttempt добавляет попытку входа и проверяет, не превышен ли лимит
func AddLoginAttempt(ip, email string) (bool, error) {
	// Получаем текущее время
	now := time.Now()

	// Проверяем, существует ли запись для данного IP и email
	var id int64
	var attempts int
	var lastTry time.Time
	var blocked bool
	var blockedAt time.Time

	err := DB.QueryRow(
		"SELECT id, attempts, last_try, blocked, blocked_at FROM login_attempts WHERE ip = ? AND email = ?",
		ip, email,
	).Scan(&id, &attempts, &lastTry, &blocked, &blockedAt)

	// Если запись не найдена, создаем новую
	if err == sql.ErrNoRows {
		_, err := DB.Exec(
			"INSERT INTO login_attempts (ip, email, attempts, last_try, blocked) VALUES (?, ?, 1, ?, 0)",
			ip, email, now,
		)
		if err != nil {
			return false, err
		}
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Если пользователь заблокирован, проверяем, не истекло ли время блокировки
	if blocked {
		if now.Sub(blockedAt) < BLOCK_DURATION {
			// Блокировка еще действует
			return true, nil
		} else {
			// Время блокировки истекло, сбрасываем счетчик
			_, err := DB.Exec(
				"UPDATE login_attempts SET attempts = 1, last_try = ?, blocked = 0, blocked_at = NULL WHERE id = ?",
				now, id,
			)
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}

	// Если прошло достаточно времени с последней попытки, сбрасываем счетчик
	if now.Sub(lastTry) > ATTEMPT_RESET_TIME {
		_, err := DB.Exec(
			"UPDATE login_attempts SET attempts = 1, last_try = ? WHERE id = ?",
			now, id,
		)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// Увеличиваем счетчик попыток
	attempts++

	// Если превышен лимит попыток, блокируем пользователя
	if attempts >= MAX_LOGIN_ATTEMPTS {
		_, err := DB.Exec(
			"UPDATE login_attempts SET attempts = ?, last_try = ?, blocked = 1, blocked_at = ? WHERE id = ?",
			attempts, now, now, id,
		)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Обновляем счетчик попыток
	_, err = DB.Exec(
		"UPDATE login_attempts SET attempts = ?, last_try = ? WHERE id = ?",
		attempts, now, id,
	)
	if err != nil {
		return false, err
	}

	return false, nil
}

// IsLoginBlocked проверяет, заблокирован ли вход для данного IP и email
func IsLoginBlocked(ip, email string) (bool, error) {
	var blocked bool
	var blockedAt time.Time

	err := DB.QueryRow(
		"SELECT blocked, blocked_at FROM login_attempts WHERE ip = ? AND email = ?",
		ip, email,
	).Scan(&blocked, &blockedAt)

	if err == sql.ErrNoRows {
		// Если записи нет, значит пользователь не заблокирован
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Если пользователь заблокирован, проверяем, не истекло ли время блокировки
	if blocked {
		if time.Now().Sub(blockedAt) < BLOCK_DURATION {
			// Блокировка еще действует
			return true, nil
		} else {
			// Время блокировки истекло, сбрасываем блокировку
			_, err := DB.Exec(
				"UPDATE login_attempts SET attempts = 0, blocked = 0, blocked_at = NULL WHERE ip = ? AND email = ?",
				ip, email,
			)
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}

	return false, nil
}

// ResetLoginAttempts сбрасывает счетчик попыток входа для данного IP и email
func ResetLoginAttempts(ip, email string) error {
	_, err := DB.Exec(
		"UPDATE login_attempts SET attempts = 0, blocked = 0, blocked_at = NULL WHERE ip = ? AND email = ?",
		ip, email,
	)
	return err
}
