package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/JustScorpio/GophKeeper/backend/internal/repositories/postgres/migrations"
	"github.com/jackc/pgx/v5"
)

type connectionConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SslMode  string
}

type DatabaseManager struct {
	DB              *pgx.Conn
	BinariesRepo    *PgBinariesRepo
	CardsRepo       *PgCardsRepo
	CredentialsRepo *PgCredentialsRepo
	TextsRepo       *PgTextsRepo
	UsersRepo       *PgUsersRepo
}

func InitDatabase(connStr string) (*pgx.Conn, error) {
	conf := extractConnectionConfig(connStr)

	// Подключение к базе данных postgres по умолчанию
	defaultDBConnStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", conf.Host, conf.User, conf.Password, "postgres", conf.Port, conf.SslMode)
	defaultDB, err := pgx.Connect(context.Background(), defaultDBConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to default database: %w", err)
	}

	defer defaultDB.Close(context.Background())

	// Проверка наличия базы данных
	var dbExists bool
	err = defaultDB.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)", conf.DBName).Scan(&dbExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check database existence: %w", err)
	}

	// Создание базы данных, если она не существует
	if !dbExists {
		_, err = defaultDB.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", conf.DBName))
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	// Подключение к базе данных
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	//Проверка подключения
	if err = db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewDatabaseManager(connStr string) (*DatabaseManager, error) {
	db, err := InitDatabase(connStr)
	if err != nil {
		return nil, err
	}

	// Запускаем миграции
	migrator := migrations.NewMigrator(db)
	if err := migrator.Migrate(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	binariesRepo, err := NewPgBinariesRepo(db)
	if err != nil {
		return nil, err
	}
	cardsRepo, err := NewPgCardsRepo(db)
	if err != nil {
		return nil, err
	}
	credentialsRepo, err := NewPgCredentialsRepo(db)
	if err != nil {
		return nil, err
	}
	textsRepo, err := NewPgTextsRepo(db)
	if err != nil {
		return nil, err
	}
	usersRepo, err := NewPgUsersRepo(db)
	if err != nil {
		return nil, err
	}

	dbManager := DatabaseManager{
		DB:              db,
		BinariesRepo:    binariesRepo,
		CardsRepo:       cardsRepo,
		CredentialsRepo: credentialsRepo,
		TextsRepo:       textsRepo,
		UsersRepo:       usersRepo,
	}

	return &dbManager, nil
}

// Только чтобы не возиться с кучей параметров при запуске приложения
// extractDBName - получить название базы данных из строки подключения
func extractConnectionConfig(connStr string) *connectionConfig {
	var result connectionConfig
	parts := strings.Split(connStr, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "host=") {
			result.Host = strings.TrimPrefix(part, "host=")
		}
		if strings.HasPrefix(part, "user=") {
			result.User = strings.TrimPrefix(part, "user=")
		}
		if strings.HasPrefix(part, "password=") {
			result.Password = strings.TrimPrefix(part, "password=")
		}
		if strings.HasPrefix(part, "dbname=") {
			result.DBName = strings.TrimPrefix(part, "dbname=")
		}
		if strings.HasPrefix(part, "port=") {
			result.Port = strings.TrimPrefix(part, "port=")
		}
		if strings.HasPrefix(part, "sslmode=") {
			result.SslMode = strings.TrimPrefix(part, "sslmode=")
		}
	}
	return &result
}
