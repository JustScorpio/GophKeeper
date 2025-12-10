package services_test

import (
	"context"
	"testing"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/frontend/internal/repositories/inmemory"
	"github.com/JustScorpio/GophKeeper/frontend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStorageService - Тесты сервиса хранилища
func TestStorageService(t *testing.T) {
	ctx := context.Background()

	// Создаем in-memory репозитории
	dbManager := inmemory.NewDatabaseManager()

	// Создаем сервис
	storageService := services.NewStorageService(
		dbManager.BinariesRepo,
		dbManager.CardsRepo,
		dbManager.CredentialsRepo,
		dbManager.TextsRepo,
	)

	t.Run("Binary CRUD operations", func(t *testing.T) {
		// Create
		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "test-binary-1",
				Metadata: "Test binary data",
			},
			Data: []byte("test data"),
		}

		created, err := storageService.CreateBinary(ctx, binary)
		require.NoError(t, err)
		assert.Equal(t, binary.ID, created.ID)
		assert.Equal(t, binary.Data, created.Data)

		// Get
		retrieved, err := storageService.GetBinary(ctx, binary.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, binary.ID, retrieved.ID)

		// GetAll
		allBinaries, err := storageService.GetAllBinaries(ctx)
		require.NoError(t, err)
		assert.Len(t, allBinaries, 1)

		// Update
		updatedBinary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       binary.ID,
				Metadata: "Updated metadata",
			},
			Data: []byte("updated data"),
		}

		updated, err := storageService.UpdateBinary(ctx, updatedBinary)
		require.NoError(t, err)
		assert.Equal(t, "Updated metadata", updated.Metadata)
		assert.Equal(t, []byte("updated data"), updated.Data)

		// Delete
		err = storageService.DeleteBinary(ctx, binary.ID)
		require.NoError(t, err)

		deleted, err := storageService.GetBinary(ctx, binary.ID)
		require.NoError(t, err)
		assert.Nil(t, deleted)
	})

	t.Run("Card CRUD operations", func(t *testing.T) {
		// Create
		card := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       "test-card-1",
				Metadata: "Personal card",
			},
			Number:         "4111111111111111",
			CardHolder:     "John Doe",
			ExpirationDate: "12/25",
			CVV:            "123",
		}

		created, err := storageService.CreateCard(ctx, card)
		require.NoError(t, err)
		assert.Equal(t, card.ID, created.ID)
		assert.Equal(t, card.Number, created.Number)

		// Get
		retrieved, err := storageService.GetCard(ctx, card.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, card.CardHolder, retrieved.CardHolder)

		// GetAll
		allCards, err := storageService.GetAllCards(ctx)
		require.NoError(t, err)
		assert.Len(t, allCards, 1)

		// Update
		updatedCard := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       card.ID,
				Metadata: "Updated card",
			},
			Number:         "5555555555554444",
			CardHolder:     "Jane Smith",
			ExpirationDate: "11/26",
			CVV:            "456",
		}

		updated, err := storageService.UpdateCard(ctx, updatedCard)
		require.NoError(t, err)
		assert.Equal(t, "Jane Smith", updated.CardHolder)

		// Delete
		err = storageService.DeleteCard(ctx, card.ID)
		require.NoError(t, err)

		deleted, err := storageService.GetCard(ctx, card.ID)
		require.NoError(t, err)
		assert.Nil(t, deleted)
	})

	t.Run("Credentials CRUD operations", func(t *testing.T) {
		// Create
		creds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       "test-creds-1",
				Metadata: "Gmail account",
			},
			Login:    "user@gmail.com",
			Password: "password123",
		}

		created, err := storageService.CreateCredentials(ctx, creds)
		require.NoError(t, err)
		assert.Equal(t, creds.ID, created.ID)
		assert.Equal(t, creds.Login, created.Login)

		// Get
		retrieved, err := storageService.GetCredentials(ctx, creds.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, creds.Password, retrieved.Password)

		// GetAll
		allCreds, err := storageService.GetAllCredentials(ctx)
		require.NoError(t, err)
		assert.Len(t, allCreds, 1)

		// Update
		updatedCreds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       creds.ID,
				Metadata: "Updated Gmail",
			},
			Login:    "updated@gmail.com",
			Password: "newpassword456",
		}

		updated, err := storageService.UpdateCredentials(ctx, updatedCreds)
		require.NoError(t, err)
		assert.Equal(t, "updated@gmail.com", updated.Login)

		// Delete
		err = storageService.DeleteCredentials(ctx, creds.ID)
		require.NoError(t, err)

		deleted, err := storageService.GetCredentials(ctx, creds.ID)
		require.NoError(t, err)
		assert.Nil(t, deleted)
	})

	t.Run("Text CRUD operations", func(t *testing.T) {
		// Create
		text := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       "test-text-1",
				Metadata: "Important notes",
			},
			Data: "This is some important text data",
		}

		created, err := storageService.CreateText(ctx, text)
		require.NoError(t, err)
		assert.Equal(t, text.ID, created.ID)
		assert.Equal(t, text.Data, created.Data)

		// Get
		retrieved, err := storageService.GetText(ctx, text.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, text.Data, retrieved.Data)

		// GetAll
		allTexts, err := storageService.GetAllTexts(ctx)
		require.NoError(t, err)
		assert.Len(t, allTexts, 1)

		// Update
		updatedText := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       text.ID,
				Metadata: "Updated notes",
			},
			Data: "This is updated text data",
		}

		updated, err := storageService.UpdateText(ctx, updatedText)
		require.NoError(t, err)
		assert.Equal(t, "This is updated text data", updated.Data)

		// Delete
		err = storageService.DeleteText(ctx, text.ID)
		require.NoError(t, err)

		deleted, err := storageService.GetText(ctx, text.ID)
		require.NoError(t, err)
		assert.Nil(t, deleted)
	})

	t.Run("Create without ID should fail", func(t *testing.T) {
		// Test with binary
		binary := &entities.BinaryData{
			Data: []byte("test"),
		}

		_, err := storageService.CreateBinary(ctx, binary)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID cannot be empty")

		// Test with card
		card := &entities.CardInformation{
			Number:         "4111111111111111",
			CardHolder:     "John Doe",
			ExpirationDate: "12/25",
			CVV:            "123",
		}

		_, err = storageService.CreateCard(ctx, card)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID cannot be empty")
	})

	t.Run("Create duplicate ID should fail", func(t *testing.T) {
		// Create first entity
		binary1 := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "duplicate-id",
				Metadata: "First",
			},
			Data: []byte("first"),
		}

		_, err := storageService.CreateBinary(ctx, binary1)
		require.NoError(t, err)

		// Try to create second with same ID
		binary2 := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "duplicate-id",
				Metadata: "Second",
			},
			Data: []byte("second"),
		}

		_, err = storageService.CreateBinary(ctx, binary2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Get non-existent entity returns nil", func(t *testing.T) {
		// Test all entity types
		binary, err := storageService.GetBinary(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, binary)

		card, err := storageService.GetCard(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, card)

		creds, err := storageService.GetCredentials(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, creds)

		text, err := storageService.GetText(ctx, "non-existent")
		require.NoError(t, err)
		assert.Nil(t, text)
	})

	t.Run("Update non-existent entity should fail", func(t *testing.T) {
		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       "non-existent",
				Metadata: "test",
			},
			Data: []byte("test"),
		}

		_, err := storageService.UpdateBinary(ctx, binary)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Delete non-existent entity should fail", func(t *testing.T) {
		err := storageService.DeleteBinary(ctx, "non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
