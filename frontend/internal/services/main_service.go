package services

import (
	"context"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

type GophkeeperService struct {
	authService  *AuthService
	apiClient    *APIClient
	localStorage *StorageService
	syncService  *SyncService
}

func NewGophkeeperService(
	authService *AuthService,
	apiClient *APIClient,
	localStorage *StorageService,
	syncService *SyncService,
) *GophkeeperService {
	return &GophkeeperService{
		authService:  authService,
		apiClient:    apiClient,
		localStorage: localStorage,
		syncService:  syncService,
	}
}

// Auth methods
func (s *GophkeeperService) Register(ctx context.Context, login, password string) (*entities.User, error) {
	return s.authService.Register(ctx, login, password)
}

func (s *GophkeeperService) Login(ctx context.Context, login, password string) (*entities.User, error) {
	user, err := s.authService.Login(ctx, login, password)
	if err != nil {
		return nil, err
	}

	// После успешного логина синхронизируем данные
	if err := s.syncService.SyncOnLogin(ctx); err != nil {
		return nil, fmt.Errorf("login successful but sync failed: %w", err)
	}

	return user, nil
}

// Create methods - создаем и локально, и на сервере
func (s *GophkeeperService) CreateBinary(ctx context.Context, dto *dtos.NewBinaryData) (*entities.BinaryData, error) {
	// Создаем на сервере
	serverBinary, err := s.apiClient.CreateBinary(ctx, dto)
	if err != nil {
		return nil, err
	}

	// Создаем локально
	localBinary, err := s.localStorage.CreateBinary(ctx, dto)
	if err != nil {
		// Если локальное создание не удалось, можно откатить серверное
		// или просто залогировать ошибку
		return serverBinary, fmt.Errorf("created on server but local failed: %w", err)
	}

	// Используем ID с сервера для локальной записи
	return &entities.BinaryData{
		Data:     localBinary.Data,
		Metadata: localBinary.Metadata,
	}, nil
}

func (s *GophkeeperService) CreateCard(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error) {
	serverCard, err := s.apiClient.CreateCard(ctx, dto)
	if err != nil {
		return nil, err
	}

	localCard, err := s.localStorage.CreateCard(ctx, dto)
	if err != nil {
		return serverCard, fmt.Errorf("created on server but local failed: %w", err)
	}

	return &entities.CardInformation{
		ID:             serverCard.ID,
		Number:         localCard.Number,
		CardHolder:     localCard.CardHolder,
		ExpirationDate: localCard.ExpirationDate,
		CVV:            localCard.CVV,
		Metadata:       localCard.Metadata,
	}, nil
}

// Read methods - читаем только из локальной БД (данные уже синхронизированы)
func (s *GophkeeperService) GetBinary(ctx context.Context, id string) (*entities.BinaryData, error) {
	return s.localStorage.GetBinary(ctx, id)
}

func (s *GophkeeperService) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	return s.localStorage.GetAllBinaries(ctx)
}

func (s *GophkeeperService) GetCard(ctx context.Context, id string) (*entities.CardInformation, error) {
	return s.localStorage.GetCard(ctx, id)
}

func (s *GophkeeperService) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	return s.localStorage.GetAllCards(ctx)
}

// Update methods - обновляем и локально, и на сервере
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

// Delete methods - удаляем и локально, и на сервере
func (s *GophkeeperService) DeleteBinary(ctx context.Context, id string) error {
	if err := s.apiClient.DeleteBinary(ctx, id); err != nil {
		return err
	}

	if err := s.localStorage.DeleteBinary(ctx, id); err != nil {
		return fmt.Errorf("deleted on server but local failed: %w", err)
	}

	return nil
}

// Manual sync - принудительная синхронизация
func (s *GophkeeperService) ForceSync(ctx context.Context) error {
	return s.syncService.SyncOnLogin(ctx)
}
