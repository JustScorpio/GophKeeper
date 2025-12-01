// SQLite Репозиторий для учётных данных
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// CredentialsRepo - репозиторий с учётными данными
type CredentialsRepo struct {
	db *sql.DB
}

// NewCredentialsRepo - инициализация репозитория
func NewCredentialsRepo(db *sql.DB) (*CredentialsRepo, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS credentials (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT NOT NULL,
			password TEXT NOT NULL,
			metadata TEXT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &CredentialsRepo{db: db}, nil
}

// GetAll - получить все сущности
func (r *CredentialsRepo) GetAll(ctx context.Context) ([]entities.Credentials, error) {
	rows, err := r.db.Query("SELECT id, login, password, metadata FROM credentials")
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
		err := rows.Scan(&cred.ID, &cred.Login, &cred.Password, &cred.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credentials: %w", err)
		}
		credentials = append(credentials, cred)
	}

	return credentials, nil
}

// Get - получить сущность по ИД
func (r *CredentialsRepo) Get(ctx context.Context, id string) (*entities.Credentials, error) {
	var cred entities.Credentials
	err := r.db.QueryRow(
		"SELECT id, login, password, metadata FROM credentials WHERE id = ?", id).Scan(&cred.ID, &cred.Login, &cred.Password, &cred.Metadata)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	return &cred, nil
}

// Create - создать сущность
func (r *CredentialsRepo) Create(ctx context.Context, dto *dtos.NewCredentials) (*entities.Credentials, error) {
	var cred entities.Credentials
	err := r.db.QueryRow("INSERT INTO credentials (login, password, metadata) VALUES (?, ?, ?) RETURNING id, login, password, metadata", dto.Login, dto.Password, dto.Metadata).Scan(&cred.ID, &cred.Login, &cred.Password, &cred.Metadata)

	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	return &cred, nil
}

// Update - изменить сущность
func (r *CredentialsRepo) Update(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	var updatedCred entities.Credentials
	err := r.db.QueryRow("UPDATE credentials SET login = ?, password = ?, metadata = ? WHERE id = ? RETURNING id, login, password, metadata", entity.Login, entity.Password, entity.Metadata, entity.ID).Scan(&updatedCred.ID, &updatedCred.Login, &updatedCred.Password, &updatedCred.Metadata)

	if err != nil {
		return nil, fmt.Errorf("failed to update credentials: %w", err)
	}

	return &updatedCred, nil
}

// Delete - удалить сущность
func (r *CredentialsRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM credentials WHERE id = ?", id)
	return err
}
