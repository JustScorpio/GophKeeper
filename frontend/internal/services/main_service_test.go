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

// MockAPIClient для GophkeeperService
type MockGophKeeperAPIClient struct {
	mock.Mock
}

func (m *MockGophKeeperAPIClient) Register(ctx context.Context, login, password string) error {
	args := m.Called(ctx, login, password)
	return args.Error(0)
}

func (m *MockGophKeeperAPIClient) Login(ctx context.Context, login, password string) error {
	args := m.Called(ctx, login, password)
	return args.Error(0)
}

func (m *MockGophKeeperAPIClient) CreateBinary(ctx context.Context, dto *dtos.NewBinaryData) (*entities.BinaryData, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BinaryData), args.Error(1)
}

func (m *MockGophKeeperAPIClient) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.BinaryData), args.Error(1)
}

func (m *MockGophKeeperAPIClient) UpdateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BinaryData), args.Error(1)
}

func (m *MockGophKeeperAPIClient) DeleteBinary(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGophKeeperAPIClient) CreateCard(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CardInformation), args.Error(1)
}

func (m *MockGophKeeperAPIClient) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.CardInformation), args.Error(1)
}

func (m *MockGophKeeperAPIClient) UpdateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CardInformation), args.Error(1)
}

func (m *MockGophKeeperAPIClient) DeleteCard(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGophKeeperAPIClient) CreateCredentials(ctx context.Context, dto *dtos.NewCredentials) (*entities.Credentials, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credentials), args.Error(1)
}

func (m *MockGophKeeperAPIClient) GetAllCredentials(ctx context.Context) ([]entities.Credentials, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Credentials), args.Error(1)
}

func (m *MockGophKeeperAPIClient) UpdateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credentials), args.Error(1)
}

func (m *MockGophKeeperAPIClient) DeleteCredentials(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGophKeeperAPIClient) CreateText(ctx context.Context, dto *dtos.NewTextData) (*entities.TextData, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TextData), args.Error(1)
}

func (m *MockGophKeeperAPIClient) GetAllTexts(ctx context.Context) ([]entities.TextData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.TextData), args.Error(1)
}

func (m *MockGophKeeperAPIClient) UpdateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	args := m.Called(ctx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TextData), args.Error(1)
}

func (m *MockGophKeeperAPIClient) DeleteText(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGophkeeperService_Auth(t *testing.T) {
	ctx := context.Background()

	t.Run("Successful registration", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		mockAPI.On("Register", ctx, "testuser", "testpass123").Return(nil)

		err := gophkeeperService.Register(ctx, "testuser", "testpass123")
		require.NoError(t, err)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Registration failed", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		mockAPI.On("Register", ctx, "existing", "pass").Return(errors.New("user already exists"))

		err := gophkeeperService.Register(ctx, "existing", "pass")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "user already exists")

		mockAPI.AssertExpectations(t)
	})

	t.Run("Successful login with sync", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		// Настраиваем успешный логин
		mockAPI.On("Login", ctx, "user", "correctpass").Return(nil)

		// Настраиваем успешную синхронизацию
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		err := gophkeeperService.Login(ctx, "user", "correctpass")
		require.NoError(t, err)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Login failed - authentication error", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		mockAPI.On("Login", ctx, "user", "wrongpass").Return(errors.New("invalid credentials"))

		err := gophkeeperService.Login(ctx, "user", "wrongpass")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")

		// Sync не должен вызываться при ошибке логина
		mockAPI.AssertNotCalled(t, "GetAllBinaries")
		mockAPI.AssertNotCalled(t, "GetAllCards")
		mockAPI.AssertNotCalled(t, "GetAllCredentials")
		mockAPI.AssertNotCalled(t, "GetAllTexts")
	})

	t.Run("Login successful but sync failed", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		mockAPI.On("Login", ctx, "user", "pass").Return(nil)
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, errors.New("sync error"))

		err := gophkeeperService.Login(ctx, "user", "pass")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sync failed")
		assert.Contains(t, err.Error(), "sync error")

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_CreateOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("Create binary - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		dto := &dtos.NewBinaryData{
			Data:            []byte("test data"),
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Test binary"},
		}

		serverResponse := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "server-id-123",
				Metadata: "Test binary",
			},
			Data: []byte("test data"),
		}

		mockAPI.On("CreateBinary", ctx, dto).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateBinary(ctx, dto)
		require.NoError(t, err)
		assert.Equal(t, "server-id-123", result.ID)
		assert.Equal(t, []byte("test data"), result.Data)

		// Проверяем, что данные сохранились локально
		localData, err := storageService.GetBinary(ctx, "server-id-123")
		require.NoError(t, err)
		assert.NotNil(t, localData)
		assert.Equal(t, "server-id-123", localData.ID)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Create binary - API error", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		dto := &dtos.NewBinaryData{
			Data: []byte("test data"),
		}

		mockAPI.On("CreateBinary", ctx, dto).Return((*entities.BinaryData)(nil), errors.New("api error"))

		result, err := gophkeeperService.CreateBinary(ctx, dto)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "api error")

		mockAPI.AssertExpectations(t)
	})

	t.Run("Create card - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		dto := &dtos.NewCardInformation{
			Number:          "4111111111111111",
			CardHolder:      "John Doe",
			ExpirationDate:  "12/25",
			CVV:             "123",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Visa"},
		}

		serverResponse := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "card-id-123",
				Metadata: "Visa",
			},
			Number:         "4111111111111111",
			CardHolder:     "John Doe",
			ExpirationDate: "12/25",
			CVV:            "123",
		}

		mockAPI.On("CreateCard", ctx, dto).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateCard(ctx, dto)
		require.NoError(t, err)
		assert.Equal(t, "card-id-123", result.ID)
		assert.Equal(t, "4111111111111111", result.Number)

		// Проверяем локальное сохранение
		localCard, err := storageService.GetCard(ctx, "card-id-123")
		require.NoError(t, err)
		assert.NotNil(t, localCard)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Create credentials - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		dto := &dtos.NewCredentials{
			Login:           "user@example.com",
			Password:        "securepass123",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Email account"},
		}

		serverResponse := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "cred-id-123",
				Metadata: "Email account",
			},
			Login:    "user@example.com",
			Password: "securepass123",
		}

		mockAPI.On("CreateCredentials", ctx, dto).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateCredentials(ctx, dto)
		require.NoError(t, err)
		assert.Equal(t, "cred-id-123", result.ID)
		assert.Equal(t, "user@example.com", result.Login)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Create text - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		dto := &dtos.NewTextData{
			Data:            "This is important text",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Notes"},
		}

		serverResponse := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       "text-id-123",
				Metadata: "Notes",
			},
			Data: "This is important text",
		}

		mockAPI.On("CreateText", ctx, dto).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateText(ctx, dto)
		require.NoError(t, err)
		assert.Equal(t, "text-id-123", result.ID)
		assert.Equal(t, "This is important text", result.Data)

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_ReadOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("Get all binaries from local storage", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		// Добавляем данные в локальное хранилище
		binary1 := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "id1"},
			Data:         []byte("data1"),
		}
		binary2 := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "id2"},
			Data:         []byte("data2"),
		}

		_, err := storageService.CreateBinary(ctx, binary1)
		require.NoError(t, err)
		_, err = storageService.CreateBinary(ctx, binary2)
		require.NoError(t, err)

		result, err := gophkeeperService.GetAllBinaries(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		expectedIDs := []string{"id1", "id2"}
		actualIDs := []string{result[0].ID, result[1].ID}
		assert.ElementsMatch(t, expectedIDs, actualIDs)

		// API не должно вызываться для операций чтения
		mockAPI.AssertNotCalled(t, "GetAllBinaries")
	})

	t.Run("Get single binary by ID", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "test-id"},
			Data:         []byte("test data"),
		}
		_, err := storageService.CreateBinary(ctx, binary)
		require.NoError(t, err)

		result, err := gophkeeperService.GetBinary(ctx, "test-id")
		require.NoError(t, err)
		assert.Equal(t, "test-id", result.ID)
		assert.Equal(t, []byte("test data"), result.Data)

		mockAPI.AssertNotCalled(t, "GetBinary")
	})

	t.Run("Get all cards", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		card := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{ID: "card1"},
			Number:       "4111111111111111",
		}
		_, err := storageService.CreateCard(ctx, card)
		require.NoError(t, err)

		result, err := gophkeeperService.GetAllCards(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "card1", result[0].ID)

		mockAPI.AssertNotCalled(t, "GetAllCards")
	})

	t.Run("Get all credentials", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		creds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{ID: "creds1"},
			Login:        "user@test.com",
		}
		_, err := storageService.CreateCredentials(ctx, creds)
		require.NoError(t, err)

		result, err := gophkeeperService.GetAllCredentials(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "creds1", result[0].ID)

		mockAPI.AssertNotCalled(t, "GetAllCredentials")
	})

	t.Run("Get all texts", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		text := &entities.TextData{
			SecureEntity: entities.SecureEntity{ID: "text1"},
			Data:         "test text",
		}
		_, err := storageService.CreateText(ctx, text)
		require.NoError(t, err)

		result, err := gophkeeperService.GetAllTexts(ctx)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "text1", result[0].ID)

		mockAPI.AssertNotCalled(t, "GetAllTexts")
	})

	t.Run("Get non-existent entity returns nil", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		result, err := gophkeeperService.GetBinary(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, result)

		card, err := gophkeeperService.GetCard(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, card)

		creds, err := gophkeeperService.GetCredentials(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, creds)

		text, err := gophkeeperService.GetText(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, text)
	})
}

func TestGophkeeperService_UpdateOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("Update binary - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		// Сначала создаем запись
		existing := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "update-id"},
			Data:         []byte("old data"),
		}
		_, err := storageService.CreateBinary(ctx, existing)
		require.NoError(t, err)

		// Обновляем
		updatedEntity := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "update-id",
				Metadata: "Updated",
			},
			Data: []byte("new data"),
		}

		serverResponse := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "update-id",
				Metadata: "Updated",
			},
			Data: []byte("new data"),
		}

		mockAPI.On("UpdateBinary", ctx, updatedEntity).Return(serverResponse, nil)

		result, err := gophkeeperService.UpdateBinary(ctx, updatedEntity)
		require.NoError(t, err)
		assert.Equal(t, "Updated", result.Metadata)
		assert.Equal(t, []byte("new data"), result.Data)

		// Проверяем локальное обновление
		localData, err := storageService.GetBinary(ctx, "update-id")
		require.NoError(t, err)
		assert.Equal(t, "Updated", localData.Metadata)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Update binary - API error", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		entity := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "non-existent"},
			Data:         []byte("data"),
		}

		mockAPI.On("UpdateBinary", ctx, entity).Return((*entities.BinaryData)(nil), errors.New("not found"))

		result, err := gophkeeperService.UpdateBinary(ctx, entity)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")

		mockAPI.AssertExpectations(t)
	})

	t.Run("Update card - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		// Создаем карту
		card := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{ID: "card-id"},
			Number:       "4111111111111111",
		}
		_, err := storageService.CreateCard(ctx, card)
		require.NoError(t, err)

		// Обновляем
		updatedCard := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "card-id",
				Metadata: "Updated Visa",
			},
			Number:         "5555555555554444",
			CardHolder:     "Jane Doe",
			ExpirationDate: "12/26",
			CVV:            "456",
		}

		serverResponse := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "card-id",
				Metadata: "Updated Visa",
			},
			Number:         "5555555555554444",
			CardHolder:     "Jane Doe",
			ExpirationDate: "12/26",
			CVV:            "456",
		}

		mockAPI.On("UpdateCard", ctx, updatedCard).Return(serverResponse, nil)

		result, err := gophkeeperService.UpdateCard(ctx, updatedCard)
		require.NoError(t, err)
		assert.Equal(t, "Updated Visa", result.Metadata)
		assert.Equal(t, "5555555555554444", result.Number)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Update credentials - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		creds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{ID: "creds-id"},
			Login:        "old@test.com",
			Password:     "oldpass",
		}
		_, err := storageService.CreateCredentials(ctx, creds)
		require.NoError(t, err)

		updatedCreds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "creds-id",
				Metadata: "Updated account",
			},
			Login:    "new@test.com",
			Password: "newpass",
		}

		serverResponse := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "creds-id",
				Metadata: "Updated account",
			},
			Login:    "new@test.com",
			Password: "newpass",
		}

		mockAPI.On("UpdateCredentials", ctx, updatedCreds).Return(serverResponse, nil)

		result, err := gophkeeperService.UpdateCredentials(ctx, updatedCreds)
		require.NoError(t, err)
		assert.Equal(t, "new@test.com", result.Login)

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_DeleteOperations(t *testing.T) {
	ctx := context.Background()

	t.Run("Delete binary - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		// Создаем запись для удаления
		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "delete-id"},
			Data:         []byte("data"),
		}
		_, err := storageService.CreateBinary(ctx, binary)
		require.NoError(t, err)

		mockAPI.On("DeleteBinary", ctx, "delete-id").Return(nil)

		err = gophkeeperService.DeleteBinary(ctx, "delete-id")
		require.NoError(t, err)

		// Проверяем, что локальная запись удалена
		localData, err := storageService.GetBinary(ctx, "delete-id")
		require.NoError(t, err)
		assert.Nil(t, localData)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Delete binary - API error", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "keep-id"},
			Data:         []byte("data"),
		}
		_, err := storageService.CreateBinary(ctx, binary)
		require.NoError(t, err)

		mockAPI.On("DeleteBinary", ctx, "keep-id").Return(errors.New("server error"))

		err = gophkeeperService.DeleteBinary(ctx, "keep-id")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "server error")

		// Проверяем, что локальная запись НЕ удалена
		localData, err := storageService.GetBinary(ctx, "keep-id")
		require.NoError(t, err)
		assert.NotNil(t, localData)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Delete card - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		card := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{ID: "card-delete"},
			Number:       "4111111111111111",
		}
		_, err := storageService.CreateCard(ctx, card)
		require.NoError(t, err)

		mockAPI.On("DeleteCard", ctx, "card-delete").Return(nil)

		err = gophkeeperService.DeleteCard(ctx, "card-delete")
		require.NoError(t, err)

		localCard, err := storageService.GetCard(ctx, "card-delete")
		require.NoError(t, err)
		assert.Nil(t, localCard)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Delete credentials - successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		creds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{ID: "creds-delete"},
			Login:        "user@test.com",
		}
		_, err := storageService.CreateCredentials(ctx, creds)
		require.NoError(t, err)

		mockAPI.On("DeleteCredentials", ctx, "creds-delete").Return(nil)

		err = gophkeeperService.DeleteCredentials(ctx, "creds-delete")
		require.NoError(t, err)

		localCreds, err := storageService.GetCredentials(ctx, "creds-delete")
		require.NoError(t, err)
		assert.Nil(t, localCreds)

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_ForceSync(t *testing.T) {
	ctx := context.Background()

	t.Run("Force sync successful", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		// Настраиваем успешную синхронизацию
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		err := gophkeeperService.ForceSync(ctx)
		require.NoError(t, err)

		mockAPI.AssertExpectations(t)
	})

	t.Run("Force sync returns sync error", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, errors.New("sync failed"))

		err := gophkeeperService.ForceSync(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sync failed")

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("Update with server success but local failure", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		// Создаем запись
		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "test-id"},
			Data:         []byte("old"),
		}
		_, err := storageService.CreateBinary(ctx, binary)
		require.NoError(t, err)

		// Обновляем с данными, которые вызовут ошибку в локальном хранилище
		// Например, с пустым ID (хотя это маловероятно в реальном сценарии)
		updatedEntity := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "", // Пустой ID вызовет ошибку
				Metadata: "Updated",
			},
			Data: []byte("new"),
		}

		serverResponse := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "", // Сервер тоже возвращает пустой ID
				Metadata: "Updated",
			},
			Data: []byte("new"),
		}

		mockAPI.On("UpdateBinary", ctx, updatedEntity).Return(serverResponse, nil)

		result, err := gophkeeperService.UpdateBinary(ctx, updatedEntity)
		// Ожидаем ошибку при локальном обновлении
		require.Error(t, err)
		assert.NotNil(t, result) // Но серверный ответ возвращается
		assert.Contains(t, err.Error(), "updated on server but local failed")

		mockAPI.AssertExpectations(t)
	})

	t.Run("Delete with server success but local failure", func(t *testing.T) {
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		// Не создаем запись локально, чтобы удаление вызвало ошибку
		mockAPI.On("DeleteBinary", ctx, "non-existent").Return(nil)

		err := gophkeeperService.DeleteBinary(ctx, "non-existent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "deleted on server but local failed")

		mockAPI.AssertExpectations(t)
	})

	t.Run("Create with server success but local storage failure scenario", func(t *testing.T) {
		// Этот тест сложно реализовать с in-memory репозиториями,
		// так как они не генерируют ошибки при создании.
		// В реальном проекте использовались бы моки репозиториев.
		mockAPI := new(MockGophKeeperAPIClient)
		dbManager := inmemory.NewDatabaseManager()
		storageService := services.NewStorageService(
			dbManager.BinariesRepo,
			dbManager.CardsRepo,
			dbManager.CredentialsRepo,
			dbManager.TextsRepo,
		)
		syncService := services.NewSyncService(mockAPI, storageService)
		gophkeeperService := services.NewGophkeeperService(mockAPI, storageService, syncService)

		dto := &dtos.NewBinaryData{
			Data:            []byte("test"),
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Test"},
		}

		serverResponse := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "test-id",
				Metadata: "Test",
			},
			Data: []byte("test"),
		}

		// Создаем дубликат, чтобы вызвать ошибку при локальном создании
		existing := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{ID: "test-id"},
			Data:         []byte("existing"),
		}
		_, err := storageService.CreateBinary(ctx, existing)
		require.NoError(t, err)

		mockAPI.On("CreateBinary", ctx, dto).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateBinary(ctx, dto)
		require.Error(t, err)
		assert.NotNil(t, result) // Серверный ответ возвращается
		assert.Contains(t, err.Error(), "created on server but local failed")
		assert.Contains(t, err.Error(), "already exists")

		mockAPI.AssertExpectations(t)
	})
}
