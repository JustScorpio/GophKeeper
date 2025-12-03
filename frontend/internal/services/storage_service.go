// Сервис для работы с локальным хранилищем
package services

import (
	"context"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/frontend/internal/repositories"
)

type StorageService struct {
	binariesRepo    repositories.IRepository[entities.BinaryData]
	cardsRepo       repositories.IRepository[entities.CardInformation]
	credentialsRepo repositories.IRepository[entities.Credentials]
	textsRepo       repositories.IRepository[entities.TextData]
}

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

// Binary methods
func (s *StorageService) CreateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	return s.binariesRepo.Create(ctx, entity)
}

func (s *StorageService) GetBinary(ctx context.Context, id string) (*entities.BinaryData, error) {
	return s.binariesRepo.Get(ctx, id)
}

func (s *StorageService) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	return s.binariesRepo.GetAll(ctx)
}

func (s *StorageService) UpdateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	return s.binariesRepo.Update(ctx, entity)
}

func (s *StorageService) DeleteBinary(ctx context.Context, id string) error {
	return s.binariesRepo.Delete(ctx, id)
}

// Card methods
func (s *StorageService) CreateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	return s.cardsRepo.Create(ctx, entity)
}

func (s *StorageService) GetCard(ctx context.Context, id string) (*entities.CardInformation, error) {
	return s.cardsRepo.Get(ctx, id)
}

func (s *StorageService) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	return s.cardsRepo.GetAll(ctx)
}

func (s *StorageService) UpdateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	return s.cardsRepo.Update(ctx, entity)
}

func (s *StorageService) DeleteCard(ctx context.Context, id string) error {
	return s.cardsRepo.Delete(ctx, id)
}

// Credentials methods
func (s *StorageService) CreateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	return s.credentialsRepo.Create(ctx, entity)
}

func (s *StorageService) GetCredentials(ctx context.Context, id string) (*entities.Credentials, error) {
	return s.credentialsRepo.Get(ctx, id)
}

func (s *StorageService) GetAllCredentials(ctx context.Context) ([]entities.Credentials, error) {
	return s.credentialsRepo.GetAll(ctx)
}

func (s *StorageService) UpdateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	return s.credentialsRepo.Update(ctx, entity)
}

func (s *StorageService) DeleteCredentials(ctx context.Context, id string) error {
	return s.credentialsRepo.Delete(ctx, id)
}

// Text methods
func (s *StorageService) CreateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	return s.textsRepo.Create(ctx, entity)
}

func (s *StorageService) GetText(ctx context.Context, id string) (*entities.TextData, error) {
	return s.textsRepo.Get(ctx, id)
}

func (s *StorageService) GetAllTexts(ctx context.Context) ([]entities.TextData, error) {
	return s.textsRepo.GetAll(ctx)
}

func (s *StorageService) UpdateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	return s.textsRepo.Update(ctx, entity)
}

func (s *StorageService) DeleteText(ctx context.Context, id string) error {
	return s.textsRepo.Delete(ctx, id)
}
