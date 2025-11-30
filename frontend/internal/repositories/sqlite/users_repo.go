// SQLite Репозиторий для пользователей (локальный кэш)
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// UsersRepo - репозиторий пользователями
type UsersRepo struct {
	db *sql.DB
}

// NewUsersRepo - инициализация репозитория
func NewUsersRepo(db *sql.DB) (*UsersRepo, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			login TEXT PRIMARY KEY,
			password TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return &UsersRepo{db: db}, nil
}

// GetAll - получить все сущности
func (r *UsersRepo) GetAll(ctx context.Context) ([]entities.User, error) {
	rows, err := r.db.Query("SELECT login, password FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
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
func (r *UsersRepo) Get(ctx context.Context, login string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow("SELECT login, password FROM users WHERE login = ?", login).Scan(&user.Login, &user.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// Create - создать сущность
func (r *UsersRepo) Create(ctx context.Context, newUser *dtos.NewUser) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow("INSERT INTO users (login, password) VALUES (?, ?)", newUser.Login, newUser.Password).Scan(&user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// Update - изменить сущность
func (r *UsersRepo) Update(ctx context.Context, user *entities.User) (*entities.User, error) {
	var updatedEntity entities.User
	err := r.db.QueryRow("UPDATE users SET password = ?, last_sync = CURRENT_TIMESTAMP WHERE login = ?", user.Password, user.Login).Scan(&updatedEntity.Login, &updatedEntity.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &updatedEntity, nil
}

// Delete - удалить сущность
func (r *UsersRepo) Delete(ctx context.Context, login string) error {
	_, err := r.db.Exec("DELETE FROM users WHERE login = ?", login)
	return err
}
