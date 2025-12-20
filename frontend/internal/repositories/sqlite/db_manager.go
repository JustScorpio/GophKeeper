// sqlite - SQLite Репозиторий
package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JustScorpio/GophKeeper/frontend/internal/repositories/sqlite/migrations"
	_ "modernc.org/sqlite"
)

type DatabaseManager struct {
	DB              *sql.DB
	BinariesRepo    *BinariesRepo
	CardsRepo       *CardsRepo
	CredentialsRepo *CredentialsRepo
	TextsRepo       *TextsRepo
}

func NewDatabaseManager(dbPath string) (*DatabaseManager, error) {
	// Создаем директорию для БД, если ее нет
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Открываем (или создаем) базу данных
	db, err := sql.Open("sqlite", "file:"+dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Включаем foreign keys и другие настройки SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;"); err != nil {
		return nil, fmt.Errorf("failed to set pragmas: %w", err)
	}

	// Запускаем миграции
	migrator := migrations.NewMigrator(db)
	if err := migrator.Migrate(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Инициализируем репозитории
	binariesRepo, err := NewBinariesRepo(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create binaries repo: %w", err)
	}

	cardsRepo, err := NewCardsRepo(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create cards repo: %w", err)
	}

	credentialsRepo, err := NewCredentialsRepo(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials repo: %w", err)
	}

	textsRepo, err := NewTextsRepo(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create texts repo: %w", err)
	}

	return &DatabaseManager{
		DB:              db,
		BinariesRepo:    binariesRepo,
		CardsRepo:       cardsRepo,
		CredentialsRepo: credentialsRepo,
		TextsRepo:       textsRepo,
	}, nil
}

func (dm *DatabaseManager) Close() error {
	return dm.DB.Close()
}
