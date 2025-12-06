package inmemory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// InMemoryCardsRepo - репозиторий с данными банковских карт в памяти
type InMemoryCardsRepo struct {
	storage map[string]entities.CardInformation
}

// NewInMemoryCardsRepo - инициализация репозитория банковских карт
func NewInMemoryCardsRepo() *InMemoryCardsRepo {
	return &InMemoryCardsRepo{
		storage: make(map[string]entities.CardInformation),
	}
}

// GetAll - получить все сущности
func (r *InMemoryCardsRepo) GetAll(ctx context.Context) ([]entities.CardInformation, error) {
	cards := make([]entities.CardInformation, 0, len(r.storage))
	for _, card := range r.storage {
		cards = append(cards, card)
	}

	return cards, nil
}

// Get - получить сущность по ИД
func (r *InMemoryCardsRepo) Get(ctx context.Context, id string) (*entities.CardInformation, error) {
	card, exists := r.storage[id]
	if !exists {
		return nil, nil
	}

	return &card, nil
}

// Create - создать сущность
func (r *InMemoryCardsRepo) Create(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	if entity.ID == "" {
		return nil, errors.New("ID cannot be empty")
	}

	// Проверяем, не существует ли уже запись с таким ID
	if _, exists := r.storage[entity.ID]; exists {
		return nil, fmt.Errorf("card with ID %s already exists", entity.ID)
	}

	r.storage[entity.ID] = *entity
	return entity, nil
}

// Update - изменить сущность
func (r *InMemoryCardsRepo) Update(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	// Валидация даты
	if _, err := time.Parse("01/06", entity.ExpirationDate); err != nil {
		return nil, fmt.Errorf("invalid expiration date format, expected MM/YY: %w", err)
	}

	// Проверяем существование
	if _, exists := r.storage[entity.ID]; !exists {
		return nil, fmt.Errorf("card with ID %s not found", entity.ID)
	}

	r.storage[entity.ID] = *entity
	return entity, nil
}

// Delete - удалить сущность
func (r *InMemoryCardsRepo) Delete(ctx context.Context, id string) error {
	if _, exists := r.storage[id]; !exists {
		return fmt.Errorf("card with ID %s not found", id)
	}

	delete(r.storage, id)
	return nil
}
