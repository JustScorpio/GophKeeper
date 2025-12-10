package services

import (
	"context"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/clients"
	"github.com/JustScorpio/GophKeeper/frontend/internal/encryption"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// GophkeeperService - главный сервис клиентской части приложения
type GophkeeperService struct {
	apiClient     clients.IAPIClient
	localStorage  *StorageService
	syncService   *SyncService
	cryptoService *encryption.CryptoService
}

// NewGophkeeperService - создать главный сервис клиентской части приложения
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

// SetEncryption - установить сервис шифрование
func (s *GophkeeperService) SetEncryption(password string) error {
	s.cryptoService = encryption.NewCryptoService(password)
	return nil
}

// IsEncryptionSet - проверить установлен ли сервис шифрование
func (s *GophkeeperService) IsEncryptionSet() bool {
	return s.cryptoService != nil
}

// Register - регистрация
func (s *GophkeeperService) Register(ctx context.Context, login, password string) error {
	err := s.apiClient.Register(ctx, login, password)
	if err != nil {
		return err
	}

	return s.SetEncryption(password)
}

// Login - аутентификация
func (s *GophkeeperService) Login(ctx context.Context, login, password string) error {
	err := s.apiClient.Login(ctx, login, password)
	if err != nil {
		return err
	}

	// Синхронизируем данные
	if err := s.syncService.Sync(ctx); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	return s.SetEncryption(password)
}

// CreateBinary - создать бинарные данные (на клиенте и сервере)
func (s *GophkeeperService) CreateBinary(ctx context.Context, dto *dtos.NewBinaryData) (*entities.BinaryData, error) {
	// Шифруем DTO перед отправкой на сервер
	if err := dto.EncryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to encrypt binary DTO: %w", err)
	}

	// Создаем на сервере
	serverBinary, err := s.apiClient.CreateBinary(ctx, dto)
	if err != nil {
		return nil, err
	}

	// Создаем локально
	localBinary, err := s.localStorage.CreateBinary(ctx, serverBinary)
	if err != nil {
		return nil, fmt.Errorf("created on server but local failed: %w", err)
	}

	// Дешифруем локальную сущность для возврата
	if err := localBinary.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt local binary: %w", err)
	}

	return localBinary, nil
}

// CreateCard - создать данные карты (на клиенте и сервере)
func (s *GophkeeperService) CreateCard(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error) {
	// Шифруем DTO перед отправкой на сервер
	if err := dto.EncryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to encrypt card DTO: %w", err)
	}

	serverCard, err := s.apiClient.CreateCard(ctx, dto)
	if err != nil {
		return nil, err
	}

	localCard, err := s.localStorage.CreateCard(ctx, serverCard)
	if err != nil {
		return nil, fmt.Errorf("created on server but local failed: %w", err)
	}

	// Дешифруем локальную сущность для возврата
	if err := localCard.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt local card: %w", err)
	}

	return localCard, nil
}

// CreateCredentials - создать учётные данные (на клиенте и сервере)
func (s *GophkeeperService) CreateCredentials(ctx context.Context, dto *dtos.NewCredentials) (*entities.Credentials, error) {
	// Шифруем только метаданные в DTO
	if err := dto.EncryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to encrypt credentials DTO: %w", err)
	}

	serverCredentials, err := s.apiClient.CreateCredentials(ctx, dto)
	if err != nil {
		return nil, err
	}

	localCredentials, err := s.localStorage.CreateCredentials(ctx, serverCredentials)
	if err != nil {
		return nil, fmt.Errorf("created on server but local failed: %w", err)
	}

	// Дешифруем локальную сущность для возврата
	if err := localCredentials.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt local credentials: %w", err)
	}

	return localCredentials, nil
}

// CreateText - создать текстовые данные (на клиенте и сервере)
func (s *GophkeeperService) CreateText(ctx context.Context, dto *dtos.NewTextData) (*entities.TextData, error) {
	// Шифруем DTO перед отправкой на сервер
	if err := dto.EncryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to encrypt text DTO: %w", err)
	}

	serverText, err := s.apiClient.CreateText(ctx, dto)
	if err != nil {
		return nil, err
	}

	localText, err := s.localStorage.CreateText(ctx, serverText)
	if err != nil {
		return nil, fmt.Errorf("created on server but local failed: %w", err)
	}

	// Дешифруем локальную сущность для возврата
	if err := localText.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt local text: %w", err)
	}

	return localText, nil
}

// GetBinary - получить бинарные данные (из локальной бд)
func (s *GophkeeperService) GetBinary(ctx context.Context, id string) (*entities.BinaryData, error) {
	binary, err := s.localStorage.GetBinary(ctx, id)
	if err != nil {
		return nil, err
	}

	// Дешифруем данные из локального хранилища
	if err := binary.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt binary: %w", err)
	}

	return binary, nil
}

// GetAllBinaries - получить все бинарные данные (из локальной бд)
func (s *GophkeeperService) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	binaries, err := s.localStorage.GetAllBinaries(ctx)
	if err != nil {
		return nil, err
	}

	// Дешифруем все бинарные данные
	for i := range binaries {
		if err := binaries[i].DecryptFields(s.cryptoService); err != nil {
			return nil, fmt.Errorf("failed to decrypt binary %s: %w", binaries[i].ID, err)
		}
	}

	return binaries, nil
}

// GetCard - получить данные карты (из локальной бд)
func (s *GophkeeperService) GetCard(ctx context.Context, id string) (*entities.CardInformation, error) {
	card, err := s.localStorage.GetCard(ctx, id)
	if err != nil {
		return nil, err
	}

	// Дешифруем данные из локального хранилища
	if err := card.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt card: %w", err)
	}

	return card, nil
}

// GetAllCards - получить данные всех карт (из локальной бд)
func (s *GophkeeperService) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	cards, err := s.localStorage.GetAllCards(ctx)
	if err != nil {
		return nil, err
	}

	// Дешифруем все данные карт
	for i := range cards {
		if err := cards[i].DecryptFields(s.cryptoService); err != nil {
			return nil, fmt.Errorf("failed to decrypt card %s: %w", cards[i].ID, err)
		}
	}

	return cards, nil
}

// GetCredentials - получить учётные данные (из локальной бд)
func (s *GophkeeperService) GetCredentials(ctx context.Context, id string) (*entities.Credentials, error) {
	credentials, err := s.localStorage.GetCredentials(ctx, id)
	if err != nil {
		return nil, err
	}

	// Дешифруем данные из локального хранилища
	if err := credentials.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt credentials: %w", err)
	}

	return credentials, nil
}

// GetAllCredentials - получить всё учётные данные (из локальной бд)
func (s *GophkeeperService) GetAllCredentials(ctx context.Context) ([]entities.Credentials, error) {
	credentials, err := s.localStorage.GetAllCredentials(ctx)
	if err != nil {
		return nil, err
	}

	// Дешифруем все учётные данные
	for i := range credentials {
		if err := credentials[i].DecryptFields(s.cryptoService); err != nil {
			return nil, fmt.Errorf("failed to decrypt credentials %s: %w", credentials[i].ID, err)
		}
	}

	return credentials, nil
}

// GetText - получить текстовые данные (из локальной бд)
func (s *GophkeeperService) GetText(ctx context.Context, id string) (*entities.TextData, error) {
	text, err := s.localStorage.GetText(ctx, id)
	if err != nil {
		return nil, err
	}

	// Дешифруем данные из локального хранилища
	if err := text.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt text: %w", err)
	}

	return text, nil
}

// GetAllTexts - получить все текстовые данные (из локальной бд)
func (s *GophkeeperService) GetAllTexts(ctx context.Context) ([]entities.TextData, error) {
	texts, err := s.localStorage.GetAllTexts(ctx)
	if err != nil {
		return nil, err
	}

	// Дешифруем все текстовые данные
	for i := range texts {
		if err := texts[i].DecryptFields(s.cryptoService); err != nil {
			return nil, fmt.Errorf("failed to decrypt text %s: %w", texts[i].ID, err)
		}
	}

	return texts, nil
}

// UpdateBinary - обновить бинарные данные
func (s *GophkeeperService) UpdateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	// Шифруем перед отправкой на сервер
	if err := entity.EncryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to encrypt binary for update: %w", err)
	}

	serverBinary, err := s.apiClient.UpdateBinary(ctx, entity)
	if err != nil {
		return nil, err
	}

	localBinary, err := s.localStorage.UpdateBinary(ctx, serverBinary)
	if err != nil {
		return serverBinary, fmt.Errorf("updated on server but local failed: %w", err)
	}

	// Дешифруем локальную сущность для возврата
	if err := localBinary.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt local binary: %w", err)
	}

	return localBinary, nil
}

// UpdateCard - обновить данные карты
func (s *GophkeeperService) UpdateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	// Шифруем перед отправкой на сервер
	if err := entity.EncryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to encrypt card for update: %w", err)
	}

	serverCard, err := s.apiClient.UpdateCard(ctx, entity)
	if err != nil {
		return nil, err
	}

	localCard, err := s.localStorage.UpdateCard(ctx, serverCard)
	if err != nil {
		return serverCard, fmt.Errorf("updated on server but local failed: %w", err)
	}

	// Дешифруем локальную сущность для возврата
	if err := localCard.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt local card: %w", err)
	}

	return localCard, nil
}

// UpdateCredentials - обновить учётные данные
func (s *GophkeeperService) UpdateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	// Шифруем метаданные перед отправкой на сервер
	if err := entity.EncryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to encrypt credentials for update: %w", err)
	}

	serverCredentials, err := s.apiClient.UpdateCredentials(ctx, entity)
	if err != nil {
		return nil, err
	}

	localCredentials, err := s.localStorage.UpdateCredentials(ctx, serverCredentials)
	if err != nil {
		return serverCredentials, fmt.Errorf("updated on server but local failed: %w", err)
	}

	// Дешифруем локальную сущность для возврата
	if err := localCredentials.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt local credentials: %w", err)
	}

	return localCredentials, nil
}

// UpdateText - обновить текстовые данные
func (s *GophkeeperService) UpdateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	// Шифруем перед отправкой на сервер
	if err := entity.EncryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to encrypt text for update: %w", err)
	}

	serverText, err := s.apiClient.UpdateText(ctx, entity)
	if err != nil {
		return nil, err
	}

	localText, err := s.localStorage.UpdateText(ctx, serverText)
	if err != nil {
		return serverText, fmt.Errorf("updated on server but local failed: %w", err)
	}

	// Дешифруем локальную сущность для возврата
	if err := localText.DecryptFields(s.cryptoService); err != nil {
		return nil, fmt.Errorf("failed to decrypt local text: %w", err)
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
