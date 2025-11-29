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
	// Создание таблицы Texts, если её нет
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS Texts (
			id SERIAL PRIMARY KEY,
			data TEXT NOT NULL,
			metadata TEXT,
			ownerid TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &PgTextsRepo{db: db}, nil
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *PgTextsRepo) GetAll(ctx context.Context) ([]entities.TextData, error) {
	userID := customcontext.GetUserID((ctx))

	rows, err := r.db.Query(ctx, "SELECT id, data, metadata, ownerid FROM Texts WHERE ownerid = $1", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var Texts []entities.TextData
	for rows.Next() {
		var textData entities.TextData
		err := rows.Scan(&textData.ID, &textData.Data, &textData.Metadata, &textData.OwnerID)
		if err != nil {
			return nil, err
		}
		Texts = append(Texts, textData)
	}

	return Texts, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *PgTextsRepo) Get(ctx context.Context, id string) (*entities.TextData, error) {
	userID := customcontext.GetUserID((ctx))

	var textData entities.TextData
	err := r.db.QueryRow(ctx, "SELECT id, data, metadata, ownerid FROM Texts WHERE id = $1 AND ownerid = $2", id, userID).Scan(&textData.ID, &textData.Data, &textData.Metadata, &textData.OwnerID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, fmt.Errorf("failed to get entity: %w", err)
	}

	return &textData, nil
}

// Create - создать сущность
func (r *PgTextsRepo) Create(ctx context.Context, textData *dtos.NewTextData) (*entities.TextData, error) {
	userID := customcontext.GetUserID(ctx)

	var entity entities.TextData
	err := r.db.QueryRow(ctx, "INSERT INTO Texts (data, metadata, ownerid) VALUES ($1, $2, $3) RETURNING id, data, metadata, ownerid", textData.Data, textData.Metadata, userID).Scan(&entity.ID, &entity.Data, &entity.Metadata, &entity.OwnerID)

	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update - изменить сущность
func (r *PgTextsRepo) Update(ctx context.Context, textData *entities.TextData) (*entities.TextData, error) {
	userID := customcontext.GetUserID(ctx)

	var updatedEntity entities.TextData
	err := r.db.QueryRow(ctx, "UPDATE Texts SET data = $2, metadata = $3, ownerid = $4 WHERE id = $1 AND ownerid = $5 RETURNING id, data, metadata, ownerid", textData.ID, textData.Data, textData.Metadata, textData.OwnerID, userID).Scan(&updatedEntity.ID, &updatedEntity.Data, &updatedEntity.Metadata, &updatedEntity.OwnerID)

	if err != nil {
		return nil, err
	}
	return &updatedEntity, nil
}

// Delete - удалить сущность
func (r *PgTextsRepo) Delete(ctx context.Context, id string) error {
	userID := customcontext.GetUserID((ctx))

	_, err := r.db.Exec(ctx, "DELETE FROM Texts WHERE id = $1 AND ownerid = $2", id, userID)
	return err
}
