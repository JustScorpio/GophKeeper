package inmemory

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
)

// InMemoryCredentialsRepo - репозиторий учетных данных в памяти
type InMemoryCredentialsRepo struct {
	storage map[string]entities.Credentials
	idSeq   int64
}

// NewInMemoryCredentialsRepo - инициализация репозитория учетных данных
func NewInMemoryCredentialsRepo() *InMemoryCredentialsRepo {
	return &InMemoryCredentialsRepo{
		storage: make(map[string]entities.Credentials),
	}
}

// generateID - генерация уникального ID
func (r *InMemoryCredentialsRepo) generateID() string {
	r.idSeq++
	return fmt.Sprintf("%d", r.idSeq)
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *InMemoryCredentialsRepo) GetAll(ctx context.Context) ([]entities.Credentials, error) {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	var creds []entities.Credentials
	for _, cred := range r.storage {
		if cred.OwnerID == userID {
			creds = append(creds, cred)
		}
	}

	return creds, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *InMemoryCredentialsRepo) Get(ctx context.Context, id string) (*entities.Credentials, error) {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	cred, exists := r.storage[id]
	if !exists {
		return nil, nil
	}

	if cred.OwnerID != userID {
		return nil, nil
	}

	return &cred, nil
}

// Create - создать сущность
func (r *InMemoryCredentialsRepo) Create(ctx context.Context, dto *dtos.NewCredentials) (*entities.Credentials, error) {
	if dto == nil {
		return nil, errors.New("dto cannot be nil")
	}

	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	id := r.generateID()
	cred := entities.Credentials{
		Login:        dto.Login,
		Password:     dto.Password,
		SecureEntity: entities.SecureEntity{ID: id, Metadata: dto.Metadata, OwnerID: userID},
	}

	r.storage[id] = cred
	return &cred, nil
}

// Update - изменить сущность
func (r *InMemoryCredentialsRepo) Update(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
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
		return nil, fmt.Errorf("credentials with ID %s not found", entity.ID)
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
func (r *InMemoryCredentialsRepo) Delete(ctx context.Context, id string) error {
	userID := customcontext.GetUserID(ctx)
	if userID == "" {
		return errors.New("user ID is required")
	}

	cred, exists := r.storage[id]
	if !exists {
		return fmt.Errorf("credentials with ID %s not found", id)
	}

	if cred.OwnerID != userID {
		return errors.New("access denied")
	}

	delete(r.storage, id)
	return nil
}
