package db

import (
	"fmt"
	"time"

	// Драйверы для различных баз данных
	_ "github.com/go-sql-driver/mysql"  // MySQL
	_ "github.com/jackc/pgx/v5/stdlib"  // PostgreSQL
	_ "github.com/mattn/go-sqlite3"     // SQLite
	_ "github.com/nakagami/firebirdsql" // Firebird
	// MongoDB
	// Расширенный SQL для реляционных БД
)

// DatabaseProvider представляет интерфейс для работы с различными базами данных
type DatabaseProvider interface {
	// Initialize инициализирует соединение с базой данных
	Initialize(config DatabaseConfig) error
	// Close закрывает соединение с базой данных
	Close() error
	// AddReport добавляет новый репорт в базу данных
	AddReport(reportedUserID, reporterID, reason string) (int64, error)
	// ConfirmReport подтверждает репорт администратором
	ConfirmReport(reportID int64, adminID string) error
	// GetReportsByUser получает все репорты на указанного пользователя
	GetReportsByUser(userID string) ([]Report, error)
	// GetReportCount получает количество подтвержденных репортов на пользователя
	GetReportCount(userID string) (int, error)
	// AddBan добавляет новый бан в базу данных
	AddBan(userID, reason, adminID string, duration *time.Duration) error
	// GetActiveBan проверяет, есть ли активный бан у пользователя
	GetActiveBan(userID string) (*Ban, error)
	// GetType возвращает тип базы данных
	GetType() string
}

// DatabaseConfig содержит конфигурацию для подключения к базе данных
type DatabaseConfig struct {
	Type     string `json:"type"` // sqlite, postgres, mysql, mongodb, mariadb, firebird, supabase, triplit
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSL      bool   `json:"ssl"`
	// Дополнительные параметры для конкретных БД
	Params map[string]string `json:"params"`
}

// Report представляет структуру репорта
type Report struct {
	ID             int64
	ReportedUserID string
	ReporterID     string
	Reason         string
	Timestamp      time.Time
	Confirmed      bool
	ConfirmedBy    string
}

// Ban представляет структуру бана
type Ban struct {
	ID        int64
	UserID    string
	Reason    string
	AdminID   string
	Timestamp time.Time
	ExpiresAt *time.Time
}

// AvailableProviders содержит список доступных провайдеров баз данных
var AvailableProviders = map[string]DatabaseProvider{
	"sqlite":   &SQLiteProvider{},
	"postgres": &PostgreSQLProvider{},
	"mysql":    &MySQLProvider{},
	"mongodb":  &MongoDBProvider{},
	"mariadb":  &MariaDBProvider{}, // Использует тот же драйвер, что и MySQL
	"firebird": &FirebirdProvider{},
	"supabase": &PostgreSQLProvider{}, // Supabase использует PostgreSQL
	"triplit":  &TriplitProvider{},
}

// CurrentProvider содержит текущий провайдер базы данных
var CurrentProvider DatabaseProvider

// Initialize инициализирует соединение с базой данных на основе конфигурации
func Initialize(config DatabaseConfig) error {
	// Получаем провайдер из списка доступных
	provider, ok := AvailableProviders[config.Type]
	if !ok {
		return fmt.Errorf("неизвестный тип базы данных: %s", config.Type)
	}

	// Инициализируем провайдер
	err := provider.Initialize(config)
	if err != nil {
		return fmt.Errorf("ошибка инициализации базы данных: %w", err)
	}

	// Устанавливаем текущий провайдер
	CurrentProvider = provider

	return nil
}

// Close закрывает соединение с базой данных
func Close() error {
	if CurrentProvider == nil {
		return nil
	}

	return CurrentProvider.Close()
}
