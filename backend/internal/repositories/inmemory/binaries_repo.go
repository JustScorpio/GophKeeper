package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
)

// InMemoryBinariesRepo - репозиторий бинарных данных в памяти
type InMemoryBinariesRepo struct {
	storage map[string]entities.BinaryData
	idSeq   int64
}

// NewInMemoryBinariesRepo - инициализация репозитория бинарных данных
func NewInMemoryBinariesRepo() *InMemoryBinariesRepo {
	return &InMemoryBinariesRepo{
		storage: make(map[string]entities.BinaryData),
	}
}

// generateID - генерация уникального ID
func (r *InMemoryBinariesRepo) generateID() string {
	r.idSeq++
	return fmt.Sprintf("%d", r.idSeq)
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *InMemoryBinariesRepo) GetAll(ctx context.Context) ([]entities.BinaryData, error) {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	var binaries []entities.BinaryData
	for _, binary := range r.storage {
		if binary.OwnerID == userID {
			binaries = append(binaries, binary)
		}
	}

	return binaries, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *InMemoryBinariesRepo) Get(ctx context.Context, id string) (*entities.BinaryData, error) {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	binary, exists := r.storage[id]
	if !exists {
		return nil, nil
	}

	if binary.OwnerID != userID {
		return nil, nil
	}

	return &binary, nil
}

// Create - создать сущность
func (r *InMemoryBinariesRepo) Create(ctx context.Context, dto *dtos.NewBinaryData) (*entities.BinaryData, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	id := r.generateID()
	binary := entities.BinaryData{
		Data:         dto.Data,
		SecureEntity: entities.SecureEntity{ID: id, Metadata: dto.Metadata, OwnerID: userID},
	}

	r.storage[id] = binary
	return &binary, nil
}

// Update - изменить сущность
func (r *InMemoryBinariesRepo) Update(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
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
		return nil, fmt.Errorf("binary with ID %s not found", entity.ID)
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
func (r *InMemoryBinariesRepo) Delete(ctx context.Context, id string) error {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}

	binary, exists := r.storage[id]
	if !exists {
		return fmt.Errorf("binary with ID %s not found", id)
	}

	if binary.OwnerID != userID {
		return errors.New("access denied")
	}

	delete(r.storage, id)
	return nil
}
