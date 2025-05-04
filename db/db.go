package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitializeDB() error {
	// Проверяем существование директории для базы данных
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err := os.Mkdir("data", 0755)
		if err != nil {
			return fmt.Errorf("ошибка создания директории для базы данных: %w", err)
		}
	}

	// Открываем соединение с базой данных
	db, err := sql.Open("sqlite3", "data/bot.db")
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	DB = db

	// Создаем необходимые таблицы
	err = createTables()
	if err != nil {
		return fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	return nil
}

// createTables создает необходимые таблицы в базе данных
func createTables() error {
	// Таблица для хранения репортов
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS reports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			reported_user_id TEXT NOT NULL,
			reporter_id TEXT NOT NULL,
			reason TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			confirmed BOOLEAN DEFAULT FALSE,
			confirmed_by TEXT
		)
	`)
	if err != nil {
		return err
	}

	// Таблица для хранения банов
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS bans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			reason TEXT NOT NULL,
			admin_id TEXT NOT NULL,
			timestamp DATETIME NOT NULL,
			expires_at DATETIME
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// AddReport добавляет новый репорт в базу данных
func AddReport(reportedUserID, reporterID, reason string) (int64, error) {
	result, err := DB.Exec(
		"INSERT INTO reports (reported_user_id, reporter_id, reason, timestamp) VALUES (?, ?, ?, ?)",
		reportedUserID, reporterID, reason, time.Now(),
	)
	if err != nil {
		return 0, err
	}

	returnID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return returnID, nil
}

// ConfirmReport подтверждает репорт администратором
func ConfirmReport(reportID int64, adminID string) error {
	_, err := DB.Exec(
		"UPDATE reports SET confirmed = TRUE, confirmed_by = ? WHERE id = ?",
		adminID, reportID,
	)
	return err
}

// GetReportCount возвращает количество подтвержденных репортов на пользователя
func GetReportCount(userID string) (int, error) {
	var count int
	err := DB.QueryRow(
		"SELECT COUNT(*) FROM reports WHERE reported_user_id = ? AND confirmed = TRUE",
		userID,
	).Scan(&count)

	return count, err
}

// AddBan добавляет запись о бане пользователя
func AddBan(userID, reason, adminID string, duration *time.Duration) error {
	var expiresAt *time.Time
	if duration != nil {
		expires := time.Now().Add(*duration)
		expiresAt = &expires
	}

	var expiresAtStr interface{} = nil
	if expiresAt != nil {
		expiresAtStr = expiresAt.Format(time.RFC3339)
	}

	_, err := DB.Exec(
		"INSERT INTO bans (user_id, reason, admin_id, timestamp, expires_at) VALUES (?, ?, ?, ?, ?)",
		userID, reason, adminID, time.Now(), expiresAtStr,
	)

	return err
}

// IsUserBanned проверяет, забанен ли пользователь
func IsUserBanned(userID string) (bool, error) {
	var count int
	err := DB.QueryRow(
		"SELECT COUNT(*) FROM bans WHERE user_id = ? AND (expires_at IS NULL OR expires_at > ?)",
		userID, time.Now().Format(time.RFC3339),
	).Scan(&count)

	return count > 0, err
}
