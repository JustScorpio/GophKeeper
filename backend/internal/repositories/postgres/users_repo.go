package postgres

import (
	"context"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/models"
	"github.com/jackc/pgx/v5"
)

type PgUsersRepo struct {
	db *pgx.Conn
}

func NewPgUsersRepo(db *pgx.Conn) (*PgUsersRepo, error) {
	// Создание таблицы users, если её нет
	_, err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY
			login TEXT NOT NULL
			password TEXT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &PgUsersRepo{db: db}, nil
}

func (r *PgUsersRepo) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.Query(ctx, "SELECT id, login, password FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Login, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *PgUsersRepo) Get(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx, "SELECT id, login, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Login, &user.Password)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PgUsersRepo) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(ctx, "INSERT INTO users (login, password,) VALUES ($1, $2)", &user.Login, &user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *PgUsersRepo) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(ctx, "UPDATE users SET login = $2, password = $3 WHERE id = $1", user.ID, user.Login, user.Password)
	return err
}

func (r *PgUsersRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
