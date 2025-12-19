// Репозиторий postgres
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/jackc/pgx/v5"
)

// PgUsersRepo - репозиторий пользователями
type PgUsersRepo struct {
	db *pgx.Conn
}

// NewPgUsersRepo - инициализация репозитория
func NewPgUsersRepo(db *pgx.Conn) (*PgUsersRepo, error) {
	return &PgUsersRepo{db: db}, nil
}

// GetAll - получить все сущности
func (r *PgUsersRepo) GetAll(ctx context.Context) ([]entities.User, error) {
	rows, err := r.db.Query(ctx, "SELECT login, password FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var users []entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.Login, &user.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Get - получить сущность по ИД
func (r *PgUsersRepo) Get(ctx context.Context, login string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow(ctx, "SELECT login, password FROM users WHERE login = $1", login).Scan(&user.Login, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, fmt.Errorf("failed to get entity: %w", err)
	}

	return &user, nil
}

// Create - создать сущность
func (r *PgUsersRepo) Create(ctx context.Context, user *dtos.NewUser) (*entities.User, error) {
	var entity entities.User
	err := r.db.QueryRow(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING login, password", user.Login, user.Password).Scan(&entity.Login, &entity.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &entity, nil
}

// Update - изменить сущность
func (r *PgUsersRepo) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	var updatedEntity entities.User
	err := r.db.QueryRow(ctx, "UPDATE users SET password = $2 WHERE login = $1 RETURNING login, password", user.Login, user.Password).Scan(&updatedEntity.Login, &updatedEntity.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, err
	}

	return &updatedEntity, nil
}

// Delete - удалить сущность
func (r *PgUsersRepo) Delete(ctx context.Context, login string) (*entities.User, error) {
	var deletedUser entities.User
	err := r.db.QueryRow(ctx, "DELETE FROM users WHERE login = $1 RETURNING login, password", login).Scan(&deletedUser.Login, &deletedUser.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, err
	}

	return &deletedUser, err
}
