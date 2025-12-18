// inmemory содержит репозиторий который хранит данные в оперативной памяти
package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
)

// InMemoryCardsRepo - репозиторий банковских карт в памяти
type InMemoryCardsRepo struct {
	storage map[string]entities.CardInformation
	idSeq   int64
}

// NewInMemoryCardsRepo - инициализация репозитория банковских карт
func NewInMemoryCardsRepo() *InMemoryCardsRepo {
	return &InMemoryCardsRepo{
		storage: make(map[string]entities.CardInformation),
	}
}

// generateID - генерация уникального ID
func (r *InMemoryCardsRepo) generateID() string {
	r.idSeq++
	return fmt.Sprintf("%d", r.idSeq)
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *InMemoryCardsRepo) GetAll(ctx context.Context) ([]entities.CardInformation, error) {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	var cards []entities.CardInformation
	for _, card := range r.storage {
		if card.OwnerID == userID {
			cards = append(cards, card)
		}
	}

	return cards, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *InMemoryCardsRepo) Get(ctx context.Context, id string) (*entities.CardInformation, error) {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	card, exists := r.storage[id]
	if !exists {
		return nil, nil
	}

	if card.OwnerID != userID {
		return nil, nil
	}

	return &card, nil
}

// Create - создать сущность
func (r *InMemoryCardsRepo) Create(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	id := r.generateID()
	card := entities.CardInformation{
		Number:         dto.Number,
		CardHolder:     dto.CardHolder,
		ExpirationDate: dto.ExpirationDate,
		CVV:            dto.CVV,
		SecureEntity:   entities.SecureEntity{ID: id, Metadata: dto.Metadata, OwnerID: userID},
	}

	r.storage[id] = card
	return &card, nil
}

// Update - изменить сущность
func (r *InMemoryCardsRepo) Update(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	if entity == nil {
		return nil, errors.New("entity cannot be nil")
	}

	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Проверяем существование и права
	existing, exists := r.storage[entity.ID]
	if !exists {
		return nil, fmt.Errorf("card with ID %s not found", entity.ID)
	}

	if existing.OwnerID != userID {
		return nil, errors.New("access denied")
	}

	// Обновляем только разрешенные поля
	entity.OwnerID = userID
	r.storage[entity.ID] = *entity

	return entity, nil
}

// Delete - удалить сущность
func (r *InMemoryCardsRepo) Delete(ctx context.Context, id string) error {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}

	card, exists := r.storage[id]
	if !exists {
		return fmt.Errorf("card with ID %s not found", id)
	}

	if card.OwnerID != userID {
		return errors.New("access denied")
	}

	delete(r.storage, id)
	return nil
}
