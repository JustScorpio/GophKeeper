package services_test

import (
	"context"
	"sort"
	"testing"

	"github.com/JustScorpio/GophKeeper/frontend/internal/encryption"
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

func (m *MockGophKeeperAPIClient) GetBinary(ctx context.Context, id string) (*entities.BinaryData, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.BinaryData), args.Error(1)
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

func (m *MockGophKeeperAPIClient) GetCard(ctx context.Context, id string) (*entities.CardInformation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CardInformation), args.Error(1)
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

func (m *MockGophKeeperAPIClient) GetCredentials(ctx context.Context, id string) (*entities.Credentials, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Credentials), args.Error(1)
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

func (m *MockGophKeeperAPIClient) GetText(ctx context.Context, id string) (*entities.TextData, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.TextData), args.Error(1)
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

func TestGophkeeperService_Encryption(t *testing.T) {
	ctx := context.Background()
	testPassword := "testpass123"

	t.Run("SetEncryption and IsEncryptionSet work correctly", func(t *testing.T) {
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

		// Initially encryption is not set
		assert.False(t, gophkeeperService.IsEncryptionSet())

		// Set encryption
		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)
		assert.True(t, gophkeeperService.IsEncryptionSet())
	})

	t.Run("Register sets encryption", func(t *testing.T) {
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

		mockAPI.On("Register", ctx, "testuser", testPassword).Return(nil)

		err := gophkeeperService.Register(ctx, "testuser", testPassword)
		require.NoError(t, err)
		assert.True(t, gophkeeperService.IsEncryptionSet())

		mockAPI.AssertExpectations(t)
	})

	t.Run("Login sets encryption", func(t *testing.T) {
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

		mockAPI.On("Login", ctx, "user", testPassword).Return(nil)
		mockAPI.On("GetAllBinaries", ctx).Return([]entities.BinaryData{}, nil)
		mockAPI.On("GetAllCards", ctx).Return([]entities.CardInformation{}, nil)
		mockAPI.On("GetAllCredentials", ctx).Return([]entities.Credentials{}, nil)
		mockAPI.On("GetAllTexts", ctx).Return([]entities.TextData{}, nil)

		err := gophkeeperService.Login(ctx, "user", testPassword)
		require.NoError(t, err)
		assert.True(t, gophkeeperService.IsEncryptionSet())

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_BinaryCRUDOperations(t *testing.T) {
	ctx := context.Background()
	testPassword := "testpass123"
	cryptoService := encryption.NewCryptoService(testPassword)

	t.Run("CreateBinary - successful with encryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		dto := &dtos.NewBinaryData{
			Data:            []byte("test binary data"),
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Test binary metadata"},
		}

		// Expected encrypted data
		encryptedData, err := cryptoService.EncryptBytes([]byte("test binary data"))
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Test binary metadata")
		require.NoError(t, err)

		serverResponse := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "binary-123",
				Metadata: encryptedMetadata,
			},
			Data: encryptedData,
		}

		mockAPI.On("CreateBinary", ctx, mock.AnythingOfType("*dtos.NewBinaryData")).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateBinary(ctx, dto)
		require.NoError(t, err)
		assert.Equal(t, "binary-123", result.ID)
		assert.Equal(t, []byte("test binary data"), result.Data)
		assert.Equal(t, "Test binary metadata", result.Metadata)

		// Verify data is encrypted in local storage
		localData, err := storageService.GetBinary(ctx, "binary-123")
		require.NoError(t, err)
		assert.Equal(t, encryptedMetadata, localData.Metadata)
		assert.Equal(t, encryptedData, localData.Data)

		mockAPI.AssertExpectations(t)
	})

	t.Run("GetBinary - successful with decryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Store encrypted data directly
		encryptedData, err := cryptoService.EncryptBytes([]byte("secret binary data"))
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("secret metadata")
		require.NoError(t, err)

		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "get-binary-123",
				Metadata: encryptedMetadata,
			},
			Data: encryptedData,
		}

		_, err = storageService.CreateBinary(ctx, binary)
		require.NoError(t, err)

		result, err := gophkeeperService.GetBinary(ctx, "get-binary-123")
		require.NoError(t, err)
		assert.Equal(t, "get-binary-123", result.ID)
		assert.Equal(t, []byte("secret binary data"), result.Data)
		assert.Equal(t, "secret metadata", result.Metadata)
	})

	t.Run("GetAllBinaries - successful with decryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create multiple encrypted binaries
		for i := 1; i <= 3; i++ {
			encryptedData, err := cryptoService.EncryptBytes([]byte("data " + string(rune('A'+i-1))))
			require.NoError(t, err)
			encryptedMetadata, err := cryptoService.Encrypt("metadata " + string(rune('A'+i-1)))
			require.NoError(t, err)

			binary := &entities.BinaryData{
				SecureEntity: entities.SecureEntity{
					ID:       "binary-" + string(rune('0'+i)),
					Metadata: encryptedMetadata,
				},
				Data: encryptedData,
			}

			_, err = storageService.CreateBinary(ctx, binary)
			require.NoError(t, err)
		}

		results, err := gophkeeperService.GetAllBinaries(ctx)
		require.NoError(t, err)
		assert.Len(t, results, 3)

		// Сортируем результаты по ID
		sort.Slice(results, func(i, j int) bool {
			return results[i].ID < results[j].ID
		})

		for i, result := range results {
			assert.Equal(t, "binary-"+string(rune('0'+i+1)), result.ID)
			assert.Equal(t, []byte("data "+string(rune('A'+i))), result.Data)
			assert.Equal(t, "metadata "+string(rune('A'+i)), result.Metadata)
		}
	})

	t.Run("UpdateBinary - successful with encryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create initial encrypted binary
		initialData, err := cryptoService.EncryptBytes([]byte("initial data"))
		require.NoError(t, err)
		initialMetadata, err := cryptoService.Encrypt("initial metadata")
		require.NoError(t, err)

		initialBinary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "update-binary-123",
				Metadata: initialMetadata,
			},
			Data: initialData,
		}

		_, err = storageService.CreateBinary(ctx, initialBinary)
		require.NoError(t, err)

		// Update with new data
		updatedBinary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "update-binary-123",
				Metadata: "updated metadata",
			},
			Data: []byte("updated data"),
		}

		// Server response with encrypted data
		updatedEncryptedData, err := cryptoService.EncryptBytes([]byte("updated data"))
		require.NoError(t, err)
		updatedEncryptedMetadata, err := cryptoService.Encrypt("updated metadata")
		require.NoError(t, err)

		serverResponse := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "update-binary-123",
				Metadata: updatedEncryptedMetadata,
			},
			Data: updatedEncryptedData,
		}

		mockAPI.On("UpdateBinary", ctx, mock.AnythingOfType("*entities.BinaryData")).Return(serverResponse, nil)

		result, err := gophkeeperService.UpdateBinary(ctx, updatedBinary)
		require.NoError(t, err)
		assert.Equal(t, "update-binary-123", result.ID)
		assert.Equal(t, []byte("updated data"), result.Data)
		assert.Equal(t, "updated metadata", result.Metadata)

		// Verify local storage has encrypted data
		localData, err := storageService.GetBinary(ctx, "update-binary-123")
		require.NoError(t, err)
		assert.Equal(t, updatedEncryptedMetadata, localData.Metadata)
		assert.Equal(t, updatedEncryptedData, localData.Data)

		mockAPI.AssertExpectations(t)
	})

	t.Run("DeleteBinary - successful", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create a binary to delete
		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "delete-binary-123",
				Metadata: "metadata",
			},
		}
		_, err = storageService.CreateBinary(ctx, binary)
		require.NoError(t, err)

		mockAPI.On("DeleteBinary", ctx, "delete-binary-123").Return(nil)

		err = gophkeeperService.DeleteBinary(ctx, "delete-binary-123")
		require.NoError(t, err)

		// Verify binary is deleted locally
		localData, err := storageService.GetBinary(ctx, "delete-binary-123")
		require.NoError(t, err)
		assert.Nil(t, localData)

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_CardCRUDOperations(t *testing.T) {
	ctx := context.Background()
	testPassword := "testpass123"
	cryptoService := encryption.NewCryptoService(testPassword)

	t.Run("CreateCard - successful with encryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		dto := &dtos.NewCardInformation{
			Number:          "4111111111111111",
			CardHolder:      "John Doe",
			ExpirationDate:  "12/25",
			CVV:             "123",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Visa Card"},
		}

		// Encrypt all fields
		encryptedNumber, err := cryptoService.Encrypt("4111111111111111")
		require.NoError(t, err)
		encryptedHolder, err := cryptoService.Encrypt("John Doe")
		require.NoError(t, err)
		encryptedExpDate, err := cryptoService.Encrypt("12/25")
		require.NoError(t, err)
		encryptedCVV, err := cryptoService.Encrypt("123")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Visa Card")
		require.NoError(t, err)

		serverResponse := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "card-123",
				Metadata: encryptedMetadata,
			},
			Number:         encryptedNumber,
			CardHolder:     encryptedHolder,
			ExpirationDate: encryptedExpDate,
			CVV:            encryptedCVV,
		}

		mockAPI.On("CreateCard", ctx, mock.AnythingOfType("*dtos.NewCardInformation")).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateCard(ctx, dto)
		require.NoError(t, err)
		assert.Equal(t, "card-123", result.ID)
		assert.Equal(t, "4111111111111111", result.Number)
		assert.Equal(t, "John Doe", result.CardHolder)
		assert.Equal(t, "12/25", result.ExpirationDate)
		assert.Equal(t, "123", result.CVV)
		assert.Equal(t, "Visa Card", result.Metadata)

		mockAPI.AssertExpectations(t)
	})

	t.Run("GetCard - successful with decryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Store encrypted card
		encryptedNumber, err := cryptoService.Encrypt("5555555555554444")
		require.NoError(t, err)
		encryptedHolder, err := cryptoService.Encrypt("Jane Doe")
		require.NoError(t, err)
		encryptedExpDate, err := cryptoService.Encrypt("06/27")
		require.NoError(t, err)
		encryptedCVV, err := cryptoService.Encrypt("789")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("MasterCard")
		require.NoError(t, err)

		card := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "get-card-123",
				Metadata: encryptedMetadata,
			},
			Number:         encryptedNumber,
			CardHolder:     encryptedHolder,
			ExpirationDate: encryptedExpDate,
			CVV:            encryptedCVV,
		}

		_, err = storageService.CreateCard(ctx, card)
		require.NoError(t, err)

		result, err := gophkeeperService.GetCard(ctx, "get-card-123")
		require.NoError(t, err)
		assert.Equal(t, "get-card-123", result.ID)
		assert.Equal(t, "5555555555554444", result.Number)
		assert.Equal(t, "Jane Doe", result.CardHolder)
		assert.Equal(t, "06/27", result.ExpirationDate)
		assert.Equal(t, "789", result.CVV)
		assert.Equal(t, "MasterCard", result.Metadata)
	})

	t.Run("GetAllCards - successful with decryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create multiple encrypted cards
		cardsData := []struct {
			id       string
			number   string
			holder   string
			expDate  string
			cvv      string
			metadata string
		}{
			{"card-1", "4111111111111111", "John A", "01/24", "111", "Visa Personal"},
			{"card-2", "5555555555554444", "Jane B", "02/25", "222", "MasterCard Business"},
			{"card-3", "378282246310005", "Bob C", "03/26", "333", "Amex Corporate"},
		}

		for _, data := range cardsData {
			encryptedNumber, err := cryptoService.Encrypt(data.number)
			require.NoError(t, err)
			encryptedHolder, err := cryptoService.Encrypt(data.holder)
			require.NoError(t, err)
			encryptedExpDate, err := cryptoService.Encrypt(data.expDate)
			require.NoError(t, err)
			encryptedCVV, err := cryptoService.Encrypt(data.cvv)
			require.NoError(t, err)
			encryptedMetadata, err := cryptoService.Encrypt(data.metadata)
			require.NoError(t, err)

			card := &entities.CardInformation{
				SecureEntity: entities.SecureEntity{
					ID:       data.id,
					Metadata: encryptedMetadata,
				},
				Number:         encryptedNumber,
				CardHolder:     encryptedHolder,
				ExpirationDate: encryptedExpDate,
				CVV:            encryptedCVV,
			}

			_, err = storageService.CreateCard(ctx, card)
			require.NoError(t, err)
		}

		results, err := gophkeeperService.GetAllCards(ctx)
		require.NoError(t, err)
		assert.Len(t, results, 3)

		// Сортируем результаты по ID
		sort.Slice(results, func(i, j int) bool {
			return results[i].ID < results[j].ID
		})

		for i, result := range results {
			data := cardsData[i]
			assert.Equal(t, data.id, result.ID)
			assert.Equal(t, data.number, result.Number)
			assert.Equal(t, data.holder, result.CardHolder)
			assert.Equal(t, data.expDate, result.ExpirationDate)
			assert.Equal(t, data.cvv, result.CVV)
			assert.Equal(t, data.metadata, result.Metadata)
		}
	})

	t.Run("UpdateCard - successful with encryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create initial encrypted card
		encryptedNumber, err := cryptoService.Encrypt("4111111111111111")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Old Visa")
		require.NoError(t, err)

		initialCard := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "update-card-123",
				Metadata: encryptedMetadata,
			},
			Number: encryptedNumber,
		}

		_, err = storageService.CreateCard(ctx, initialCard)
		require.NoError(t, err)

		// Update card
		updatedCard := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "update-card-123",
				Metadata: "Updated Visa Premium",
			},
			Number:         "5555666677778888",
			CardHolder:     "Updated Holder",
			ExpirationDate: "12/28",
			CVV:            "999",
		}

		// Server response with encrypted data
		updatedEncryptedNumber, err := cryptoService.Encrypt("5555666677778888")
		require.NoError(t, err)
		updatedEncryptedHolder, err := cryptoService.Encrypt("Updated Holder")
		require.NoError(t, err)
		updatedEncryptedExpDate, err := cryptoService.Encrypt("12/28")
		require.NoError(t, err)
		updatedEncryptedCVV, err := cryptoService.Encrypt("999")
		require.NoError(t, err)
		updatedEncryptedMetadata, err := cryptoService.Encrypt("Updated Visa Premium")
		require.NoError(t, err)

		serverResponse := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "update-card-123",
				Metadata: updatedEncryptedMetadata,
			},
			Number:         updatedEncryptedNumber,
			CardHolder:     updatedEncryptedHolder,
			ExpirationDate: updatedEncryptedExpDate,
			CVV:            updatedEncryptedCVV,
		}

		mockAPI.On("UpdateCard", ctx, mock.AnythingOfType("*entities.CardInformation")).Return(serverResponse, nil)

		result, err := gophkeeperService.UpdateCard(ctx, updatedCard)
		require.NoError(t, err)
		assert.Equal(t, "update-card-123", result.ID)
		assert.Equal(t, "5555666677778888", result.Number)
		assert.Equal(t, "Updated Holder", result.CardHolder)
		assert.Equal(t, "12/28", result.ExpirationDate)
		assert.Equal(t, "999", result.CVV)
		assert.Equal(t, "Updated Visa Premium", result.Metadata)

		mockAPI.AssertExpectations(t)
	})

	t.Run("DeleteCard - successful", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		card := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "delete-card-123",
				Metadata: "Card to delete",
			},
		}
		_, err = storageService.CreateCard(ctx, card)
		require.NoError(t, err)

		mockAPI.On("DeleteCard", ctx, "delete-card-123").Return(nil)

		err = gophkeeperService.DeleteCard(ctx, "delete-card-123")
		require.NoError(t, err)

		localCard, err := storageService.GetCard(ctx, "delete-card-123")
		require.NoError(t, err)
		assert.Nil(t, localCard)

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_CredentialsCRUDOperations(t *testing.T) {
	ctx := context.Background()
	testPassword := "testpass123"
	cryptoService := encryption.NewCryptoService(testPassword)

	t.Run("CreateCredentials - successful with encryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		dto := &dtos.NewCredentials{
			Login:           "user@example.com",
			Password:        "strongpassword123",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Work Email"},
		}

		// Encrypt fields
		encryptedLogin, err := cryptoService.Encrypt("user@example.com")
		require.NoError(t, err)
		encryptedPassword, err := cryptoService.Encrypt("strongpassword123")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Work Email")
		require.NoError(t, err)

		serverResponse := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "creds-123",
				Metadata: encryptedMetadata,
			},
			Login:    encryptedLogin,
			Password: encryptedPassword,
		}

		mockAPI.On("CreateCredentials", ctx, mock.AnythingOfType("*dtos.NewCredentials")).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateCredentials(ctx, dto)
		require.NoError(t, err)
		assert.Equal(t, "creds-123", result.ID)
		assert.Equal(t, "user@example.com", result.Login)
		assert.Equal(t, "strongpassword123", result.Password)
		assert.Equal(t, "Work Email", result.Metadata)

		mockAPI.AssertExpectations(t)
	})

	t.Run("GetCredentials - successful with decryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Store encrypted credentials
		encryptedLogin, err := cryptoService.Encrypt("admin@company.com")
		require.NoError(t, err)
		encryptedPassword, err := cryptoService.Encrypt("AdminPass123!")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Admin Account")
		require.NoError(t, err)

		creds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "get-creds-123",
				Metadata: encryptedMetadata,
			},
			Login:    encryptedLogin,
			Password: encryptedPassword,
		}

		_, err = storageService.CreateCredentials(ctx, creds)
		require.NoError(t, err)

		result, err := gophkeeperService.GetCredentials(ctx, "get-creds-123")
		require.NoError(t, err)
		assert.Equal(t, "get-creds-123", result.ID)
		assert.Equal(t, "admin@company.com", result.Login)
		assert.Equal(t, "AdminPass123!", result.Password)
		assert.Equal(t, "Admin Account", result.Metadata)
	})

	t.Run("GetAllCredentials - successful with decryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create multiple encrypted credentials
		credsData := []struct {
			id       string
			login    string
			password string
			metadata string
		}{
			{"creds-1", "user1@test.com", "pass1", "Personal Email"},
			{"creds-2", "user2@work.com", "pass2", "Work Account"},
			{"creds-3", "admin@system.com", "admin123", "System Admin"},
		}

		for _, data := range credsData {
			encryptedLogin, err := cryptoService.Encrypt(data.login)
			require.NoError(t, err)
			encryptedPassword, err := cryptoService.Encrypt(data.password)
			require.NoError(t, err)
			encryptedMetadata, err := cryptoService.Encrypt(data.metadata)
			require.NoError(t, err)

			creds := &entities.Credentials{
				SecureEntity: entities.SecureEntity{
					ID:       data.id,
					Metadata: encryptedMetadata,
				},
				Login:    encryptedLogin,
				Password: encryptedPassword,
			}

			_, err = storageService.CreateCredentials(ctx, creds)
			require.NoError(t, err)
		}

		results, err := gophkeeperService.GetAllCredentials(ctx)
		require.NoError(t, err)
		assert.Len(t, results, 3)

		// Сортируем результаты по ID
		sort.Slice(results, func(i, j int) bool {
			return results[i].ID < results[j].ID
		})

		for i, result := range results {
			data := credsData[i]
			assert.Equal(t, data.id, result.ID)
			assert.Equal(t, data.login, result.Login)
			assert.Equal(t, data.password, result.Password)
			assert.Equal(t, data.metadata, result.Metadata)
		}
	})

	t.Run("UpdateCredentials - successful with encryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create initial encrypted credentials
		encryptedLogin, err := cryptoService.Encrypt("old@test.com")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Old Account")
		require.NoError(t, err)

		initialCreds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "update-creds-123",
				Metadata: encryptedMetadata,
			},
			Login: encryptedLogin,
		}

		_, err = storageService.CreateCredentials(ctx, initialCreds)
		require.NoError(t, err)

		// Update credentials
		updatedCreds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "update-creds-123",
				Metadata: "Updated Account",
			},
			Login:    "new@test.com",
			Password: "newpassword456",
		}

		// Server response with encrypted data
		updatedEncryptedLogin, err := cryptoService.Encrypt("new@test.com")
		require.NoError(t, err)
		updatedEncryptedPassword, err := cryptoService.Encrypt("newpassword456")
		require.NoError(t, err)
		updatedEncryptedMetadata, err := cryptoService.Encrypt("Updated Account")
		require.NoError(t, err)

		serverResponse := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "update-creds-123",
				Metadata: updatedEncryptedMetadata,
			},
			Login:    updatedEncryptedLogin,
			Password: updatedEncryptedPassword,
		}

		mockAPI.On("UpdateCredentials", ctx, mock.AnythingOfType("*entities.Credentials")).Return(serverResponse, nil)

		result, err := gophkeeperService.UpdateCredentials(ctx, updatedCreds)
		require.NoError(t, err)
		assert.Equal(t, "update-creds-123", result.ID)
		assert.Equal(t, "new@test.com", result.Login)
		assert.Equal(t, "newpassword456", result.Password)
		assert.Equal(t, "Updated Account", result.Metadata)

		mockAPI.AssertExpectations(t)
	})

	t.Run("DeleteCredentials - successful", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		creds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "delete-creds-123",
				Metadata: "Account to delete",
			},
		}
		_, err = storageService.CreateCredentials(ctx, creds)
		require.NoError(t, err)

		mockAPI.On("DeleteCredentials", ctx, "delete-creds-123").Return(nil)

		err = gophkeeperService.DeleteCredentials(ctx, "delete-creds-123")
		require.NoError(t, err)

		localCreds, err := storageService.GetCredentials(ctx, "delete-creds-123")
		require.NoError(t, err)
		assert.Nil(t, localCreds)

		mockAPI.AssertExpectations(t)
	})
}

func TestGophkeeperService_TextCRUDOperations(t *testing.T) {
	ctx := context.Background()
	testPassword := "testpass123"
	cryptoService := encryption.NewCryptoService(testPassword)

	t.Run("CreateText - successful with encryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		dto := &dtos.NewTextData{
			Data:            "This is a secret note that should be encrypted.",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "Secret Note"},
		}

		// Encrypt fields
		encryptedData, err := cryptoService.Encrypt("This is a secret note that should be encrypted.")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Secret Note")
		require.NoError(t, err)

		serverResponse := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       "text-123",
				Metadata: encryptedMetadata,
			},
			Data: encryptedData,
		}

		mockAPI.On("CreateText", ctx, mock.AnythingOfType("*dtos.NewTextData")).Return(serverResponse, nil)

		result, err := gophkeeperService.CreateText(ctx, dto)
		require.NoError(t, err)
		assert.Equal(t, "text-123", result.ID)
		assert.Equal(t, "This is a secret note that should be encrypted.", result.Data)
		assert.Equal(t, "Secret Note", result.Metadata)

		mockAPI.AssertExpectations(t)
	})

	t.Run("GetText - successful with decryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Store encrypted text
		encryptedData, err := cryptoService.Encrypt("This is very sensitive information that must be protected.")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Sensitive Info")
		require.NoError(t, err)

		text := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       "get-text-123",
				Metadata: encryptedMetadata,
			},
			Data: encryptedData,
		}

		_, err = storageService.CreateText(ctx, text)
		require.NoError(t, err)

		result, err := gophkeeperService.GetText(ctx, "get-text-123")
		require.NoError(t, err)
		assert.Equal(t, "get-text-123", result.ID)
		assert.Equal(t, "This is very sensitive information that must be protected.", result.Data)
		assert.Equal(t, "Sensitive Info", result.Metadata)
	})

	t.Run("GetAllTexts - successful with decryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create multiple encrypted texts
		textsData := []struct {
			id       string
			data     string
			metadata string
		}{
			{"text-1", "First secret note", "Note 1"},
			{"text-2", "Second secret note with more details", "Note 2"},
			{"text-3", "Third note containing important information", "Note 3"},
		}

		for _, data := range textsData {
			encryptedData, err := cryptoService.Encrypt(data.data)
			require.NoError(t, err)
			encryptedMetadata, err := cryptoService.Encrypt(data.metadata)
			require.NoError(t, err)

			text := &entities.TextData{
				SecureEntity: entities.SecureEntity{
					ID:       data.id,
					Metadata: encryptedMetadata,
				},
				Data: encryptedData,
			}

			_, err = storageService.CreateText(ctx, text)
			require.NoError(t, err)
		}

		results, err := gophkeeperService.GetAllTexts(ctx)
		require.NoError(t, err)
		assert.Len(t, results, 3)

		// Сортируем результаты по ID
		sort.Slice(results, func(i, j int) bool {
			return results[i].ID < results[j].ID
		})

		for i, result := range results {
			data := textsData[i]
			assert.Equal(t, data.id, result.ID)
			assert.Equal(t, data.data, result.Data)
			assert.Equal(t, data.metadata, result.Metadata)
		}
	})

	t.Run("UpdateText - successful with encryption", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		// Create initial encrypted text
		encryptedData, err := cryptoService.Encrypt("Old text content")
		require.NoError(t, err)
		encryptedMetadata, err := cryptoService.Encrypt("Old Note")
		require.NoError(t, err)

		initialText := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       "update-text-123",
				Metadata: encryptedMetadata,
			},
			Data: encryptedData,
		}

		_, err = storageService.CreateText(ctx, initialText)
		require.NoError(t, err)

		// Update text
		updatedText := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       "update-text-123",
				Metadata: "Updated Important Note",
			},
			Data: "This is the updated text content with important changes that need to be secured.",
		}

		// Server response with encrypted data
		updatedEncryptedData, err := cryptoService.Encrypt("This is the updated text content with important changes that need to be secured.")
		require.NoError(t, err)
		updatedEncryptedMetadata, err := cryptoService.Encrypt("Updated Important Note")
		require.NoError(t, err)

		serverResponse := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       "update-text-123",
				Metadata: updatedEncryptedMetadata,
			},
			Data: updatedEncryptedData,
		}

		mockAPI.On("UpdateText", ctx, mock.AnythingOfType("*entities.TextData")).Return(serverResponse, nil)

		result, err := gophkeeperService.UpdateText(ctx, updatedText)
		require.NoError(t, err)
		assert.Equal(t, "update-text-123", result.ID)
		assert.Equal(t, "This is the updated text content with important changes that need to be secured.", result.Data)
		assert.Equal(t, "Updated Important Note", result.Metadata)

		mockAPI.AssertExpectations(t)
	})

	t.Run("DeleteText - successful", func(t *testing.T) {
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

		err := gophkeeperService.SetEncryption(testPassword)
		require.NoError(t, err)

		text := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       "delete-text-123",
				Metadata: "Text to delete",
			},
		}
		_, err = storageService.CreateText(ctx, text)
		require.NoError(t, err)

		mockAPI.On("DeleteText", ctx, "delete-text-123").Return(nil)

		err = gophkeeperService.DeleteText(ctx, "delete-text-123")
		require.NoError(t, err)

		localText, err := storageService.GetText(ctx, "delete-text-123")
		require.NoError(t, err)
		assert.Nil(t, localText)

		mockAPI.AssertExpectations(t)
	})
}
