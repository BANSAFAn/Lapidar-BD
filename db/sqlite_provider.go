package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteProvider представляет провайдер для работы с SQLite
type SQLiteProvider struct {
	db *sql.DB
}

// Initialize инициализирует соединение с базой данных SQLite
func (p *SQLiteProvider) Initialize(config DatabaseConfig) error {
	// Формируем строку подключения
	dbPath := config.Database
	if dbPath == "" {
		dbPath = "data/bot.db"
	}

	// Открываем соединение с базой данных
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных SQLite: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных SQLite: %w", err)
	}

	p.db = db

	// Создаем необходимые таблицы
	err = p.createTables()
	if err != nil {
		return fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	return nil
}

// createTables создает необходимые таблицы в базе данных
func (p *SQLiteProvider) createTables() error {
	// Таблица для хранения репортов
	_, err := p.db.Exec(`
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
	_, err = p.db.Exec(`
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

// Close закрывает соединение с базой данных
func (p *SQLiteProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// AddReport добавляет новый репорт в базу данных
func (p *SQLiteProvider) AddReport(reportedUserID, reporterID, reason string) (int64, error) {
	result, err := p.db.Exec(
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
func (p *SQLiteProvider) ConfirmReport(reportID int64, adminID string) error {
	_, err := p.db.Exec(
		"UPDATE reports SET confirmed = TRUE, confirmed_by = ? WHERE id = ?",
		adminID, reportID,
	)
	return err
}

// GetReportsByUser получает все репорты на указанного пользователя
func (p *SQLiteProvider) GetReportsByUser(userID string) ([]Report, error) {
	rows, err := p.db.Query(
		"SELECT id, reported_user_id, reporter_id, reason, timestamp, confirmed, confirmed_by FROM reports WHERE reported_user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []Report
	for rows.Next() {
		var r Report
		var timestamp string
		var confirmedBy sql.NullString
		err := rows.Scan(&r.ID, &r.ReportedUserID, &r.ReporterID, &r.Reason, &timestamp, &r.Confirmed, &confirmedBy)
		if err != nil {
			return nil, err
		}

		t, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, err
		}
		r.Timestamp = t

		if confirmedBy.Valid {
			r.ConfirmedBy = confirmedBy.String
		}

		reports = append(reports, r)
	}

	return reports, nil
}

// GetReportCount получает количество подтвержденных репортов на пользователя
func (p *SQLiteProvider) GetReportCount(userID string) (int, error) {
	var count int
	err := p.db.QueryRow(
		"SELECT COUNT(*) FROM reports WHERE reported_user_id = ? AND confirmed = TRUE",
		userID,
	).Scan(&count)

	return count, err
}

// AddBan добавляет новый бан в базу данных
func (p *SQLiteProvider) AddBan(userID, reason, adminID string, duration *time.Duration) error {
	var expiresAt *time.Time
	if duration != nil {
		expires := time.Now().Add(*duration)
		expiresAt = &expires
	}

	var expiresAtStr interface{} = nil
	if expiresAt != nil {
		expiresAtStr = expiresAt.Format(time.RFC3339)
	}

	_, err := p.db.Exec(
		"INSERT INTO bans (user_id, reason, admin_id, timestamp, expires_at) VALUES (?, ?, ?, ?, ?)",
		userID, reason, adminID, time.Now().Format(time.RFC3339), expiresAtStr,
	)

	return err
}

// GetActiveBan проверяет, есть ли активный бан у пользователя
func (p *SQLiteProvider) GetActiveBan(userID string) (*Ban, error) {
	row := p.db.QueryRow(
		"SELECT id, user_id, reason, admin_id, timestamp, expires_at FROM bans WHERE user_id = ? AND (expires_at IS NULL OR expires_at > ?)",
		userID, time.Now().Format(time.RFC3339),
	)

	var ban Ban
	var timestamp string
	var expiresAt sql.NullString

	err := row.Scan(&ban.ID, &ban.UserID, &ban.Reason, &ban.AdminID, &timestamp, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, err
	}
	ban.Timestamp = t

	if expiresAt.Valid {
		exp, err := time.Parse(time.RFC3339, expiresAt.String)
		if err != nil {
			return nil, err
		}
		ban.ExpiresAt = &exp
	}

	return &ban, nil
}

// GetType возвращает тип базы данных
func (p *SQLiteProvider) GetType() string {
	return "sqlite"
}
