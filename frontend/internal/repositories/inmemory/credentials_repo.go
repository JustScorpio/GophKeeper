package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// InMemoryCredentialsRepo - репозиторий с учётными данными в памяти
type InMemoryCredentialsRepo struct {
	storage map[string]entities.Credentials
}

// NewInMemoryCredentialsRepo - инициализация репозитория учетных данных
func NewInMemoryCredentialsRepo() *InMemoryCredentialsRepo {
	return &InMemoryCredentialsRepo{
		storage: make(map[string]entities.Credentials),
	}
}

// GetAll - получить все сущности
func (r *InMemoryCredentialsRepo) GetAll(ctx context.Context) ([]entities.Credentials, error) {
	creds := make([]entities.Credentials, 0, len(r.storage))
	for _, cred := range r.storage {
		creds = append(creds, cred)
	}

	return creds, nil
}

// Get - получить сущность по ИД
func (r *InMemoryCredentialsRepo) Get(ctx context.Context, id string) (*entities.Credentials, error) {
	cred, exists := r.storage[id]
	if !exists {
		return nil, nil
	}

	return &cred, nil
}

// Create - создать сущность
func (r *InMemoryCredentialsRepo) Create(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	if entity.ID == "" {
		return nil, errors.New("ID cannot be empty")
	}

	// Проверяем, не существует ли уже запись с таким ID
	if _, exists := r.storage[entity.ID]; exists {
		return nil, fmt.Errorf("credentials with ID %s already exists", entity.ID)
	}

	r.storage[entity.ID] = *entity
	return entity, nil
}

// Update - изменить сущность
func (r *InMemoryCredentialsRepo) Update(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	// Проверяем существование
	if _, exists := r.storage[entity.ID]; !exists {
		return nil, fmt.Errorf("credentials with ID %s not found", entity.ID)
	}

	r.storage[entity.ID] = *entity
	return entity, nil
}

// Delete - удалить сущность
func (r *InMemoryCredentialsRepo) Delete(ctx context.Context, id string) error {
	if _, exists := r.storage[id]; !exists {
		return fmt.Errorf("credentials with ID %s not found", id)
	}

	delete(r.storage, id)
	return nil
}
