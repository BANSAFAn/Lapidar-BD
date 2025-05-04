package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/nakagami/firebirdsql"
)

// FirebirdProvider представляет провайдер для работы с Firebird
type FirebirdProvider struct {
	db *sql.DB
}

// Initialize инициализирует соединение с базой данных Firebird
func (p *FirebirdProvider) Initialize(config DatabaseConfig) error {
	// Формируем строку подключения
	connStr := fmt.Sprintf("%s:%s@%s:%d/%s",
		config.User, config.Password, config.Host, config.Port, config.Database)

	// Открываем соединение с базой данных
	db, err := sql.Open("firebirdsql", connStr)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных Firebird: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных Firebird: %w", err)
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
func (p *FirebirdProvider) createTables() error {
	// Таблица для хранения репортов
	_, err := p.db.Exec(`
		CREATE TABLE IF NOT EXISTS reports (
			id INTEGER NOT NULL PRIMARY KEY,
			reported_user_id VARCHAR(255) NOT NULL,
			reporter_id VARCHAR(255) NOT NULL,
			reason VARCHAR(1000) NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			confirmed BOOLEAN DEFAULT FALSE,
			confirmed_by VARCHAR(255)
		)
	`)
	if err != nil {
		return err
	}

	// Создаем генератор последовательности для ID репортов
	_, err = p.db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS reports_id_seq
	`)
	if err != nil {
		return err
	}

	// Таблица для хранения банов
	_, err = p.db.Exec(`
		CREATE TABLE IF NOT EXISTS bans (
			id INTEGER NOT NULL PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			reason VARCHAR(1000) NOT NULL,
			admin_id VARCHAR(255) NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			expires_at TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Создаем генератор последовательности для ID банов
	_, err = p.db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS bans_id_seq
	`)
	if err != nil {
		return err
	}

	return nil
}

// Close закрывает соединение с базой данных
func (p *FirebirdProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// AddReport добавляет новый репорт в базу данных
func (p *FirebirdProvider) AddReport(reportedUserID, reporterID, reason string) (int64, error) {
	// Получаем следующее значение из последовательности
	var nextID int64
	err := p.db.QueryRow("SELECT NEXT VALUE FOR reports_id_seq FROM RDB$DATABASE").Scan(&nextID)
	if err != nil {
		return 0, err
	}

	// Вставляем запись
	_, err = p.db.Exec(
		"INSERT INTO reports (id, reported_user_id, reporter_id, reason, timestamp, confirmed) VALUES (?, ?, ?, ?, ?, FALSE)",
		nextID, reportedUserID, reporterID, reason, time.Now(),
	)
	if err != nil {
		return 0, err
	}

	return nextID, nil
}

// ConfirmReport подтверждает репорт администратором
func (p *FirebirdProvider) ConfirmReport(reportID int64, adminID string) error {
	_, err := p.db.Exec(
		"UPDATE reports SET confirmed = TRUE, confirmed_by = ? WHERE id = ?",
		adminID, reportID,
	)
	return err
}

// GetReportsByUser получает все репорты на указанного пользователя
func (p *FirebirdProvider) GetReportsByUser(userID string) ([]Report, error) {
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
		var timestamp time.Time
		var confirmedBy sql.NullString
		err := rows.Scan(&r.ID, &r.ReportedUserID, &r.ReporterID, &r.Reason, &timestamp, &r.Confirmed, &confirmedBy)
		if err != nil {
			return nil, err
		}

		r.Timestamp = timestamp

		if confirmedBy.Valid {
			r.ConfirmedBy = confirmedBy.String
		}

		reports = append(reports, r)
	}

	return reports, nil
}

// GetReportCount получает количество подтвержденных репортов на пользователя
func (p *FirebirdProvider) GetReportCount(userID string) (int, error) {
	var count int
	err := p.db.QueryRow(
		"SELECT COUNT(*) FROM reports WHERE reported_user_id = ? AND confirmed = TRUE",
		userID,
	).Scan(&count)

	return count, err
}

// AddBan добавляет новый бан в базу данных
func (p *FirebirdProvider) AddBan(userID, reason, adminID string, duration *time.Duration) error {
	// Получаем следующее значение из последовательности
	var nextID int64
	err := p.db.QueryRow("SELECT NEXT VALUE FOR bans_id_seq FROM RDB$DATABASE").Scan(&nextID)
	if err != nil {
		return err
	}

	var expiresAt *time.Time
	if duration != nil {
		expires := time.Now().Add(*duration)
		expiresAt = &expires
	}

	// Вставляем запись
	_, err = p.db.Exec(
		"INSERT INTO bans (id, user_id, reason, admin_id, timestamp, expires_at) VALUES (?, ?, ?, ?, ?, ?)",
		nextID, userID, reason, adminID, time.Now(), expiresAt,
	)

	return err
}

// GetActiveBan проверяет, есть ли активный бан у пользователя
func (p *FirebirdProvider) GetActiveBan(userID string) (*Ban, error) {
	row := p.db.QueryRow(
		"SELECT id, user_id, reason, admin_id, timestamp, expires_at FROM bans WHERE user_id = ? AND (expires_at IS NULL OR expires_at > ?)",
		userID, time.Now(),
	)

	var ban Ban
	var expiresAt sql.NullTime

	err := row.Scan(&ban.ID, &ban.UserID, &ban.Reason, &ban.AdminID, &ban.Timestamp, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		ban.ExpiresAt = &expiresAt.Time
	}

	return &ban, nil
}

// GetType возвращает тип базы данных
func (p *FirebirdProvider) GetType() string {
	return "firebird"
}
