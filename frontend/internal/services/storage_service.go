// services - сервисы
package services

import (
	"context"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/frontend/internal/repositories"
)

// StorageService - сервис для работы с хранилищем
type StorageService struct {
	binariesRepo    repositories.IRepository[entities.BinaryData]
	cardsRepo       repositories.IRepository[entities.CardInformation]
	credentialsRepo repositories.IRepository[entities.Credentials]
	textsRepo       repositories.IRepository[entities.TextData]
}

// NewStorageService - создать сервис для работы с хранилищем
func NewStorageService(
	binariesRepo repositories.IRepository[entities.BinaryData],
	cardsRepo repositories.IRepository[entities.CardInformation],
	credentialsRepo repositories.IRepository[entities.Credentials],
	textsRepo repositories.IRepository[entities.TextData],
) *StorageService {
	return &StorageService{
		binariesRepo:    binariesRepo,
		cardsRepo:       cardsRepo,
		credentialsRepo: credentialsRepo,
		textsRepo:       textsRepo,
	}
}

// CreateBinary - создать запись с бинарными данными
func (s *StorageService) CreateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	return s.binariesRepo.Create(ctx, entity)
}

// GetBinary - получить бинарные данные
func (s *StorageService) GetBinary(ctx context.Context, id string) (*entities.BinaryData, error) {
	return s.binariesRepo.Get(ctx, id)
}

// GetAllBinaries - получить все бинарные данные
func (s *StorageService) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	return s.binariesRepo.GetAll(ctx)
}

// UpdateBinary - обновить бинарные данные
func (s *StorageService) UpdateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	return s.binariesRepo.Update(ctx, entity)
}

// DeleteBinary - удалить бинарные данные
func (s *StorageService) DeleteBinary(ctx context.Context, id string) error {
	return s.binariesRepo.Delete(ctx, id)
}

// CreateCard - создать запись с данными карты
func (s *StorageService) CreateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	return s.cardsRepo.Create(ctx, entity)
}

// GetCard - получить данные карты
func (s *StorageService) GetCard(ctx context.Context, id string) (*entities.CardInformation, error) {
	return s.cardsRepo.Get(ctx, id)
}

// GetAllCards - получить все данные карт
func (s *StorageService) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	return s.cardsRepo.GetAll(ctx)
}

// UpdateCard - обновить данные карты
func (s *StorageService) UpdateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	return s.cardsRepo.Update(ctx, entity)
}

// DeleteCard - удалить данные карты
func (s *StorageService) DeleteCard(ctx context.Context, id string) error {
	return s.cardsRepo.Delete(ctx, id)
}

// CreateCredentials - создать запись с учётными данными
func (s *StorageService) CreateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	return s.credentialsRepo.Create(ctx, entity)
}

// GetCredentials - получить учётные данные
func (s *StorageService) GetCredentials(ctx context.Context, id string) (*entities.Credentials, error) {
	return s.credentialsRepo.Get(ctx, id)
}

// GetAllCredentials - получить все учётные данные
func (s *StorageService) GetAllCredentials(ctx context.Context) ([]entities.Credentials, error) {
	return s.credentialsRepo.GetAll(ctx)
}

// UpdateCredentials - обновить учётные данные
func (s *StorageService) UpdateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	return s.credentialsRepo.Update(ctx, entity)
}

// DeleteCredentials - удалить учётные данные
func (s *StorageService) DeleteCredentials(ctx context.Context, id string) error {
	return s.credentialsRepo.Delete(ctx, id)
}

// CreateText - создать запись с текстовыми данными
func (s *StorageService) CreateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	return s.textsRepo.Create(ctx, entity)
}

// GetText - получить текстовые данные
func (s *StorageService) GetText(ctx context.Context, id string) (*entities.TextData, error) {
	return s.textsRepo.Get(ctx, id)
}

// GetAllTexts - получить все текстовые данные
func (s *StorageService) GetAllTexts(ctx context.Context) ([]entities.TextData, error) {
	return s.textsRepo.GetAll(ctx)
}

// UpdateText - обновить текстовые данные
func (s *StorageService) UpdateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	return s.textsRepo.Update(ctx, entity)
}

// DeleteText - удалить текстовые данные
func (s *StorageService) DeleteText(ctx context.Context, id string) error {
	return s.textsRepo.Delete(ctx, id)
}
