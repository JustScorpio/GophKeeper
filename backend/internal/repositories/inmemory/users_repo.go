package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
)

// InMemoryUsersRepo - репозиторий пользователей в памяти
type InMemoryUsersRepo struct {
	storage map[string]entities.User
	idSeq   int64
}

// NewInMemoryUsersRepo - инициализация репозитория пользователей
func NewInMemoryUsersRepo() *InMemoryUsersRepo {
	return &InMemoryUsersRepo{
		storage: make(map[string]entities.User),
	}
}

// generateID - генерация уникального ID
func (r *InMemoryUsersRepo) generateID() string {
	r.idSeq++
	return fmt.Sprintf("%d", r.idSeq)
}

// GetAll - получить все сущности
func (r *InMemoryUsersRepo) GetAll(ctx context.Context) ([]entities.User, error) {
	users := make([]entities.User, 0, len(r.storage))
	for _, user := range r.storage {
		users = append(users, user)
	}

	return users, nil
}

// Get - получить сущность по логину
func (r *InMemoryUsersRepo) Get(ctx context.Context, login string) (*entities.User, error) {
	user, exists := r.storage[login]
	if !exists {
		return nil, nil
	}

	return &user, nil
}

// Create - создать сущность
func (r *InMemoryUsersRepo) Create(ctx context.Context, dto *dtos.NewUser) (*entities.User, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	// Проверка существования пользователя
	if _, exists := r.storage[dto.Login]; exists {
		return nil, fmt.Errorf("user with login %s already exists", dto.Login)
	}

	user := entities.User{
		Login:    dto.Login,
		Password: dto.Password,
	}

	r.storage[dto.Login] = user
	return &user, nil
}

// Update - изменить сущность
func (r *InMemoryUsersRepo) Update(ctx context.Context, entity *entities.User) (*entities.User, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	// Проверка существования пользователя
	if _, exists := r.storage[entity.Login]; !exists {
		return nil, fmt.Errorf("user with login %s not found", entity.Login)
	}

	r.storage[entity.Login] = *entity
	return entity, nil
}

// Delete - удалить сущность
func (r *InMemoryUsersRepo) Delete(ctx context.Context, login string) error {
	if _, exists := r.storage[login]; !exists {
		return fmt.Errorf("user with login %s not found", login)
	}

	delete(r.storage, login)
	return nil
}
