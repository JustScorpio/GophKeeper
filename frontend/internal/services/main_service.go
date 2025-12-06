package services

import (
	"context"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/clients"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

type GophkeeperService struct {
	apiClient    clients.IAPIClient
	localStorage *StorageService
	syncService  *SyncService
}

func NewGophkeeperService(
	apiClient clients.IAPIClient,
	localStorage *StorageService,
	syncService *SyncService,
) *GophkeeperService {
	return &GophkeeperService{
		apiClient:    apiClient,
		localStorage: localStorage,
		syncService:  syncService,
	}
}

// Auth methods
func (s *GophkeeperService) Register(ctx context.Context, login, password string) error {
	return s.apiClient.Register(ctx, login, password)
}

func (s *GophkeeperService) Login(ctx context.Context, login, password string) error {
	err := s.apiClient.Login(ctx, login, password)
	if err != nil {
		return err
	}

	// Синхронизируем данные
	if err := s.syncService.Sync(ctx); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	return nil
}

// CreateBinary - создать бинарные данные (на клиенте и сервере)
func (s *GophkeeperService) CreateBinary(ctx context.Context, dto *dtos.NewBinaryData) (*entities.BinaryData, error) {
	// Создаем на сервере
	serverBinary, err := s.apiClient.CreateBinary(ctx, dto)
	if err != nil {
		return nil, err
	}

	// Создаем локально
	localBinary, err := s.localStorage.CreateBinary(ctx, serverBinary)
	if err != nil {
		// TODO: Что делать локальное создание не удалось? Маловероятно да и правильного решения как будто нет - поэтому игнорим.
		return serverBinary, fmt.Errorf("created on server but local failed: %w", err)
	}

	return localBinary, nil
}

// CreateCard - создать данные карты (на клиенте и сервере)
func (s *GophkeeperService) CreateCard(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error) {
	serverCard, err := s.apiClient.CreateCard(ctx, dto)
	if err != nil {
		return nil, err
	}

	localCard, err := s.localStorage.CreateCard(ctx, serverCard)
	if err != nil {
		return serverCard, fmt.Errorf("created on server but local failed: %w", err)
	}

	return localCard, nil
}

// CreateCredentials - создать учётные данные (на клиенте и сервере)
func (s *GophkeeperService) CreateCredentials(ctx context.Context, dto *dtos.NewCredentials) (*entities.Credentials, error) {
	serverCredentials, err := s.apiClient.CreateCredentials(ctx, dto)
	if err != nil {
		return nil, err
	}

	localCredentials, err := s.localStorage.CreateCredentials(ctx, serverCredentials)
	if err != nil {
		return serverCredentials, fmt.Errorf("created on server but local failed: %w", err)
	}

	return localCredentials, nil
}

// CreateText - создать данные карты (на клиенте и сервере)
func (s *GophkeeperService) CreateText(ctx context.Context, dto *dtos.NewTextData) (*entities.TextData, error) {
	serverText, err := s.apiClient.CreateText(ctx, dto)
	if err != nil {
		return nil, err
	}

	localText, err := s.localStorage.CreateText(ctx, serverText)
	if err != nil {
		return serverText, fmt.Errorf("created on server but local failed: %w", err)
	}

	return localText, nil
}

// GetBinary - получить бинарные данные (из локальной бд)
func (s *GophkeeperService) GetBinary(ctx context.Context, id string) (*entities.BinaryData, error) {
	return s.localStorage.GetBinary(ctx, id)
}

// GetAllBinaries - получить все бинарные данные (из локальной бд)
func (s *GophkeeperService) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	return s.localStorage.GetAllBinaries(ctx)
}

// GetCard - получить данные карты (из локальной бд)
func (s *GophkeeperService) GetCard(ctx context.Context, id string) (*entities.CardInformation, error) {
	return s.localStorage.GetCard(ctx, id)
}

// GetAllCards - получить данные всех карт (из локальной бд)
func (s *GophkeeperService) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	return s.localStorage.GetAllCards(ctx)
}

// GetCredentials - получить учётные данные (из локальной бд)
func (s *GophkeeperService) GetCredentials(ctx context.Context, id string) (*entities.Credentials, error) {
	return s.localStorage.GetCredentials(ctx, id)
}

// GetAllCredentials - получить всё учётные данные (из локальной бд)
func (s *GophkeeperService) GetAllCredentials(ctx context.Context) ([]entities.Credentials, error) {
	return s.localStorage.GetAllCredentials(ctx)
}

// GetText - получить текстовые данные (из локальной бд)
func (s *GophkeeperService) GetText(ctx context.Context, id string) (*entities.TextData, error) {
	return s.localStorage.GetText(ctx, id)
}

// GetAllTexts - получить все текстовые данные (из локальной бд)
func (s *GophkeeperService) GetAllTexts(ctx context.Context) ([]entities.TextData, error) {
	return s.localStorage.GetAllTexts(ctx)
}

// UpdateBinary - обновить бинарные данные
func (s *GophkeeperService) UpdateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	serverBinary, err := s.apiClient.UpdateBinary(ctx, entity)
	if err != nil {
		return nil, err
	}

	localBinary, err := s.localStorage.UpdateBinary(ctx, entity)
	if err != nil {
		return serverBinary, fmt.Errorf("updated on server but local failed: %w", err)
	}

	return localBinary, nil
}

// UpdateCard - обновить данные карты
func (s *GophkeeperService) UpdateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	serverCard, err := s.apiClient.UpdateCard(ctx, entity)
	if err != nil {
		return nil, err
	}

	localCard, err := s.localStorage.UpdateCard(ctx, entity)
	if err != nil {
		return serverCard, fmt.Errorf("updated on server but local failed: %w", err)
	}

	return localCard, nil
}

// UpdateCredentials - обновить учётные данные
func (s *GophkeeperService) UpdateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	serverCredentials, err := s.apiClient.UpdateCredentials(ctx, entity)
	if err != nil {
		return nil, err
	}

	localCredentials, err := s.localStorage.UpdateCredentials(ctx, entity)
	if err != nil {
		return serverCredentials, fmt.Errorf("updated on server but local failed: %w", err)
	}

	return localCredentials, nil
}

// UpdateText - обновить текстовые данные
func (s *GophkeeperService) UpdateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	serverText, err := s.apiClient.UpdateText(ctx, entity)
	if err != nil {
		return nil, err
	}

	localText, err := s.localStorage.UpdateText(ctx, entity)
	if err != nil {
		return serverText, fmt.Errorf("updated on server but local failed: %w", err)
	}

	return localText, nil
}

// DeleteBinary - удалить бинарные данные
func (s *GophkeeperService) DeleteBinary(ctx context.Context, id string) error {
	if err := s.apiClient.DeleteBinary(ctx, id); err != nil {
		return err
	}

	if err := s.localStorage.DeleteBinary(ctx, id); err != nil {
		return fmt.Errorf("deleted on server but local failed: %w", err)
	}

	return nil
}

// DeleteCard - удалить данные карты
func (s *GophkeeperService) DeleteCard(ctx context.Context, id string) error {
	if err := s.apiClient.DeleteCard(ctx, id); err != nil {
		return err
	}

	if err := s.localStorage.DeleteCard(ctx, id); err != nil {
		return fmt.Errorf("deleted on server but local failed: %w", err)
	}

	return nil
}

// DeleteCredentials - удалить учётные данные
func (s *GophkeeperService) DeleteCredentials(ctx context.Context, id string) error {
	if err := s.apiClient.DeleteCredentials(ctx, id); err != nil {
		return err
	}

	if err := s.localStorage.DeleteCredentials(ctx, id); err != nil {
		return fmt.Errorf("deleted on server but local failed: %w", err)
	}

	return nil
}

// DeleteText - удалить текстовые данные
func (s *GophkeeperService) DeleteText(ctx context.Context, id string) error {
	if err := s.apiClient.DeleteText(ctx, id); err != nil {
		return err
	}

	if err := s.localStorage.DeleteText(ctx, id); err != nil {
		return fmt.Errorf("deleted on server but local failed: %w", err)
	}

	return nil
}

// ForceSync - синхронизация данных с сервером
func (s *GophkeeperService) ForceSync(ctx context.Context) error {
	return s.syncService.Sync(ctx)
}
