// Репозиторий postgres
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/jackc/pgx/v5"
)

// PgTextsRepo - репозиторий с текстовыми данными
type PgTextsRepo struct {
	db *pgx.Conn
}

// NewPgTextsRepo - инициализация репозитория
func NewPgTextsRepo(db *pgx.Conn) (*PgTextsRepo, error) {
	return &PgTextsRepo{db: db}, nil
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *PgTextsRepo) GetAll(ctx context.Context) ([]entities.TextData, error) {
	userID := customcontext.GetUserID((ctx))

	rows, err := r.db.Query(ctx, "SELECT id, data, metadata, ownerid FROM Texts WHERE ownerid = $1", userID)
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
		err := rows.Scan(&text.ID, &text.Data, &text.Metadata, &text.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan text: %w", err)
		}
		texts = append(texts, text)
	}

	return texts, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *PgTextsRepo) Get(ctx context.Context, id string) (*entities.TextData, error) {
	userID := customcontext.GetUserID((ctx))

	var text entities.TextData
	err := r.db.QueryRow(ctx, "SELECT id, data, metadata, ownerid FROM Texts WHERE id = $1 AND ownerid = $2", id, userID).Scan(&text.ID, &text.Data, &text.Metadata, &text.OwnerID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, fmt.Errorf("failed to get text: %w", err)
	}

	return &text, nil
}

// Create - создать сущность
func (r *PgTextsRepo) Create(ctx context.Context, text *dtos.NewTextData) (*entities.TextData, error) {
	userID := customcontext.GetUserID(ctx)

	var entity entities.TextData
	err := r.db.QueryRow(ctx, "INSERT INTO Texts (data, metadata, ownerid) VALUES ($1, $2, $3) RETURNING id, data, metadata, ownerid", text.Data, text.Metadata, userID).Scan(&entity.ID, &entity.Data, &entity.Metadata, &entity.OwnerID)

	if err != nil {
		return nil, fmt.Errorf("failed to create text: %w", err)
	}
	return &entity, nil
}

// Update - изменить сущность
func (r *PgTextsRepo) Update(ctx context.Context, text *entities.TextData) (*entities.TextData, error) {
	userID := customcontext.GetUserID(ctx)

	var updatedEntity entities.TextData
	err := r.db.QueryRow(ctx, "UPDATE Texts SET data = $2, metadata = $3 WHERE id = $1 AND ownerid = $4 RETURNING id, data, metadata, ownerid", text.ID, text.Data, text.Metadata, userID).Scan(&updatedEntity.ID, &updatedEntity.Data, &updatedEntity.Metadata, &updatedEntity.OwnerID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, err
	}

	return &updatedEntity, nil
}

// Delete - удалить сущность
func (r *PgTextsRepo) Delete(ctx context.Context, id string) (*entities.TextData, error) {
	userID := customcontext.GetUserID((ctx))

	var deletedText entities.TextData
	err := r.db.QueryRow(ctx, "DELETE FROM Texts WHERE id = $1 AND ownerid = $2 RETURNING id, data, metadata, ownerid", id, userID).Scan(&deletedText.ID, &deletedText.Data, &deletedText.Metadata, &deletedText.OwnerID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, err
	}

	return &deletedText, err
}
