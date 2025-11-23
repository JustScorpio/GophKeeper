// Репозиторий postgres
package postgres

import (
	"context"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/jackc/pgx/v5"
)

// PgBinariesRepo - репозиторий с бинарными данными
type PgBinariesRepo struct {
	db *pgx.Conn
}

// NewPgBinariesRepo - инициализация репозитория
func NewPgBinariesRepo(db *pgx.Conn) (*PgBinariesRepo, error) {
	// Создание таблицы Binaries, если её нет
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS Binaries (
			id SERIAL PRIMARY KEY,
			data BYTEA NOT NULL,
			metadata TEXT,
			ownerid TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &PgBinariesRepo{db: db}, nil
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *PgBinariesRepo) GetAll(ctx context.Context) ([]entities.BinaryData, error) {
	userID := customcontext.GetUserID((ctx))

	rows, err := r.db.Query(ctx, "SELECT id, data, metadata, ownerid FROM Binaries WHERE ownerid = $1", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var Binaries []entities.BinaryData
	for rows.Next() {
		var binaryData entities.BinaryData
		err := rows.Scan(&binaryData.ID, &binaryData.Data, &binaryData.Metadata, &binaryData.OwnerID)
		if err != nil {
			return nil, err
		}
		Binaries = append(Binaries, binaryData)
	}

	return Binaries, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *PgBinariesRepo) Get(ctx context.Context, id string) (*entities.BinaryData, error) {
	userID := customcontext.GetUserID((ctx))

	var binaryData entities.BinaryData
	err := r.db.QueryRow(ctx, "SELECT id, data, metadata, ownerid FROM Binaries WHERE id = $1 AND ownerid = $2", id, userID).Scan(&binaryData.ID, &binaryData.Data, &binaryData.Metadata, &binaryData.OwnerID)

	if err != nil {
		return nil, err
	}
	return &binaryData, nil
}

// Create - создать сущность
func (r *PgBinariesRepo) Create(ctx context.Context, binaryData *dtos.NewBinaryData) (*entities.BinaryData, error) {
	userID := customcontext.GetUserID(ctx)

	var entity entities.BinaryData
	err := r.db.QueryRow(ctx, "INSERT INTO Binaries (data, metadata, ownerid) VALUES ($1, $2, $3) RETURNING id, data, metadata, ownerid", binaryData.Data, binaryData.Metadata, userID).Scan(&entity.ID, &entity.Data, &entity.Metadata, &entity.OwnerID)

	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update - изменить сущность
func (r *PgBinariesRepo) Update(ctx context.Context, binaryData *entities.BinaryData) (*entities.BinaryData, error) {
	userID := customcontext.GetUserID(ctx)

	var updatedEntity entities.BinaryData
	err := r.db.QueryRow(ctx, "UPDATE Binaries SET data = $2, metadata = $3, ownerid = $4 WHERE id = $1 AND ownerid = $5 RETURNING id, data, metadata, ownerid", binaryData.ID, binaryData.Data, binaryData.Metadata, binaryData.OwnerID, userID).Scan(&updatedEntity.ID, &updatedEntity.Data, &updatedEntity.Metadata, &updatedEntity.OwnerID)

	if err != nil {
		return nil, err
	}
	return &updatedEntity, nil
}

// Delete - удалить сущность
func (r *PgBinariesRepo) Delete(ctx context.Context, id string) error {
	userID := customcontext.GetUserID((ctx))

	_, err := r.db.Exec(ctx, "DELETE FROM Binaries WHERE id = $1 AND ownerid = $2", id, userID)
	return err
}
