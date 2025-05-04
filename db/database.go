package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/nakagami/firebirdsql"
)

type DatabaseProvider interface {
	Initialize(config DatabaseConfig) error
	Close() error
	AddReport(reportedUserID, reporterID, reason string) (int64, error)
	ConfirmReport(reportID int64, adminID string) error
	GetReportsByUser(userID string) ([]Report, error)
	GetReportCount(userID string) (int, error)
	AddBan(userID, reason, adminID string, duration *time.Duration) error
	GetActiveBan(userID string) (*Ban, error)
	GetType() string
}

type DatabaseConfig struct {
	Type     string            `json:"type"`
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	User     string            `json:"user"`
	Password string            `json:"password"`
	Database string            `json:"database"`
	SSL      bool              `json:"ssl"`
	Params   map[string]string `json:"params"`
}

type Report struct {
	ID             int64
	ReportedUserID string
	ReporterID     string
	Reason         string
	Timestamp      time.Time
	Confirmed      bool
	ConfirmedBy    string
}

type Ban struct {
	ID        int64
	UserID    string
	Reason    string
	AdminID   string
	Timestamp time.Time
	ExpiresAt *time.Time
}

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

var CurrentProvider DatabaseProvider

func Initialize(config DatabaseConfig) error {
	provider, ok := AvailableProviders[config.Type]
	if !ok {
		return fmt.Errorf("неизвестный тип базы данных: %s", config.Type)
	}

	err := provider.Initialize(config)
	if err != nil {
		return fmt.Errorf("ошибка инициализации базы данных: %w", err)
	}

	CurrentProvider = provider

	return nil
}

func Close() error {
	if CurrentProvider == nil {
		return nil
	}

	return CurrentProvider.Close()
}
