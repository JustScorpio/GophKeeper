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

// InMemoryTextsRepo - репозиторий текстовых данных в памяти
type InMemoryTextsRepo struct {
	storage map[string]entities.TextData
	idSeq   int64
}

// NewInMemoryTextsRepo - инициализация репозитория текстовых данных
func NewInMemoryTextsRepo() *InMemoryTextsRepo {
	return &InMemoryTextsRepo{
		storage: make(map[string]entities.TextData),
	}
}

// generateID - генерация уникального ID
func (r *InMemoryTextsRepo) generateID() string {
	r.idSeq++
	return fmt.Sprintf("%d", r.idSeq)
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *InMemoryTextsRepo) GetAll(ctx context.Context) ([]entities.TextData, error) {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	var texts []entities.TextData
	for _, text := range r.storage {
		if text.OwnerID == userID {
			texts = append(texts, text)
		}
	}

	return texts, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *InMemoryTextsRepo) Get(ctx context.Context, id string) (*entities.TextData, error) {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	text, exists := r.storage[id]
	if !exists {
		return nil, nil
	}

	if text.OwnerID != userID {
		return nil, nil
	}

	return &text, nil
}

// Create - создать сущность
func (r *InMemoryTextsRepo) Create(ctx context.Context, dto *dtos.NewTextData) (*entities.TextData, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	id := r.generateID()
	text := entities.TextData{
		Data:         dto.Data,
		SecureEntity: entities.SecureEntity{ID: id, Metadata: dto.Metadata, OwnerID: userID},
	}

	r.storage[id] = text
	return &text, nil
}

// Update - изменить сущность
func (r *InMemoryTextsRepo) Update(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
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
		return nil, fmt.Errorf("text with ID %s not found", entity.ID)
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
func (r *InMemoryTextsRepo) Delete(ctx context.Context, id string) error {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}

	text, exists := r.storage[id]
	if !exists {
		return fmt.Errorf("text with ID %s not found", id)
	}

	if text.OwnerID != userID {
		return errors.New("access denied")
	}

	delete(r.storage, id)
	return nil
}
