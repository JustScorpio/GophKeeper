package postgres

import (
	"context"
	"fmt"
	"strings"

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

func NewDBConnection(connStr string) (*pgx.Conn, error) {
	conf := extractConnectionConfig(connStr)

	// Создание базы данных
	defaultDBConnStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", conf.Host, conf.User, conf.Password, "postgres", conf.Port, conf.SslMode)
	defaultDB, err := pgx.Connect(context.Background(), defaultDBConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to default database: %w", err)
	}
	defer defaultDB.Close(context.Background())

	// Проверка и создание базы данных
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
