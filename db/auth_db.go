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

// HashPassword хеширует пароль
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// GenerateJWT генерирует JWT токен
func GenerateJWT(email string, sessionID string) (string, error) {
	// Загружаем конфигурацию администратора для получения секретного ключа
	adminConfig, err := config.LoadAdminConfig()
	if err != nil {
		return "", err
	}

	// Создаем новый токен
	token := jwt.New(jwt.SigningMethodHS256)

	// Устанавливаем claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["session_id"] = sessionID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Токен действителен 24 часа

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(adminConfig.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyJWT проверяет JWT токен
func VerifyJWT(tokenString string) (jwt.MapClaims, error) {
	// Загружаем конфигурацию администратора для получения секретного ключа
	adminConfig, err := config.LoadAdminConfig()
	if err != nil {
		return nil, err
	}

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}

		return []byte(adminConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
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
