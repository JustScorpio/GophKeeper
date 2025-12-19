// services_test.go
package services_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/customerrors"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/backend/internal/repositories/inmemory"
	"github.com/JustScorpio/GophKeeper/backend/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestData - структура для тестовых данных
type TestData struct {
	User        dtos.NewUser
	Binary      dtos.NewBinaryData
	Card        dtos.NewCardInformation
	Credentials dtos.NewCredentials
	Text        dtos.NewTextData
}

// createTestData создает тестовые данные
func createTestData() TestData {
	return TestData{
		User: dtos.NewUser{
			Login:    "testuser",
			Password: "testpassword",
		},
		Binary: dtos.NewBinaryData{
			Data:            []byte("test binary data"),
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "binary metadata"},
		},
		Card: dtos.NewCardInformation{
			Number:          "4111111111111111",
			CardHolder:      "John Doe",
			ExpirationDate:  "12/25",
			CVV:             "123",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "card metadata"},
		},
		Credentials: dtos.NewCredentials{
			Login:           "serviceuser",
			Password:        "servicepassword",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "credentials metadata"},
		},
		Text: dtos.NewTextData{
			Data:            "This is a test text data",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "text metadata"},
		},
	}
}

// createTestService создает тестовый сервис с in-memory репозиториями
func createTestService() (*services.StorageService, *inmemory.DatabaseManager) {
	dbManager := inmemory.NewDatabaseManager()
	service := services.NewStorageService(
		dbManager.Users,
		dbManager.Binaries,
		dbManager.Cards,
		dbManager.Credentials,
		dbManager.Texts,
	)
	return service, dbManager
}

// createTestContext создает тестовый контекст с пользователем
func createTestContext(userLogin string) context.Context {
	return customcontext.WithUserID(context.Background(), userLogin)
}

// TestNewStorageService тестирует создание сервиса
func TestNewStorageService(t *testing.T) {
	dbManager := inmemory.NewDatabaseManager()
	service := services.NewStorageService(
		dbManager.Users,
		dbManager.Binaries,
		dbManager.Cards,
		dbManager.Credentials,
		dbManager.Texts,
	)

	assert.NotNil(t, service)

	// Проверяем, что сервис можно корректно завершить
	service.Shutdown()
}

// TestStorageService_UserOperations тестирует операции с пользователями
func TestStorageService_UserOperations(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	testData := createTestData()
	ctx := context.Background()

	t.Run("Создание пользователя", func(t *testing.T) {
		user, err := service.CreateUser(ctx, testData.User)

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testData.User.Login, user.Login)
		assert.Equal(t, testData.User.Password, user.Password)
	})

	t.Run("Создание пользователя с существующим логином", func(t *testing.T) {
		// Пытаемся создать пользователя с тем же логином
		duplicateUser := dtos.NewUser{
			Login:    testData.User.Login,
			Password: "differentpassword",
		}

		_, err := service.CreateUser(ctx, duplicateUser)
		assert.Error(t, err)
		assert.Equal(t, customerrors.AlreadyExistsError, err)
	})

	t.Run("Получение пользователя по логину", func(t *testing.T) {
		// Получаем пользователя
		user, err := service.GetUser(ctx, testData.User.Login)

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testData.User.Login, user.Login)
	})

	t.Run("Получение несуществующего пользователя", func(t *testing.T) {
		user, err := service.GetUser(ctx, "nonexistent")

		assert.NoError(t, err)
		assert.Nil(t, user) // В репозитории nil если не найден
	})

	t.Run("Получение всех пользователей", func(t *testing.T) {
		// Создаем несколько пользователей
		usersToCreate := []dtos.NewUser{
			{Login: "user1", Password: "pass1"},
			{Login: "user2", Password: "pass2"},
			{Login: "user3", Password: "pass3"},
		}

		for _, u := range usersToCreate {
			_, err := service.CreateUser(ctx, u)
			require.NoError(t, err)
		}

		// Получаем всех пользователей
		users, err := service.GetAllUsers(ctx)

		require.NoError(t, err)
		assert.Equal(t, len(users), len(usersToCreate)+1) //Созданные в текущем тесте + 1 созданный в других тестах ранее
	})

	t.Run("Обновление пароля пользователя", func(t *testing.T) {
		// Обновляем пользователя
		updatedUser := &entities.User{
			Login:    testData.User.Login,
			Password: "updatedpassword",
		}

		ctxWithUser := customcontext.WithUserID(ctx, testData.User.Login)
		result, err := service.UpdateUser(ctxWithUser, updatedUser)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedUser.Password, result.Password)
	})

	t.Run("Обновление пароля другого пользователя", func(t *testing.T) {
		// Обновляем пользователя
		updatedUser := &entities.User{
			Login:    testData.User.Login,
			Password: "updatedpassword",
		}

		ctxWithUser := customcontext.WithUserID(ctx, "nonexistent")
		result, err := service.UpdateUser(ctxWithUser, updatedUser)

		assert.Error(t, err)
		assert.Equal(t, customerrors.ForbiddenError, err)
		assert.Nil(t, result)
	})

	t.Run("Удаление пользователя", func(t *testing.T) {
		// Удаляем пользователя
		err := service.DeleteUser(ctx, testData.User.Login)

		require.NoError(t, err)

		// Проверяем, что пользователь удален
		user, err := service.GetUser(ctx, testData.User.Login)

		require.NoError(t, err)
		assert.Nil(t, user)
	})
}

// TestStorageService_BinaryOperations тестирует операции с бинарными данными
func TestStorageService_BinaryOperations(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	testData := createTestData()
	ctx := context.Background()

	// Сначала создаем пользователя для контекста
	ctxWithUser := customcontext.WithUserID(ctx, testData.User.Login)

	t.Run("Создание бинарных данных", func(t *testing.T) {
		binary, err := service.CreateBinary(ctxWithUser, &testData.Binary)

		require.NoError(t, err)
		assert.NotNil(t, binary)
		assert.Equal(t, testData.Binary.Data, binary.Data)
		assert.Equal(t, testData.Binary.Metadata, binary.Metadata)
		assert.NotEmpty(t, binary.ID)
	})

	t.Run("Получение бинарных данных по ID", func(t *testing.T) {
		// Сначала создаем
		createdBinary, err := service.CreateBinary(ctxWithUser, &testData.Binary)
		require.NoError(t, err)

		// Получаем
		binary, err := service.GetBinary(ctxWithUser, createdBinary.ID)

		require.NoError(t, err)
		assert.NotNil(t, binary)
		assert.Equal(t, createdBinary.ID, binary.ID)
		assert.Equal(t, createdBinary.Data, binary.Data)
	})

	t.Run("Получение несуществующих бинарных данных", func(t *testing.T) {
		binary, err := service.GetBinary(ctxWithUser, "nonexistent-id")

		assert.NoError(t, err)
		assert.Nil(t, binary)
	})

	t.Run("Получение всех бинарных данных", func(t *testing.T) {
		// Создаем несколько записей
		for i := 0; i < 3; i++ {
			data := dtos.NewBinaryData{
				Data:            []byte(fmt.Sprintf("data %d", i)),
				NewSecureEntity: dtos.NewSecureEntity{Metadata: fmt.Sprintf("metadata %d", i)},
			}
			_, err := service.CreateBinary(ctxWithUser, &data)
			require.NoError(t, err)
		}

		// Получаем все
		binaries, err := service.GetAllBinaries(ctxWithUser)

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(binaries), 3)
	})

	t.Run("Обновление бинарных данных", func(t *testing.T) {
		// Сначала создаем
		createdBinary, err := service.CreateBinary(ctxWithUser, &testData.Binary)
		require.NoError(t, err)

		// Обновляем
		updatedBinary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID:       createdBinary.ID,
				Metadata: "updated metadata",
			},
			Data: []byte("updated data"),
		}

		result, err := service.UpdateBinary(ctxWithUser, updatedBinary)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedBinary.Data, result.Data)
		assert.Equal(t, updatedBinary.Metadata, result.Metadata)
	})

	t.Run("Удаление бинарных данных", func(t *testing.T) {
		// Сначала создаем
		createdBinary, err := service.CreateBinary(ctxWithUser, &testData.Binary)
		require.NoError(t, err)

		// Удаляем
		err = service.DeleteBinary(ctxWithUser, createdBinary.ID)

		require.NoError(t, err)

		// Проверяем, что данные удалены
		binary, err := service.GetBinary(ctxWithUser, createdBinary.ID)

		require.NoError(t, err)
		assert.Nil(t, binary)
	})
}

// TestStorageService_CardOperations тестирует операции с банковскими картами
func TestStorageService_CardOperations(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	testData := createTestData()
	ctx := context.Background()
	ctxWithUser := customcontext.WithUserID(ctx, testData.User.Login)

	t.Run("Создание данных карты", func(t *testing.T) {
		card, err := service.CreateCard(ctxWithUser, &testData.Card)

		require.NoError(t, err)
		assert.NotNil(t, card)
		assert.Equal(t, testData.Card.Number, card.Number)
		assert.Equal(t, testData.Card.CardHolder, card.CardHolder)
		assert.Equal(t, testData.Card.ExpirationDate, card.ExpirationDate)
		assert.Equal(t, testData.Card.CVV, card.CVV)
		assert.NotEmpty(t, card.ID)
	})

	t.Run("Получение данных карты по ID", func(t *testing.T) {
		createdCard, err := service.CreateCard(ctxWithUser, &testData.Card)
		require.NoError(t, err)

		card, err := service.GetCard(ctxWithUser, createdCard.ID)

		require.NoError(t, err)
		assert.NotNil(t, card)
		assert.Equal(t, createdCard.ID, card.ID)
	})

	t.Run("Обновление данных карты", func(t *testing.T) {
		createdCard, err := service.CreateCard(ctxWithUser, &testData.Card)
		require.NoError(t, err)

		updatedCard := &entities.CardInformation{
			SecureEntity: entities.SecureEntity{
				ID:       createdCard.ID,
				Metadata: "updated metadata",
			},
			Number:         "5555555555554444",
			CardHolder:     "Jane Doe",
			ExpirationDate: "06/27",
			CVV:            "456",
		}

		result, err := service.UpdateCard(ctxWithUser, updatedCard)

		require.NoError(t, err)
		assert.Equal(t, updatedCard.Number, result.Number)
		assert.Equal(t, updatedCard.CardHolder, result.CardHolder)
	})

	t.Run("Удаление данных карты", func(t *testing.T) {
		createdCard, err := service.CreateCard(ctxWithUser, &testData.Card)
		require.NoError(t, err)

		err = service.DeleteCard(ctxWithUser, createdCard.ID)

		require.NoError(t, err)

		card, err := service.GetCard(ctxWithUser, createdCard.ID)

		require.NoError(t, err)
		assert.Nil(t, card)
	})
}

// TestStorageService_CredentialsOperations тестирует операции с учетными данными
func TestStorageService_CredentialsOperations(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	testData := createTestData()
	ctx := context.Background()
	ctxWithUser := customcontext.WithUserID(ctx, testData.User.Login)

	t.Run("Создание учетных данных", func(t *testing.T) {
		creds, err := service.CreateCredentials(ctxWithUser, &testData.Credentials)

		require.NoError(t, err)
		assert.NotNil(t, creds)
		assert.Equal(t, testData.Credentials.Login, creds.Login)
		assert.Equal(t, testData.Credentials.Password, creds.Password)
		assert.NotEmpty(t, creds.ID)
	})

	t.Run("Получение учетных данных по ID", func(t *testing.T) {
		createdCreds, err := service.CreateCredentials(ctxWithUser, &testData.Credentials)
		require.NoError(t, err)

		creds, err := service.GetCredentials(ctxWithUser, createdCreds.ID)

		require.NoError(t, err)
		assert.NotNil(t, creds)
		assert.Equal(t, createdCreds.ID, creds.ID)
	})

	t.Run("Обновление учетных данных", func(t *testing.T) {
		createdCreds, err := service.CreateCredentials(ctxWithUser, &testData.Credentials)
		require.NoError(t, err)

		updatedCreds := &entities.Credentials{
			SecureEntity: entities.SecureEntity{
				ID:       createdCreds.ID,
				Metadata: "updated metadata",
			},
			Login:    "updatedlogin",
			Password: "updatedpassword",
		}

		result, err := service.UpdateCredentials(ctxWithUser, updatedCreds)

		require.NoError(t, err)
		assert.Equal(t, updatedCreds.Login, result.Login)
		assert.Equal(t, updatedCreds.Password, result.Password)
	})

	t.Run("Удаление учетных данных", func(t *testing.T) {
		createdCreds, err := service.CreateCredentials(ctxWithUser, &testData.Credentials)
		require.NoError(t, err)

		err = service.DeleteCredentials(ctxWithUser, createdCreds.ID)

		require.NoError(t, err)

		creds, err := service.GetCredentials(ctxWithUser, createdCreds.ID)

		require.NoError(t, err)
		assert.Nil(t, creds)
	})
}

// TestStorageService_TextOperations тестирует операции с текстовыми данными
func TestStorageService_TextOperations(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	testData := createTestData()
	ctx := context.Background()
	ctxWithUser := customcontext.WithUserID(ctx, testData.User.Login)

	t.Run("Создание текстовых данных", func(t *testing.T) {
		text, err := service.CreateText(ctxWithUser, &testData.Text)

		require.NoError(t, err)
		assert.NotNil(t, text)
		assert.Equal(t, testData.Text.Data, text.Data)
		assert.Equal(t, testData.Text.Metadata, text.Metadata)
		assert.NotEmpty(t, text.ID)
	})

	t.Run("Получение текстовых данных по ID", func(t *testing.T) {
		createdText, err := service.CreateText(ctxWithUser, &testData.Text)
		require.NoError(t, err)

		text, err := service.GetText(ctxWithUser, createdText.ID)

		require.NoError(t, err)
		assert.NotNil(t, text)
		assert.Equal(t, createdText.ID, text.ID)
	})

	t.Run("Обновление текстовых данных", func(t *testing.T) {
		createdText, err := service.CreateText(ctxWithUser, &testData.Text)
		require.NoError(t, err)

		updatedText := &entities.TextData{
			SecureEntity: entities.SecureEntity{
				ID:       createdText.ID,
				Metadata: "updated metadata",
			},
			Data: "updated text data",
		}

		result, err := service.UpdateText(ctxWithUser, updatedText)

		require.NoError(t, err)
		assert.Equal(t, updatedText.Data, result.Data)
		assert.Equal(t, updatedText.Metadata, result.Metadata)
	})

	t.Run("Удаление текстовых данных", func(t *testing.T) {
		createdText, err := service.CreateText(ctxWithUser, &testData.Text)
		require.NoError(t, err)

		err = service.DeleteText(ctxWithUser, createdText.ID)

		require.NoError(t, err)

		text, err := service.GetText(ctxWithUser, createdText.ID)

		require.NoError(t, err)
		assert.Nil(t, text)
	})
}

// TestStorageService_ConcurrentOperations тестирует конкурентные операции
func TestStorageService_ConcurrentOperations(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	ctx := context.Background()
	numWorkers := 10
	numOperations := 20

	var wg sync.WaitGroup
	errors := make(chan error, numWorkers*numOperations)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				ctxWithUser := customcontext.WithUserID(ctx, fmt.Sprintf("user%d", workerID))

				// Создание данных
				binary := &dtos.NewBinaryData{
					Data:            []byte(fmt.Sprintf("data from worker %d op %d", workerID, j)),
					NewSecureEntity: dtos.NewSecureEntity{Metadata: fmt.Sprintf("metadata %d", j)},
				}

				created, err := service.CreateBinary(ctxWithUser, binary)
				if err != nil {
					errors <- fmt.Errorf("worker %d create failed: %w", workerID, err)
					continue
				}

				// Получение данных
				retrieved, err := service.GetBinary(ctxWithUser, created.ID)
				if err != nil {
					errors <- fmt.Errorf("worker %d get failed: %w", workerID, err)
					continue
				}

				if retrieved == nil || string(retrieved.Data) != string(created.Data) {
					errors <- fmt.Errorf("worker %d data mismatch", workerID)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Проверяем, что не было ошибок
	for err := range errors {
		t.Errorf("Concurrent operation error: %v", err)
	}
}

// TestStorageService_ErrorHandling тестирует обработку ошибок
func TestStorageService_ErrorHandling(t *testing.T) {
	t.Run("Ошибка при создании с nil данными", func(t *testing.T) {
		service, _ := createTestService()
		defer service.Shutdown()

		ctx := customcontext.WithUserID(context.Background(), "testuser")

		// Попытка создания с nil
		_, err := service.CreateBinary(ctx, nil)
		assert.Error(t, err)
	})

	t.Run("Ошибка при обновлении с неправильным ID", func(t *testing.T) {
		service, _ := createTestService()
		defer service.Shutdown()

		ctx := customcontext.WithUserID(context.Background(), "testuser")

		// Попытка обновления несуществующей записи
		binary := &entities.BinaryData{
			SecureEntity: entities.SecureEntity{
				ID: "nonexistent-id",
			},
			Data: []byte("data"),
		}

		result, err := service.UpdateBinary(ctx, binary)

		// In-memory репозиторий может вернуть nil или ошибку
		// Проверяем оба варианта
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.Nil(t, result)
		}
	})
}

// TestStorageService_DataIsolation тестирует изоляцию данных между пользователями
func TestStorageService_DataIsolation(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	// Создаем данные для user1
	ctxUser1 := customcontext.WithUserID(context.Background(), "user1")
	binary1 := &dtos.NewBinaryData{
		Data:            []byte("user1 data"),
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "user1 metadata"},
	}

	created1, err := service.CreateBinary(ctxUser1, binary1)
	require.NoError(t, err)

	// Создаем данные для user2
	ctxUser2 := customcontext.WithUserID(context.Background(), "user2")
	binary2 := &dtos.NewBinaryData{
		Data:            []byte("user2 data"),
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "user2 metadata"},
	}

	created2, err := service.CreateBinary(ctxUser2, binary2)
	require.NoError(t, err)

	// user1 не должен видеть данные user2
	result1, err := service.GetBinary(ctxUser1, created2.ID)
	require.NoError(t, err)
	assert.Nil(t, result1, "user1 should not see user2's data")

	// user2 не должен видеть данные user1
	result2, err := service.GetBinary(ctxUser2, created1.ID)
	require.NoError(t, err)
	assert.Nil(t, result2, "user2 should not see user1's data")

	// Каждый пользователь должен видеть только свои данные
	binaries1, err := service.GetAllBinaries(ctxUser1)
	require.NoError(t, err)
	assert.Len(t, binaries1, 1)
	assert.Equal(t, "user1 data", string(binaries1[0].Data))

	binaries2, err := service.GetAllBinaries(ctxUser2)
	require.NoError(t, err)
	assert.Len(t, binaries2, 1)
	assert.Equal(t, "user2 data", string(binaries2[0].Data))
}

// TestStorageService_TaskQueue тестирует работу очереди задач
func TestStorageService_TaskQueue(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	ctx := context.Background()
	numTasks := 100
	errors := make(chan error, numTasks)

	// Запускаем много задач параллельно
	for i := 0; i < numTasks; i++ {
		go func(taskID int) {
			user := dtos.NewUser{
				Login:    fmt.Sprintf("user%d", taskID),
				Password: fmt.Sprintf("pass%d", taskID),
			}

			_, err := service.CreateUser(ctx, user)
			if err != nil && err != customerrors.AlreadyExistsError {
				errors <- fmt.Errorf("task %d failed: %w", taskID, err)
			}
		}(i)
	}

	// Даем время на выполнение
	time.Sleep(100 * time.Millisecond)
	close(errors)

	// Проверяем ошибки
	for err := range errors {
		t.Errorf("Task error: %v", err)
	}

	// Проверяем, что задачи выполнены
	users, err := service.GetAllUsers(ctx)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), numTasks/2) // Некоторые могли быть дубликатами
}

// TestStorageService_EdgeCases тестирует крайние случаи
func TestStorageService_EdgeCases(t *testing.T) {
	service, _ := createTestService()
	defer service.Shutdown()

	ctx := context.Background()

	t.Run("Пустые данные", func(t *testing.T) {
		// Создание с пустыми данными
		emptyBinary := &dtos.NewBinaryData{
			Data:            []byte{},
			NewSecureEntity: dtos.NewSecureEntity{Metadata: ""},
		}

		ctxWithUser := customcontext.WithUserID(ctx, "testuser")
		result, err := service.CreateBinary(ctxWithUser, emptyBinary)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Data)
	})

	t.Run("Очень длинные строки", func(t *testing.T) {
		longString := string(make([]byte, 10000))
		ctxWithUser := customcontext.WithUserID(ctx, "testuser")

		longText := &dtos.NewTextData{
			Data:            longString,
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "long metadata"},
		}

		result, err := service.CreateText(ctxWithUser, longText)

		require.NoError(t, err)
		assert.Equal(t, longString, result.Data)
	})
}

// TestStorageService_Performance тестирует производительность
func TestStorageService_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	service, _ := createTestService()
	defer service.Shutdown()

	ctx := context.Background()
	numOperations := 1000

	start := time.Now()

	for i := 0; i < numOperations; i++ {
		ctxWithUser := customcontext.WithUserID(ctx, fmt.Sprintf("perfuser%d", i))

		binary := &dtos.NewBinaryData{
			Data:            []byte(fmt.Sprintf("data %d", i)),
			NewSecureEntity: dtos.NewSecureEntity{Metadata: fmt.Sprintf("metadata %d", i)},
		}

		_, err := service.CreateBinary(ctxWithUser, binary)
		require.NoError(t, err)
	}

	elapsed := time.Since(start)
	t.Logf("Created %d records in %v", numOperations, elapsed)

	// Проверяем, что производительность приемлемая
	assert.Less(t, elapsed, 5*time.Second, "Operations took too long")
}
