package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgreSQLProvider представляет провайдер для работы с PostgreSQL
type PostgreSQLProvider struct {
	db *sql.DB
}

// Initialize инициализирует соединение с базой данных PostgreSQL
func (p *PostgreSQLProvider) Initialize(config DatabaseConfig) error {
	// Формируем строку подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.User, config.Password, config.Host, config.Port, config.Database)

	// Добавляем параметр SSL, если необходимо
	if !config.SSL {
		connStr += "?sslmode=disable"
	}

	// Открываем соединение с базой данных
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных PostgreSQL: %w", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных PostgreSQL: %w", err)
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
func (p *PostgreSQLProvider) createTables() error {
	// Таблица для хранения репортов
	_, err := p.db.Exec(`
		CREATE TABLE IF NOT EXISTS reports (
			id SERIAL PRIMARY KEY,
			reported_user_id TEXT NOT NULL,
			reporter_id TEXT NOT NULL,
			reason TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL,
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
			id SERIAL PRIMARY KEY,
			user_id TEXT NOT NULL,
			reason TEXT NOT NULL,
			admin_id TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			expires_at TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// Close закрывает соединение с базой данных
func (p *PostgreSQLProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// AddReport добавляет новый репорт в базу данных
func (p *PostgreSQLProvider) AddReport(reportedUserID, reporterID, reason string) (int64, error) {
	var id int64
	err := p.db.QueryRow(
		"INSERT INTO reports (reported_user_id, reporter_id, reason, timestamp) VALUES ($1, $2, $3, $4) RETURNING id",
		reportedUserID, reporterID, reason, time.Now(),
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ConfirmReport подтверждает репорт администратором
func (p *PostgreSQLProvider) ConfirmReport(reportID int64, adminID string) error {
	_, err := p.db.Exec(
		"UPDATE reports SET confirmed = TRUE, confirmed_by = $1 WHERE id = $2",
		adminID, reportID,
	)
	return err
}

// GetReportsByUser получает все репорты на указанного пользователя
func (p *PostgreSQLProvider) GetReportsByUser(userID string) ([]Report, error) {
	rows, err := p.db.Query(
		"SELECT id, reported_user_id, reporter_id, reason, timestamp, confirmed, confirmed_by FROM reports WHERE reported_user_id = $1",
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
func (p *PostgreSQLProvider) GetReportCount(userID string) (int, error) {
	var count int
	err := p.db.QueryRow(
		"SELECT COUNT(*) FROM reports WHERE reported_user_id = $1 AND confirmed = TRUE",
		userID,
	).Scan(&count)

	return count, err
}

// AddBan добавляет новый бан в базу данных
func (p *PostgreSQLProvider) AddBan(userID, reason, adminID string, duration *time.Duration) error {
	var expiresAt *time.Time
	if duration != nil {
		expires := time.Now().Add(*duration)
		expiresAt = &expires
	}

	_, err := p.db.Exec(
		"INSERT INTO bans (user_id, reason, admin_id, timestamp, expires_at) VALUES ($1, $2, $3, $4, $5)",
		userID, reason, adminID, time.Now(), expiresAt,
	)

	return err
}

// GetActiveBan проверяет, есть ли активный бан у пользователя
func (p *PostgreSQLProvider) GetActiveBan(userID string) (*Ban, error) {
	row := p.db.QueryRow(
		"SELECT id, user_id, reason, admin_id, timestamp, expires_at FROM bans WHERE user_id = $1 AND (expires_at IS NULL OR expires_at > $2)",
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
func (p *PostgreSQLProvider) GetType() string {
	return "postgres"
}
