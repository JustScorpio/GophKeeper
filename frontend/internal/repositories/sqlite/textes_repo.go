// sqlite - SQLite Репозиторий
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// TextsRepo - репозиторий с текстовыми данными
type TextsRepo struct {
	db *sql.DB
}

// NewTextsRepo - инициализация репозитория
func NewTextsRepo(db *sql.DB) (*TextsRepo, error) {
	return &TextsRepo{db: db}, nil
}

// GetAll - получить все сущности
func (r *TextsRepo) GetAll(ctx context.Context) ([]entities.TextData, error) {
	rows, err := r.db.Query("SELECT id, data, metadata FROM texts")
	if err != nil {
		return nil, fmt.Errorf("failed to get texts: %w", err)
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var texts []entities.TextData
	for rows.Next() {
		var text entities.TextData
		err := rows.Scan(&text.ID, &text.Data, &text.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to scan text: %w", err)
		}
		texts = append(texts, text)
	}

	return texts, nil
}

// Get - получить сущность по ИД
func (r *TextsRepo) Get(ctx context.Context, id string) (*entities.TextData, error) {
	var text entities.TextData
	err := r.db.QueryRow("SELECT id, data, metadata FROM texts WHERE id = ?", id).Scan(&text.ID, &text.Data, &text.Metadata)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get text: %w", err)
	}

	return &text, nil
}

// Create - создать сущность
func (r *TextsRepo) Create(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	var text entities.TextData
	err := r.db.QueryRow("INSERT INTO texts (id, data, metadata) VALUES (?, ?, ?) RETURNING id, data, metadata", entity.ID, entity.Data, entity.Metadata).Scan(&text.ID, &text.Data, &text.Metadata)

	if err != nil {
		return nil, fmt.Errorf("failed to create text: %w", err)
	}

	return &text, nil
}

// Update - изменить сущность
func (r *TextsRepo) Update(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	var updatedText entities.TextData
	err := r.db.QueryRow("UPDATE texts SET data = ?, metadata = ? WHERE id = ? RETURNING id, data, metadata", entity.Data, entity.Metadata, entity.ID).Scan(&updatedText.ID, &updatedText.Data, &updatedText.Metadata)

	if err != nil {
		return nil, fmt.Errorf("failed to update text: %w", err)
	}

	return &updatedText, nil
}

func (r *TextsRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM texts WHERE id = ?", id)
	return err
}
