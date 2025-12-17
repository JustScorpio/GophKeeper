package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Migrator управляет миграциями базы данных
type Migrator struct {
	db *sql.DB
}

// NewMigrator создает новый мигратор
func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

// Migrate - применяет все миграции
func (m *Migrator) Migrate(ctx context.Context) error {
	// Создаем таблицу для отслеживания миграций
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := m.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Получаем список примененных миграций
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Читаем все файлы в директории с миграциями
	migrationsDir := "../../internal/repositories/sqlite/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	migrations := make(map[int]os.DirEntry)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			version, err := extractVersion(file.Name())
			if err != nil {
				return err
			}

			migrations[version] = file
		}
	}

	// Сортируем версии по порядку
	versions := make([]int, 0, len(migrations))
	for k := range migrations {
		versions = append(versions, k)
	}

	sort.Ints(versions)

	// Применяем миграции по порядку
	for _, version := range versions {
		if !applied[version] {
			migration := migrations[version]
			if err := m.applyMigration(ctx, filepath.Join(migrationsDir, migration.Name()), version); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Name(), err)
			}

			log.Printf("Migration applied: %s", migration.Name())
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// getAppliedMigrations - возвращает map примененных миграций
func (m *Migrator) getAppliedMigrations(ctx context.Context) (map[int]bool, error) {
	applied := make(map[int]bool)

	rows, err := m.db.QueryContext(ctx, "SELECT version FROM migrations")
	if err != nil {
		// Если таблицы нет, возвращаем пустой map
		return applied, nil
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// applyMigration - применяет миграцию
func (m *Migrator) applyMigration(ctx context.Context, filePath string, version int) error {
	// Читаем SQL файл
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Начинаем транзакцию
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Выполняем SQL миграции
	if _, err := tx.ExecContext(ctx, string(sqlBytes)); err != nil {
		tx.Rollback() // Откатываем при ошибке
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Записываем факт применения миграции
	_, err = tx.ExecContext(ctx, "INSERT INTO migrations (version) VALUES (?)", version)
	if err != nil {
		tx.Rollback() // Откатываем при ошибке
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	return nil
}

// extractVersion - извлекает номер миграции из имени файла
func extractVersion(filename string) (int, error) {
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid filename format: %s", filename)
	}
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid version in filename %s: %w", filename, err)
	}
	return version, nil
}
