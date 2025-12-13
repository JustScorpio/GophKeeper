package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// InMemoryTextsRepo - репозиторий с текстовыми данными в памяти
type InMemoryTextsRepo struct {
	storage map[string]entities.TextData
}

// NewInMemoryTextsRepo - инициализация репозитория текстовых данных
func NewInMemoryTextsRepo() *InMemoryTextsRepo {
	return &InMemoryTextsRepo{
		storage: make(map[string]entities.TextData),
	}
}

// GetAll - получить все сущности
func (r *InMemoryTextsRepo) GetAll(ctx context.Context) ([]entities.TextData, error) {
	texts := make([]entities.TextData, 0, len(r.storage))
	for _, text := range r.storage {
		texts = append(texts, text)
	}

	return texts, nil
}

// Get - получить сущность по ИД
func (r *InMemoryTextsRepo) Get(ctx context.Context, id string) (*entities.TextData, error) {
	text, exists := r.storage[id]
	if !exists {
		return nil, nil
	}

	return &text, nil
}

// Create - создать сущность
func (r *InMemoryTextsRepo) Create(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	if entity.ID == "" {
		return nil, errors.New("ID cannot be empty")
	}

	// Проверяем, не существует ли уже запись с таким ID
	if _, exists := r.storage[entity.ID]; exists {
		return nil, fmt.Errorf("text with ID %s already exists", entity.ID)
	}

	r.storage[entity.ID] = *entity
	return entity, nil
}

// Update - изменить сущность
func (r *InMemoryTextsRepo) Update(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	// Проверяем существование
	if _, exists := r.storage[entity.ID]; !exists {
		return nil, fmt.Errorf("text with ID %s not found", entity.ID)
	}

	r.storage[entity.ID] = *entity
	return entity, nil
}

// Delete - удалить сущность
func (r *InMemoryTextsRepo) Delete(ctx context.Context, id string) error {
	if _, exists := r.storage[id]; !exists {
		return fmt.Errorf("text with ID %s not found", id)
	}

	delete(r.storage, id)
	return nil
}
