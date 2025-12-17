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

// PgCredentialsRepo - репозиторий с учётными данными
type PgCredentialsRepo struct {
	db *pgx.Conn
}

// NewPgCredentialsRepo - инициализация репозитория
func NewPgCredentialsRepo(db *pgx.Conn) (*PgCredentialsRepo, error) {
	return &PgCredentialsRepo{db: db}, nil
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *PgCredentialsRepo) GetAll(ctx context.Context) ([]entities.Credentials, error) {
	userID := customcontext.GetUserID((ctx))

	rows, err := r.db.Query(ctx, "SELECT id, login, password, metadata, ownerid FROM Credentials WHERE ownerid = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var credentials []entities.Credentials
	for rows.Next() {
		var cred entities.Credentials
		err := rows.Scan(&cred.ID, &cred.Login, &cred.Password, &cred.Metadata, &cred.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credentials: %w", err)
		}
		credentials = append(credentials, cred)
	}

	return credentials, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *PgCredentialsRepo) Get(ctx context.Context, id string) (*entities.Credentials, error) {
	userID := customcontext.GetUserID((ctx))

	var credentials entities.Credentials
	err := r.db.QueryRow(ctx, "SELECT id, login, password, metadata, ownerid FROM Credentials WHERE id = $1 AND ownerID = $2", id, userID).Scan(&credentials.ID, &credentials.Login, &credentials.Password, &credentials.Metadata, &credentials.OwnerID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	return &credentials, nil
}

// Create - создать сущность
func (r *PgCredentialsRepo) Create(ctx context.Context, credentials *dtos.NewCredentials) (*entities.Credentials, error) {
	userID := customcontext.GetUserID(ctx)

	var entity entities.Credentials
	err := r.db.QueryRow(ctx, "INSERT INTO Credentials (login, password, metadata, ownerid) VALUES ($1, $2, $3, $4) RETURNING id, login, password, metadata, ownerid", credentials.Login, credentials.Password, credentials.Metadata, userID).Scan(&entity.ID, &entity.Login, &entity.Password, &entity.Metadata, &entity.OwnerID)

	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}
	return &entity, nil
}

// Update - изменить сущность
func (r *PgCredentialsRepo) Update(ctx context.Context, credentials *entities.Credentials) (*entities.Credentials, error) {
	userID := customcontext.GetUserID(ctx)

	var updatedEntity entities.Credentials
	err := r.db.QueryRow(ctx, "UPDATE Credentials SET login = $2, password = $3, metadata = $4 WHERE id = $1 AND ownerid = $5 RETURNING id, login, password, metadata, ownerid", credentials.ID, credentials.Login, credentials.Password, credentials.Metadata, userID).Scan(&updatedEntity.ID, &updatedEntity.Login, &updatedEntity.Password, &updatedEntity.Metadata, &updatedEntity.OwnerID)

	if err != nil {
		return nil, fmt.Errorf("failed to update credentials: %w", err)
	}
	return &updatedEntity, nil
}

// Delete - удалить сущность
func (r *PgCredentialsRepo) Delete(ctx context.Context, id string) error {
	userID := customcontext.GetUserID((ctx))

	_, err := r.db.Exec(ctx, "DELETE FROM Credentials WHERE id = $1 AND ownerid = $2", id, userID)
	return err
}
