package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/handlers"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/backend/internal/repositories/inmemory"
	"github.com/JustScorpio/GophKeeper/backend/internal/services"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Создание тестовых данных
var (
	testUser = dtos.NewUser{
		Login:    "testuser",
		Password: "password123",
	}

	testBinary = dtos.NewBinaryData{
		Data:            []byte("test binary data"),
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "test metadata"},
	}

	testCard = dtos.NewCardInformation{
		Number:          "4111111111111111",
		CardHolder:      "John Doe",
		ExpirationDate:  "12/25",
		CVV:             "123",
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "personal card"},
	}

	testCredentials = dtos.NewCredentials{
		Login:           "serviceuser",
		Password:        "servicepassword",
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "test credentials"},
	}

	testText = dtos.NewTextData{
		Data:            "This is a test text data",
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "test text metadata"},
	}
)

// InMemoryRepositories - структура с прямым доступом к репозиториям
type InMemoryRepositories struct {
	Users       *inmemory.InMemoryUsersRepo
	Binaries    *inmemory.InMemoryBinariesRepo
	Cards       *inmemory.InMemoryCardsRepo
	Credentials *inmemory.InMemoryCredentialsRepo
	Texts       *inmemory.InMemoryTextsRepo
}

// createTestHandler - создать тестовый хэндлер
func createTestHandler() (*handlers.GophkeeperHandler, *InMemoryRepositories) {
	// Создаем тестовый хендлер
	repos := &InMemoryRepositories{
		Users:       inmemory.NewInMemoryUsersRepo(),
		Binaries:    inmemory.NewInMemoryBinariesRepo(),
		Cards:       inmemory.NewInMemoryCardsRepo(),
		Credentials: inmemory.NewInMemoryCredentialsRepo(),
		Texts:       inmemory.NewInMemoryTextsRepo(),
	}

	service := services.NewStorageService(
		repos.Users,
		repos.Binaries,
		repos.Cards,
		repos.Credentials,
		repos.Texts,
	)
	return handlers.NewGophkeeperHandler(service), repos
}

// createTestRequest - создать тестовый запрос
func createTestRequest(method, path string, body interface{}, addAuth bool) *http.Request {
	var req *http.Request

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	if addAuth {
		ctx := customcontext.WithUserID(req.Context(), "testuser")
		req = req.WithContext(ctx)
	}

	return req
}

// TestRegisterAndLogin = ТЕСТЫ РЕГИСТРАЦИИ И АУТЕНТИФИКАЦИИ
func TestRegisterAndLogin(t *testing.T) {
	handler, repos := createTestHandler()

	t.Run("Успешная регистрация пользователя", func(t *testing.T) {
		req := createTestRequest("POST", "/register", testUser, false)
		w := httptest.NewRecorder()

		handler.Register(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Set-Cookie"), "jwt_token")

		// Проверяем, что пользователь создался
		user, err := repos.Users.Get(req.Context(), testUser.Login)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, testUser.Login, user.Login)
	})

	t.Run("Успешная аутентификация", func(t *testing.T) {
		// Пытаемся залогиниться
		loginReq := map[string]string{
			"login":    testUser.Login,
			"password": testUser.Password,
		}

		req := createTestRequest("POST", "/login", loginReq, false)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Set-Cookie"), "jwt_token")
	})

	t.Run("Неуспешная аутентификация - неверный пароль", func(t *testing.T) {
		loginReq := map[string]string{
			"login":    testUser.Login,
			"password": "wrongpassword",
		}

		req := createTestRequest("POST", "/login", loginReq, false)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})
}

// TestBinaryDataCRUD - ТЕСТЫ БИНАРНЫХ ДАННЫХ
func TestBinaryDataCRUD(t *testing.T) {
	handler, repos := createTestHandler()

	// Создаем пользователя для тестов
	ctx := context.Background()
	_, err := repos.Users.Create(ctx, &testUser)
	require.NoError(t, err)

	t.Run("Создание бинарных данных", func(t *testing.T) {
		req := createTestRequest("POST", "/binaries", testBinary, true)
		w := httptest.NewRecorder()

		handler.CreateBinary(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.BinaryData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testBinary.Data, response.Data)
		assert.Equal(t, testBinary.Metadata, response.Metadata)
		assert.Equal(t, "testuser", response.OwnerID)
	})

	t.Run("Получение бинарных данных по ID", func(t *testing.T) {
		// Сначала создаем запись
		ctxWithUser := customcontext.WithUserID(context.Background(), "testuser")
		binary, err := repos.Binaries.Create(ctxWithUser, &testBinary)
		require.NoError(t, err)

		req := httptest.NewRequest("GET", fmt.Sprintf("/binaries/%s", binary.ID), nil)
		ctx := customcontext.WithUserID(req.Context(), "testuser")
		req = req.WithContext(ctx)

		// Добавляем параметр ID для chi router
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", binary.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.GetBinary(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.BinaryData
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, binary.ID, response.ID)
		assert.Equal(t, binary.Data, response.Data)
	})

	t.Run("Получение всех бинарных данных пользователя", func(t *testing.T) {
		req := createTestRequest("GET", "/binaries", nil, true)
		w := httptest.NewRecorder()

		handler.GetAllBinaries(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entities.BinaryData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, len(response) > 0)
	})

	t.Run("Обновление бинарных данных", func(t *testing.T) {
		// Сначала создаем запись
		ctxWithUser := customcontext.WithUserID(context.Background(), "testuser")
		binary, err := repos.Binaries.Create(ctxWithUser, &testBinary)
		require.NoError(t, err)

		// Обновляем
		updateData := entities.BinaryData{
			Data:         []byte("updated data"),
			SecureEntity: entities.SecureEntity{ID: binary.ID, Metadata: "updated metadata"},
		}

		req := createTestRequest("PUT", fmt.Sprintf("/binaries/%s", binary.ID), updateData, true)

		// Добавляем параметр ID для chi router
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", binary.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.UpdateBinary(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.BinaryData
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, []byte("updated data"), response.Data)
		assert.Equal(t, "updated metadata", response.Metadata)
	})

	t.Run("Удаление бинарных данных", func(t *testing.T) {
		// Сначала создаем запись
		ctxWithUser := customcontext.WithUserID(context.Background(), "testuser")
		binary, err := repos.Binaries.Create(ctxWithUser, &testBinary)
		require.NoError(t, err)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/binaries/%s", binary.ID), nil)
		ctx := customcontext.WithUserID(req.Context(), "testuser")
		req = req.WithContext(ctx)

		// Добавляем параметр ID для chi router
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", binary.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.DeleteBinary(w, req)

		assert.Equal(t, http.StatusGone, w.Code)

		// Проверяем, что данные удалились
		deletedBinary, err := repos.Binaries.Get(ctxWithUser, binary.ID)
		require.NoError(t, err)
		assert.Nil(t, deletedBinary)
	})
}

// TestCardDataCRUD - ТЕСТЫ БАНКОВСКИХ КАРТ
func TestCardDataCRUD(t *testing.T) {
	handler, repos := createTestHandler()

	// Создаем пользователя для тестов
	ctx := context.Background()
	_, err := repos.Users.Create(ctx, &testUser)
	require.NoError(t, err)

	t.Run("Создание данных банковской карты", func(t *testing.T) {
		req := createTestRequest("POST", "/cards", testCard, true)
		w := httptest.NewRecorder()

		handler.CreateCard(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.CardInformation
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testCard.Number, response.Number)
		assert.Equal(t, testCard.CardHolder, response.CardHolder)
		assert.Equal(t, testCard.ExpirationDate, response.ExpirationDate)
		assert.Equal(t, testCard.CVV, response.CVV)
		assert.Equal(t, "testuser", response.OwnerID)
	})

	t.Run("Получение данных карты по ID", func(t *testing.T) {
		// Сначала создаем запись
		ctxWithUser := customcontext.WithUserID(context.Background(), "testuser")
		card, err := repos.Cards.Create(ctxWithUser, &testCard)
		require.NoError(t, err)

		req := httptest.NewRequest("GET", fmt.Sprintf("/cards/%s", card.ID), nil)
		ctx := customcontext.WithUserID(req.Context(), "testuser")
		req = req.WithContext(ctx)

		// Добавляем параметр ID для chi router
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", card.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.GetCard(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.CardInformation
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, card.ID, response.ID)
		assert.Equal(t, card.Number, response.Number)
	})

	t.Run("Изоляция данных между пользователями", func(t *testing.T) {
		// Создаем второго пользователя
		user2 := dtos.NewUser{
			Login:    "user2",
			Password: "password2",
		}
		_, err := repos.Users.Create(context.Background(), &user2)
		require.NoError(t, err)

		// user1 создает карту
		ctxUser1 := customcontext.WithUserID(context.Background(), "testuser")
		card, err := repos.Cards.Create(ctxUser1, &testCard)
		require.NoError(t, err)

		// user2 пытается получить карту user1
		req := httptest.NewRequest("GET", fmt.Sprintf("/cards/%s", card.ID), nil)
		ctx := customcontext.WithUserID(req.Context(), "user2")
		req = req.WithContext(ctx)

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", card.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.GetCard(w, req)

		// user2 не должен видеть карту user1
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "not found")
	})
}

// TestCredentialsCRUD - ТЕСТЫ УЧЕТНЫХ ДАННЫХ
func TestCredentialsCRUD(t *testing.T) {
	handler, repos := createTestHandler()

	// Создаем пользователя для тестов
	ctx := context.Background()
	_, err := repos.Users.Create(ctx, &testUser)
	require.NoError(t, err)

	t.Run("Создание учетных данных", func(t *testing.T) {
		req := createTestRequest("POST", "/credentials", testCredentials, true)
		w := httptest.NewRecorder()

		handler.CreateCredentials(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.Credentials
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testCredentials.Login, response.Login)
		assert.Equal(t, testCredentials.Password, response.Password)
		assert.Equal(t, testCredentials.Metadata, response.Metadata)
		assert.Equal(t, "testuser", response.OwnerID)
	})

	t.Run("Получение всех учетных данных пользователя", func(t *testing.T) {
		req := createTestRequest("GET", "/credentials", nil, true)
		w := httptest.NewRecorder()

		handler.GetAllCredentials(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entities.Credentials
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, len(response) > 0)

		// Проверяем, что все данные принадлежат текущему пользователю
		for _, cred := range response {
			assert.Equal(t, "testuser", cred.OwnerID)
		}
	})
}

// TestTextDataCRUD - ТЕСТЫ ТЕКСТОВЫХ ДАННЫХ
func TestTextDataCRUD(t *testing.T) {
	handler, repos := createTestHandler()

	// Создаем пользователя для тестов
	ctx := context.Background()
	_, err := repos.Users.Create(ctx, &testUser)
	require.NoError(t, err)

	t.Run("Создание текстовых данных", func(t *testing.T) {
		req := createTestRequest("POST", "/texts", testText, true)
		w := httptest.NewRecorder()

		handler.CreateText(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.TextData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testText.Data, response.Data)
		assert.Equal(t, testText.Metadata, response.Metadata)
		assert.Equal(t, "testuser", response.OwnerID)
	})

	t.Run("Удаление текстовых данных", func(t *testing.T) {
		// Сначала создаем запись
		ctxWithUser := customcontext.WithUserID(context.Background(), "testuser")
		text, err := repos.Texts.Create(ctxWithUser, &testText)
		require.NoError(t, err)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/texts/%s", text.ID), nil)
		ctx := customcontext.WithUserID(req.Context(), "testuser")
		req = req.WithContext(ctx)

		// Добавляем параметр ID для chi router
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", text.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.DeleteText(w, req)

		assert.Equal(t, http.StatusGone, w.Code)

		// Проверяем, что данные удалились
		deletedText, err := repos.Texts.Get(ctxWithUser, text.ID)
		require.NoError(t, err)
		assert.Nil(t, deletedText)
	})
}

// TestUnauthorizedAccess - ТЕСТЫ БЕЗ АВТОРИЗАЦИИ
func TestUnauthorizedAccess(t *testing.T) {
	handler, _ := createTestHandler()

	tests := []struct {
		name   string
		method string
		path   string
		body   interface{}
	}{
		{"CreateBinary без авторизации", "POST", "/binaries", testBinary},
		{"GetAllBinaries без авторизации", "GET", "/binaries", nil},
		{"CreateCard без авторизации", "POST", "/cards", testCard},
		{"CreateCredentials без авторизации", "POST", "/credentials", testCredentials},
		{"CreateText без авторизации", "POST", "/texts", testText},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createTestRequest(tt.method, tt.path, tt.body, false)
			w := httptest.NewRecorder()

			switch tt.method {
			case "POST":
				switch tt.path {
				case "/binaries":
					handler.CreateBinary(w, req)
				case "/cards":
					handler.CreateCard(w, req)
				case "/credentials":
					handler.CreateCredentials(w, req)
				case "/texts":
					handler.CreateText(w, req)
				}
			case "GET":
				if tt.path == "/binaries" {
					handler.GetAllBinaries(w, req)
				}
			}

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Contains(t, w.Body.String(), "Authentication required")
		})
	}
}

// TestCompleteUserScenario - ПОЛНЫЙ ЦИКЛ ОПЕРАЦИЙ
func TestCompleteUserScenario(t *testing.T) {
	handler, repos := createTestHandler()

	t.Run("Полный сценарий работы пользователя", func(t *testing.T) {
		// 1. Регистрация нового пользователя
		newUser := dtos.NewUser{
			Login:    "newuser",
			Password: "newpassword",
		}

		registerReq := createTestRequest("POST", "/register", newUser, false)
		registerW := httptest.NewRecorder()
		handler.Register(registerW, registerReq)
		assert.Equal(t, http.StatusOK, registerW.Code)

		// 2. Создание различных данных
		ctxWithUser := customcontext.WithUserID(context.Background(), "newuser")

		// Бинарные данные
		binary, err := repos.Binaries.Create(ctxWithUser, &testBinary)
		require.NoError(t, err)

		// Банковская карта
		card, err := repos.Cards.Create(ctxWithUser, &testCard)
		require.NoError(t, err)

		// Учетные данные
		cred, err := repos.Credentials.Create(ctxWithUser, &testCredentials)
		require.NoError(t, err)

		// Текстовые данные
		text, err := repos.Texts.Create(ctxWithUser, &testText)
		require.NoError(t, err)

		// 3. Проверка, что данные созданы
		assert.NotEmpty(t, binary.ID)
		assert.NotEmpty(t, card.ID)
		assert.NotEmpty(t, cred.ID)
		assert.NotEmpty(t, text.ID)

		// 4. Проверка изоляции - другой пользователь не видит данные
		anotherUserCtx := customcontext.WithUserID(context.Background(), "anotheruser")

		// Пытаемся получить данные нового пользователя
		anotherBinary, err := repos.Binaries.Get(anotherUserCtx, binary.ID)
		require.NoError(t, err)
		assert.Nil(t, anotherBinary)

		anotherCard, err := repos.Cards.Get(anotherUserCtx, card.ID)
		require.NoError(t, err)
		assert.Nil(t, anotherCard)

		// 5. Удаление данных
		err = repos.Binaries.Delete(ctxWithUser, binary.ID)
		require.NoError(t, err)

		err = repos.Cards.Delete(ctxWithUser, card.ID)
		require.NoError(t, err)

		// 6. Проверка, что данные удалены
		deletedBinary, err := repos.Binaries.Get(ctxWithUser, binary.ID)
		require.NoError(t, err)
		assert.Nil(t, deletedBinary)

		deletedCard, err := repos.Cards.Get(ctxWithUser, card.ID)
		require.NoError(t, err)
		assert.Nil(t, deletedCard)

		// 7. Остальные данные остались
		remainingCred, err := repos.Credentials.Get(ctxWithUser, cred.ID)
		require.NoError(t, err)
		assert.NotNil(t, remainingCred)

		remainingText, err := repos.Texts.Get(ctxWithUser, text.ID)
		require.NoError(t, err)
		assert.NotNil(t, remainingText)
	})
}

// TestErrorScenarios - ТЕСТЫ ОШИБОЧНЫХ СЦЕНАРИЕВ
func TestErrorScenarios(t *testing.T) {
	handler, repos := createTestHandler()

	// Создаем пользователя для тестов
	ctx := context.Background()
	_, err := repos.Users.Create(ctx, &testUser)
	require.NoError(t, err)

	t.Run("Получение несуществующих данных", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/binaries/999999", nil)
		ctx := customcontext.WithUserID(req.Context(), "testuser")
		req = req.WithContext(ctx)

		// Добавляем параметр ID для chi router
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.GetBinary(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "not found")
	})

	t.Run("Обновление несуществующих данных", func(t *testing.T) {
		updateData := entities.BinaryData{
			Data:         []byte("updated data"),
			SecureEntity: entities.SecureEntity{ID: "999999", Metadata: "updated metadata"},
		}

		req := createTestRequest("PUT", "/binaries/999999", updateData, true)

		// Добавляем параметр ID для chi router
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "999999")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		w := httptest.NewRecorder()
		handler.UpdateBinary(w, req)

		// Должна быть ошибка
		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

// TestMultiUserEnvironment - ТЕСТЫ МНОГОПОЛЬЗОВАТЕЛЬСКОЙ СРЕДЫ
func TestMultiUserEnvironment(t *testing.T) {
	handler, repos := createTestHandler()

	t.Run("Несколько пользователей создают и получают свои данные через хендлер", func(t *testing.T) {
		users := []struct {
			login string
			count int
		}{
			{"user1", 3},
			{"user2", 2},
			{"user3", 1},
		}

		// Храним ID созданных данных для каждого пользователя
		userDataIDs := make(map[string][]string)

		// Создаем пользователей
		for _, u := range users {
			userDTO := dtos.NewUser{
				Login:    u.login,
				Password: "password",
			}
			_, err := repos.Users.Create(context.Background(), &userDTO)
			require.NoError(t, err)
		}

		// Каждый пользователь создает свои данные
		for _, u := range users {
			for i := 0; i < u.count; i++ {
				binary := dtos.NewBinaryData{
					Data:            []byte(fmt.Sprintf("data from %s #%d", u.login, i+1)),
					NewSecureEntity: dtos.NewSecureEntity{Metadata: fmt.Sprintf("metadata from %s", u.login)},
				}

				body, _ := json.Marshal(binary)
				req := httptest.NewRequest("POST", "/binaries", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				// Добавляем контекст с пользователем
				ctx := customcontext.WithUserID(req.Context(), u.login)
				req = req.WithContext(ctx)

				w := httptest.NewRecorder()
				handler.CreateBinary(w, req)

				// Проверяем успешность создания
				assert.Equal(t, http.StatusCreated, w.Code,
					"User %s should be able to create binary data", u.login)

				// Парсим ответ и сохраняем ID
				var response entities.BinaryData
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				userDataIDs[u.login] = append(userDataIDs[u.login], response.ID)
			}
		}

		// Теперь каждый пользователь получает свои данные
		for _, u := range users {
			// Создаем запрос для получения всех данных пользователя
			req := httptest.NewRequest("GET", "/binaries", nil)

			// Добавляем контекст с пользователем
			ctx := customcontext.WithUserID(req.Context(), u.login)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.GetAllBinaries(w, req)

			// Проверяем успешность получения
			assert.Equal(t, http.StatusOK, w.Code,
				"User %s should be able to get all binaries", u.login)

			// Парсим ответ
			var response []entities.BinaryData
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Проверяем, что пользователь получил правильное количество данных
			assert.Len(t, response, u.count,
				"User %s should have %d binaries, got %d", u.login, u.count, len(response))

			// Проверяем, что все данные принадлежат текущему пользователю
			for _, b := range response {
				assert.Equal(t, u.login, b.OwnerID,
					"Binary should belong to user %s, but owner is %s", u.login, b.OwnerID)
			}

			// Проверяем, что получены именно те данные, которые мы создали (проверяем по содержимому metadata)
			for _, b := range response {
				assert.Contains(t, b.Metadata, u.login,
					"Metadata should contain user login: %s", u.login)
			}
		}

		// Теперь тестируем получение конкретных данных
		for _, u := range users {
			// Для каждого ID данных пользователя
			for _, dataID := range userDataIDs[u.login] {
				// Создаем запрос на получение конкретной записи
				req := httptest.NewRequest("GET", "/binaries/"+dataID, nil)

				// Добавляем контекст с пользователем
				ctx := customcontext.WithUserID(req.Context(), u.login)
				req = req.WithContext(ctx)

				// Добавляем параметр ID для chi router
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", dataID)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				w := httptest.NewRecorder()
				handler.GetBinary(w, req)

				// Проверяем успешность получения
				assert.Equal(t, http.StatusOK, w.Code,
					"User %s should be able to get binary with ID %s", u.login, dataID)

				// Парсим ответ
				var response entities.BinaryData
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Проверяем, что получены правильные данные
				assert.Equal(t, dataID, response.ID)
				assert.Equal(t, u.login, response.OwnerID)
			}
		}

		// Тестируем изоляцию данных - пользователь не должен видеть чужие данные
		for i, currentUser := range users {
			// Пытаемся получить данные другого пользователя
			otherUser := users[(i+1)%len(users)] // Берем следующего пользователя в списке

			if len(userDataIDs[otherUser.login]) > 0 {
				otherUserDataID := userDataIDs[otherUser.login][0] // Берем первый ID данных другого пользователя

				// Создаем запрос от имени текущего пользователя на данные другого пользователя
				req := httptest.NewRequest("GET", "/binaries/"+otherUserDataID, nil)

				// Добавляем контекст с текущим пользователем
				ctx := customcontext.WithUserID(req.Context(), currentUser.login)
				req = req.WithContext(ctx)

				// Добавляем параметр ID для chi router
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("id", otherUserDataID)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

				w := httptest.NewRecorder()
				handler.GetBinary(w, req)

				// Текущий пользователь НЕ должен видеть данные другого пользователя
				// В текущей реализации хендлер возвращает 404 Not Found
				assert.Equal(t, http.StatusNotFound, w.Code,
					"User %s should NOT be able to get binary %s of user %s",
					currentUser.login, otherUserDataID, otherUser.login)
			}
		}
	})
}
