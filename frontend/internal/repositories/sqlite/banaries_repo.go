// SQLite Репозиторий
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// BinariesRepo - репозиторий с бинарными данными
type BinariesRepo struct {
	db *sql.DB
}

// NewBinariesRepo - инициализация репозитория
func NewBinariesRepo(db *sql.DB) (*BinariesRepo, error) {
	// Создание таблицы Binaries, если её нет
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS binaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data BLOB NOT NULL,
			metadata TEXT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &BinariesRepo{db: db}, nil
}

// GetAll - получить все сущности
func (r *BinariesRepo) GetAll(ctx context.Context) ([]entities.BinaryData, error) {
	rows, err := r.db.Query("SELECT id, data, metadata FROM binaries")
	if err != nil {
		return nil, fmt.Errorf("failed to get entities: %w", err)
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var binaries []entities.BinaryData
	for rows.Next() {
		var binary entities.BinaryData
		err := rows.Scan(&binary.ID, &binary.Data, &binary.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entity: %w", err)
		}
		binaries = append(binaries, binary)
	}

	return binaries, nil
}

// Get - получить сущность по ИД
func (r *BinariesRepo) Get(ctx context.Context, id string) (*entities.BinaryData, error) {
	var binaryData entities.BinaryData
	err := r.db.QueryRow("SELECT id, data, metadata FROM binaries WHERE id = ?", id).Scan(&binaryData.ID, &binaryData.Data, &binaryData.Metadata)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get entity: %w", err)
	}

	return &binaryData, nil
}

// Create - создать сущность
func (r *BinariesRepo) Create(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	var binary entities.BinaryData
	err := r.db.QueryRow("INSERT INTO binaries (id, data, metadata) VALUES (?, ?, ?) RETURNING id, data, metadata", entity.ID, entity.Data, entity.Metadata).Scan(&binary.ID, &binary.Data, &binary.Metadata)

	if err != nil {
		return nil, fmt.Errorf("failed to create entity: %w", err)
	}

	return &binary, nil
}

// Update - изменить сущность
func (r *BinariesRepo) Update(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	var updatedBinary entities.BinaryData
	err := r.db.QueryRow("UPDATE binaries SET data = ?, metadata = ? WHERE id = ? RETURNING id, data, metadata", entity.Data, entity.Metadata, entity.ID).Scan(&updatedBinary.ID, &updatedBinary.Data, &updatedBinary.Metadata)

	if err != nil {
		return nil, fmt.Errorf("failed to update entity: %w", err)
	}

	return &updatedBinary, nil
}

// Delete - удалить сущность
func (r *BinariesRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM binaries WHERE id = ?", id)
	return err
}
