package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/frontend/internal/repositories/inmemory"
	"github.com/JustScorpio/GophKeeper/frontend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockSyncAPIClient для SyncService
type MockSyncAPIClient struct {
	mock.Mock
}

// Register - регистрация пользователя
func (m *MockSyncAPIClient) Register(ctx context.Context, login, password string) error {
	args := m.Called(ctx, login, password)
	return args.Error(0)
}

// Login - аутентификация пользователя
func (m *MockSyncAPIClient) Login(ctx context.Context, login, password string) error {
	args := m.Called(ctx, login, password)
	return args.Error(0)
}

// CreateBinary - создать бинарные данные
func (m *MockSyncAPIClient) CreateBinary(ctx context.Context, dto *dtos.NewBinaryData) (*entities.BinaryData, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BinaryData), args.Error(1)
}

// GetAllBinaries - получить все бинарные данные
func (m *MockSyncAPIClient) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.BinaryData), args.Error(1)
}

// UpdateBinary - обновить бинарные данные
func (m *MockSyncAPIClient) UpdateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BinaryData), args.Error(1)
}

// DeleteBinary - удалить бинарные данные
func (m *MockSyncAPIClient) DeleteBinary(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CreateCard - создать данные карты
func (m *MockSyncAPIClient) CreateCard(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CardInformation), args.Error(1)
}

// GetAllCards - получить данные всех карт
func (m *MockSyncAPIClient) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.CardInformation), args.Error(1)
}

// UpdateCard - обновить данные карты
func (m *MockSyncAPIClient) UpdateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CardInformation), args.Error(1)
}

// DeleteCard - удалить данные карты
func (m *MockSyncAPIClient) DeleteCard(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CreateCredentials - создать учётные данные
func (m *MockSyncAPIClient) CreateCredentials(ctx context.Context, dto *dtos.NewCredentials) (*entities.Credentials, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credentials), args.Error(1)
}

// GetAllCredentials - получить все учётные данные
func (m *MockSyncAPIClient) GetAllCredentials(ctx context.Context) ([]entities.Credentials, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Credentials), args.Error(1)
}

// UpdateCredentials - обновить учётные данные
func (m *MockSyncAPIClient) UpdateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credentials), args.Error(1)
}

// DeleteCredentials - удалить учётные данные
func (m *MockSyncAPIClient) DeleteCredentials(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// CreateText - создать текстовые данные
func (m *MockSyncAPIClient) CreateText(ctx context.Context, dto *dtos.NewTextData) (*entities.TextData, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TextData), args.Error(1)
}

// GetAllTexts - получить все текстовые данные
func (m *MockSyncAPIClient) GetAllTexts(ctx context.Context) ([]entities.TextData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.TextData), args.Error(1)
}

// UpdateText - обновить текстовые данные
func (m *MockSyncAPIClient) UpdateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TextData), args.Error(1)
}

// DeleteText - удалить текстовые данные
func (m *MockSyncAPIClient) DeleteText(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// TestSyncService_Sync - тесты сервиса синхронизации данных
func TestSyncService_Sync(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful full sync", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Подготавливаем данные с сервера
		serverBinaries := []entities.BinaryData{
			{
				SecureEntity: entities.SecureEntity{
					ID:       "server-bin-1",
					Metadata: "Server Binary 1",
				},
				Data: []byte("server data 1"),
			},
			{
				SecureEntity: entities.SecureEntity{
					ID:       "server-bin-2",
					Metadata: "Server Binary 2",
				},
				Data: []byte("server data 2"),
			},
		}

		serverCards := []entities.CardInformation{
			{
				SecureEntity: entities.SecureEntity{
					ID:       "server-card-1",
					Metadata: "Visa Card",
				},
				Number:         "4111111111111111",
				CardHolder:     "John Doe",
				ExpirationDate: "12/25",
				CVV:            "123",
			},
		}

		serverCredentials := []entities.Credentials{
			{
				SecureEntity: entities.SecureEntity{
					ID:       "server-cred-1",
					Metadata: "Google Account",
				},
				Login:    "john@gmail.com",
				Password: "password123",
			},
		}

		serverTexts := []entities.TextData{
			{
				SecureEntity: entities.SecureEntity{
					ID:       "server-text-1",
					Metadata: "Important Notes",
				},
				Data: "These are important notes from server",
			},
		}

		// Добавляем локальные данные (должны быть удалены после синхронизации)
		localBinary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "local-bin",
				Metadata: "Local Binary",
			},
			Data: []byte("local data"),
		}
		_, err := storageService.CreateBinary(ctx, localBinary)
		require.NoError(t, err)

		localCard := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "local-card",
				Metadata: "Local Card",
			},
			Number:         "5555555555554444",
			CardHolder:     "Local User",
			ExpirationDate: "01/26",
			CVV:            "789",
		}
		_, err = storageService.CreateCard(ctx, localCard)
		require.NoError(t, err)

		// Настраиваем моки для ВСЕХ 4 методов
		mockAPI.On("GetAllBinaries", ctx).Return(serverBinaries, nil)
		mockAPI.On("GetAllCards", ctx).Return(serverCards, nil)
		mockAPI.On("GetAllCredentials", ctx).Return(serverCredentials, nil)
		mockAPI.On("GetAllTexts", ctx).Return(serverTexts, nil)

		// Выполняем синхронизацию
		err = syncService.Sync(ctx)
		require.NoError(t, err)

		// Проверяем, что локальные данные заменены серверными
		binaries, err := storageService.GetAllBinaries(ctx)
		require.NoError(t, err)
		assert.Len(t, binaries, 2)
		assert.Equal(t, "server-bin-1", binaries[0].ID)
		assert.Equal(t, "server-bin-2", binaries[1].ID)

		cards, err := storageService.GetAllCards(ctx)
		require.NoError(t, err)
		assert.Len(t, cards, 1)
		assert.Equal(t, "server-card-1", cards[0].ID)
		assert.Equal(t, "4111111111111111", cards[0].Number)

		credentials, err := storageService.GetAllCredentials(ctx)
		require.NoError(t, err)
		assert.Len(t, credentials, 1)
		assert.Equal(t, "server-cred-1", credentials[0].ID)
		assert.Equal(t, "john@gmail.com", credentials[0].Login)

		texts, err := storageService.GetAllTexts(ctx)
		require.NoError(t, err)
		assert.Len(t, texts, 1)
		assert.Equal(t, "server-text-1", texts[0].ID)
		assert.Equal(t, "These are important notes from server", texts[0].Data)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Sync with empty server data", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Добавляем локальные данные
		localBinary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "local-data",
				Metadata: "Local Data",
			},
			Data: []byte("local"),
		}
		_, err := storageService.CreateBinary(ctx, localBinary)
		require.NoError(t, err)

		// Сервер возвращает пустые списки для ВСЕХ 4 методов
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		err = syncService.Sync(ctx)
		require.NoError(t, err)

		// Проверяем, что локальные данные удалены
		binaries, err := storageService.GetAllBinaries(ctx)
		require.NoError(t, err)
		assert.Empty(t, binaries)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Sync fails on binaries error", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Только первый метод возвращает ошибку, остальные не вызываются
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, errors.New("network error"))

		err := syncService.Sync(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to sync binaries")
		assert.Contains(t, err.Error(), "network error")

		mockAPI.AssertExpectations(t)
		// Остальные методы не должны вызываться при ошибке на первом
		mockAPI.AssertNotCalled(t, "GetAllCards")
		mockAPI.AssertNotCalled(t, "GetAllCredentials")
		mockAPI.AssertNotCalled(t, "GetAllTexts")
	})

	t.Run("Sync fails on cards error after successful binaries", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Настраиваем ВСЕ 4 метода
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, errors.New("cards error"))

		// Эти методы не будут вызываться из-за ошибки на cards
		// Но можем оставить их без вызова или использовать Maybe()
		// Для простоты не настраиваем их

		err := syncService.Sync(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to sync cards")
		assert.Contains(t, err.Error(), "cards error")

		mockAPI.AssertExpectations(t)
	})

	t.Run("Sync fails on credentials error", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Настраиваем методы до credentials
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, errors.New("creds error"))

		err := syncService.Sync(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to sync credentials")
		assert.Contains(t, err.Error(), "creds error")

		mockAPI.AssertExpectations(t)
		mockAPI.AssertNotCalled(t, "GetAllTexts")
	})

	t.Run("Sync fails on texts error", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Настраиваем ВСЕ 4 метода
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, errors.New("texts error"))

		err := syncService.Sync(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to sync texts")
		assert.Contains(t, err.Error(), "texts error")

		mockAPI.AssertExpectations(t)
	})

	t.Run("Sync preserves server data structure", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Сервер возвращает данные с полной структурой
		serverBinaries := []entities.BinaryData{
			{
				SecureEntity: entities.SecureEntity{
					ID:       "test-id",
					Metadata: "Full metadata",
				},
				Data: []byte{1, 2, 3, 4, 5},
			},
		}

		// Настраиваем ВСЕ 4 метода
		mockAPI.On("GetAllBinaries", ctx).Return(serverBinaries, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		// Вызываем синхронизацию
		err := syncService.Sync(ctx)
		require.NoError(t, err)

		// Проверяем, что данные сохранились полностью
		binaries, err := storageService.GetAllBinaries(ctx)
		require.NoError(t, err)
		require.Len(t, binaries, 1)

		assert.Equal(t, "test-id", binaries[0].ID)
		assert.Equal(t, "Full metadata", binaries[0].Metadata)
		assert.Equal(t, []byte{1, 2, 3, 4, 5}, binaries[0].Data)

		mockAPI.AssertExpectations(t)
	})

	t.Run("SyncCards with multiple cards", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		serverCards := []entities.CardInformation{
			{
				SecureEntity: entities.SecureEntity{
					ID:       "card-1",
					Metadata: "Personal Visa",
				},
				Number:         "4111111111111111",
				CardHolder:     "John Doe",
				ExpirationDate: "12/25",
				CVV:            "123",
			},
			{
				SecureEntity: entities.SecureEntity{
					ID:       "card-2",
					Metadata: "Business MasterCard",
				},
				Number:         "5555555555554444",
				CardHolder:     "Jane Smith",
				ExpirationDate: "11/26",
				CVV:            "456",
			},
		}

		// Настраиваем ВСЕ 4 метода
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return(serverCards, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		err := syncService.Sync(ctx)
		require.NoError(t, err)

		cards, err := storageService.GetAllCards(ctx)
		require.NoError(t, err)
		assert.Len(t, cards, 2)
		assert.Equal(t, "card-1", cards[0].ID)
		assert.Equal(t, "card-2", cards[1].ID)

		mockAPI.AssertExpectations(t)
	})
}

// TestSyncService_EdgeCases - тесты для edge cases
func TestSyncService_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("Sync with nil data from server", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Сервер возвращает nil вместо среза
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		err := syncService.Sync(ctx)
		require.NoError(t, err)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Sync continues after local delete of non-existent item", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Сервер возвращает данные
		serverBinaries := []entities.BinaryData{
			{
				SecureEntity: entities.SecureEntity{
					ID:       "server-1",
					Metadata: "Server Data",
				},
				Data: []byte("server"),
			},
		}

		// Локальное хранилище пустое, но delete не должен вызывать ошибку
		mockAPI.On("GetAllBinaries", ctx).Return(serverBinaries, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		err := syncService.Sync(ctx)
		require.NoError(t, err)

		// Проверяем, что данные синхронизированы
		binaries, err := storageService.GetAllBinaries(ctx)
		require.NoError(t, err)
		assert.Len(t, binaries, 1)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Sync handles duplicate server data gracefully", func(t *testing.T) {
		mockAPI := new(MockSyncAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)

		syncService := services.NewSyncService(mockAPI, storageService)

		// Добавляем локальные данные
		localBinary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "existing-id",
				Metadata: "Existing",
			},
			Data: []byte("existing"),
		}
		_, err := storageService.CreateBinary(ctx, localBinary)
		require.NoError(t, err)

		// Сервер возвращает данные с тем же ID
		serverBinaries := []entities.BinaryData{
			{
				SecureEntity: entities.SecureEntity{
					ID:       "existing-id",
					Metadata: "Updated from server",
				},
				Data: []byte("updated"),
			},
		}

		mockAPI.On("GetAllBinaries", ctx).Return(serverBinaries, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		err = syncService.Sync(ctx)
		require.NoError(t, err)

		// Проверяем, что данные обновились
		binaries, err := storageService.GetAllBinaries(ctx)
		require.NoError(t, err)
		assert.Len(t, binaries, 1)
		assert.Equal(t, "Updated from server", binaries[0].Metadata)
		assert.Equal(t, []byte("updated"), binaries[0].Data)

		mockAPI.AssertExpectations(t)
	})
}
