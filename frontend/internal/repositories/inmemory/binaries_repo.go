// inmemory - репозиторий хранящий данные воперативной памяти
package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// InMemoryBinariesRepo - репозиторий с бинарными данными в памяти
type InMemoryBinariesRepo struct {
	storage map[string]entities.BinaryData
}

// NewInMemoryBinariesRepo - инициализация репозитория бинарных данных
func NewInMemoryBinariesRepo() *InMemoryBinariesRepo {
	return &InMemoryBinariesRepo{
		storage: make(map[string]entities.BinaryData),
	}
}

// GetAll - получить все сущности
func (r *InMemoryBinariesRepo) GetAll(ctx context.Context) ([]entities.BinaryData, error) {
	binaries := make([]entities.BinaryData, 0, len(r.storage))
	for _, binary := range r.storage {
		binaries = append(binaries, binary)
	}

	return binaries, nil
}

// Get - получить сущность по ИД
func (r *InMemoryBinariesRepo) Get(ctx context.Context, id string) (*entities.BinaryData, error) {
	binary, exists := r.storage[id]
	if !exists {
		return nil, nil
	}

	return &binary, nil
}

// Create - создать сущность
func (r *InMemoryBinariesRepo) Create(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	if entity.ID == "" {
		return nil, errors.New("ID cannot be empty")
	}

	// Проверяем, не существует ли уже запись с таким ID
	if _, exists := r.storage[entity.ID]; exists {
		return nil, fmt.Errorf("binary with ID %s already exists", entity.ID)
	}

	r.storage[entity.ID] = *entity
	return entity, nil
}

// Update - изменить сущность
func (r *InMemoryBinariesRepo) Update(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	// Проверяем существование
	if _, exists := r.storage[entity.ID]; !exists {
		return nil, fmt.Errorf("binary with ID %s not found", entity.ID)
	}

	r.storage[entity.ID] = *entity
	return entity, nil
}

// Delete - удалить сущность
func (r *InMemoryBinariesRepo) Delete(ctx context.Context, id string) error {
	if _, exists := r.storage[id]; !exists {
		return fmt.Errorf("binary with ID %s not found", id)
	}

	delete(r.storage, id)
	return nil
}
