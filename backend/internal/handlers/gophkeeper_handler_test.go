package handlers_test

import (
	"bytes"
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
	"github.com/JustScorpio/GophKeeper/backend/internal/utils"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестовые пользователи
var testUsers = map[string]string{
	"user1": "password123",
	"user2": "password456",
	"user3": "password789",
}

// getTestData - возвращает тестовые данные
func getTestData() struct {
	binary      dtos.NewBinaryData
	card        dtos.NewCardInformation
	credentials dtos.NewCredentials
	text        dtos.NewTextData
} {
	return struct {
		binary      dtos.NewBinaryData
		card        dtos.NewCardInformation
		credentials dtos.NewCredentials
		text        dtos.NewTextData
	}{
		binary: dtos.NewBinaryData{
			Data:            []byte("test binary data"),
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "test metadata"},
		},
		card: dtos.NewCardInformation{
			Number:          "4111111111111111",
			CardHolder:      "John Doe",
			ExpirationDate:  "12/25",
			CVV:             "123",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "personal card"},
		},
		credentials: dtos.NewCredentials{
			Login:           "serviceuser",
			Password:        "servicepassword",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "test credentials"},
		},
		text: dtos.NewTextData{
			Data:            "This is a test text data",
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "test text metadata"},
		},
	}
}

// createTestHandlerAndRouter - создает тестовый хэндлер и роутер
func createTestHandlerAndRouter() (*chi.Mux, *inmemory.DatabaseManager) {
	dbManager := inmemory.NewDatabaseManager()
	service := services.NewStorageService(
		dbManager.Users,
		dbManager.Binaries,
		dbManager.Cards,
		dbManager.Credentials,
		dbManager.Texts,
	)
	handler := handlers.NewGophkeeperHandler(service)

	// Создаем роутер
	router := chi.NewRouter()

	// Публичные маршруты
	router.Post("/register", handler.Register)
	router.Post("/login", handler.Login)

	// Защищенные маршруты
	router.Route("/api/user", func(r chi.Router) {
		// Используем middleware для аутентификации
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// TODO: добавить проверку токенов (для тестов пропускаем)
				next.ServeHTTP(w, r)
			})
		})

		// Binary data endpoints
		r.Post("/binaries", handler.CreateBinary)
		r.Get("/binaries", handler.GetAllBinaries)
		r.Get("/binaries/{id}", handler.GetBinary)
		r.Put("/binaries", handler.UpdateBinary)
		r.Delete("/binaries/{id}", handler.DeleteBinary)

		// Card information endpoints
		r.Post("/card", handler.CreateCard)
		r.Get("/cards", handler.GetAllCards)
		r.Get("/cards/{id}", handler.GetCard)
		r.Put("/card", handler.UpdateCard)
		r.Delete("/cards/{id}", handler.DeleteCard)

		// Credentials endpoints
		r.Post("/credentials", handler.CreateCredentials)
		r.Get("/credentials", handler.GetAllCredentials)
		r.Get("/credentials/{id}", handler.GetCredentials)
		r.Put("/credentials", handler.UpdateCredentials)
		r.Delete("/credentials/{id}", handler.DeleteCredentials)

		// Text data endpoints
		r.Post("/texts", handler.CreateText)
		r.Get("/texts", handler.GetAllTexts)
		r.Get("/texts/{id}", handler.GetText)
		r.Put("/texts", handler.UpdateText)
		r.Delete("/texts/{id}", handler.DeleteText)
	})

	return router, dbManager
}

// createTestRequest - создать тестовый запрос с авторизацией
func createTestRequest(method, url string, body interface{}, authRequired bool, userLogin string) *http.Request {
	var req *http.Request

	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, url, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	if authRequired {
		ctx := customcontext.WithUserID(req.Context(), userLogin)
		req = req.WithContext(ctx)
	}

	return req
}

// registerTestUser - регистрирует пользователя через роутер
func registerTestUser(t *testing.T, router *chi.Mux, login, password string) {
	newUser := dtos.NewUser{
		Login:    login,
		Password: password,
	}

	req := createTestRequest("POST", "/register", newUser, false, "")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "Failed to register user %s", login)
}

// createBinary - создает бинарные данные через роутер
func createBinary(t *testing.T, router *chi.Mux, userLogin string, data dtos.NewBinaryData) entities.BinaryData {
	req := createTestRequest("POST", "/api/user/binaries", data, true, userLogin)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code, "Failed to create binary data for user %s", userLogin)

	var response entities.BinaryData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return response
}

// createCard - создает данные карты через роутер
func createCard(t *testing.T, router *chi.Mux, userLogin string, data dtos.NewCardInformation) entities.CardInformation {
	req := createTestRequest("POST", "/api/user/card", data, true, userLogin)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code, "Failed to create card data for user %s", userLogin)

	var response entities.CardInformation
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return response
}

// createCredentials - создает учетные данные через роутер
func createCredentials(t *testing.T, router *chi.Mux, userLogin string, data dtos.NewCredentials) entities.Credentials {
	req := createTestRequest("POST", "/api/user/credentials", data, true, userLogin)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code, "Failed to create credentials for user %s", userLogin)

	var response entities.Credentials
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return response
}

// createText - создает текстовые данные через роутер
func createText(t *testing.T, router *chi.Mux, userLogin string, data dtos.NewTextData) entities.TextData {
	req := createTestRequest("POST", "/api/user/texts", data, true, userLogin)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code, "Failed to create text data for user %s", userLogin)

	var response entities.TextData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	return response
}

// TestRegisterAndLogin - ТЕСТЫ РЕГИСТРАЦИИ И АУТЕНТИФИКАЦИИ
func TestRegisterAndLogin(t *testing.T) {
	router, dbManager := createTestHandlerAndRouter()

	t.Run("Успешная регистрация пользователя", func(t *testing.T) {
		newUser := dtos.NewUser{
			Login:    "newuser1",
			Password: "newpassword123",
		}

		req := createTestRequest("POST", "/register", newUser, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Set-Cookie"), "jwt_token")

		// Проверяем через репозиторий, что пользователь создан
		user, err := dbManager.Users.Get(req.Context(), newUser.Login)
		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, newUser.Login, user.Login)
		assert.NotEqual(t, newUser.Password, user.Password)
		assert.True(t, utils.CheckPasswordHash(newUser.Password, user.Password))
	})

	t.Run("Регистрация с существующим логином", func(t *testing.T) {
		// Сначала регистрируем пользователя
		registerTestUser(t, router, "existinguser", "password123")

		// Пытаемся зарегистрироваться с тем же логином
		existingUser := dtos.NewUser{
			Login:    "existinguser",
			Password: "differentpassword",
		}

		req := createTestRequest("POST", "/register", existingUser, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "already exists")
	})

	t.Run("Успешная аутентификация", func(t *testing.T) {
		// Сначала регистрируем пользователя
		registerTestUser(t, router, "authuser", "authpassword")

		// Пытаемся залогиниться
		loginReq := map[string]string{
			"login":    "authuser",
			"password": "authpassword",
		}

		req := createTestRequest("POST", "/login", loginReq, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Set-Cookie"), "jwt_token")
	})

	t.Run("Неуспешная аутентификация - неверный пароль", func(t *testing.T) {
		registerTestUser(t, router, "wrongpassuser", "correctpass")

		loginReq := map[string]string{
			"login":    "wrongpassuser",
			"password": "wrongpassword",
		}

		req := createTestRequest("POST", "/login", loginReq, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})

	t.Run("Неуспешная аутентификация - пользователь не существует", func(t *testing.T) {
		loginReq := map[string]string{
			"login":    "nonexistentuser",
			"password": "somepassword",
		}

		req := createTestRequest("POST", "/login", loginReq, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})
}

// TestBinaryDataCRUD - ТЕСТЫ БИНАРНЫХ ДАННЫХ
func TestBinaryDataCRUD(t *testing.T) {
	router, _ := createTestHandlerAndRouter()
	testData := getTestData()

	// Регистрируем пользователя
	registerTestUser(t, router, "user1", testUsers["user1"])

	t.Run("Создание бинарных данных", func(t *testing.T) {
		req := createTestRequest("POST", "/api/user/binaries", testData.binary, true, "user1")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.BinaryData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testData.binary.Data, response.Data)
		assert.Equal(t, testData.binary.Metadata, response.Metadata)
	})

	t.Run("Получение бинарных данных по ID", func(t *testing.T) {
		// Сначала создаем запись
		binary := createBinary(t, router, "user1", testData.binary)

		// Получаем данные по ID
		req := createTestRequest("GET", fmt.Sprintf("/api/user/binaries/%s", binary.ID),
			nil, true, "user1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.BinaryData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, binary.ID, response.ID)
		assert.Equal(t, binary.Data, response.Data)
	})

	t.Run("Создание нескольких бинарных данных и получение всех", func(t *testing.T) {
		// Создаем несколько записей
		for i := 0; i < 3; i++ {
			binaryData := dtos.NewBinaryData{
				Data:            []byte(fmt.Sprintf("data %d", i)),
				NewSecureEntity: dtos.NewSecureEntity{Metadata: fmt.Sprintf("metadata %d", i)},
			}
			createBinary(t, router, "user1", binaryData)
		}

		// Получаем все данные
		req := createTestRequest("GET", "/api/user/binaries", nil, true, "user1")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entities.BinaryData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, len(response) >= 3) // Могут быть данные из предыдущих тестов
	})

	t.Run("Обновление бинарных данных", func(t *testing.T) {
		// Сначала создаем запись
		binary := createBinary(t, router, "user1", testData.binary)

		// Обновляем
		updateData := entities.BinaryData{
			Data:         []byte("updated data"),
			SecureEntity: entities.SecureEntity{ID: binary.ID, Metadata: "updated metadata"},
		}

		req := createTestRequest("PUT", "/api/user/binaries", updateData, true, "user1")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entities.BinaryData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, []byte("updated data"), response.Data)
		assert.Equal(t, "updated metadata", response.Metadata)
	})

	t.Run("Удаление бинарных данных", func(t *testing.T) {
		// Сначала создаем запись
		binary := createBinary(t, router, "user1", testData.binary)

		// Удаляем
		req := createTestRequest("DELETE", fmt.Sprintf("/api/user/binaries/%s", binary.ID),
			nil, true, "user1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusGone, w.Code)

		// Пытаемся получить удаленные данные
		getReq := createTestRequest("GET", fmt.Sprintf("/api/user/binaries/%s", binary.ID),
			nil, true, "user1")
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		assert.Equal(t, http.StatusNotFound, getW.Code)
	})
}

// TestCardDataCRUD - ТЕСТЫ БАНКОВСКИХ КАРТ
func TestCardDataCRUD(t *testing.T) {
	router, _ := createTestHandlerAndRouter()
	testData := getTestData()

	// Регистрируем пользователей
	registerTestUser(t, router, "user1", testUsers["user1"])
	registerTestUser(t, router, "user2", testUsers["user2"])

	t.Run("Создание данных банковской карты", func(t *testing.T) {
		req := createTestRequest("POST", "/api/user/card", testData.card, true, "user1")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.CardInformation
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testData.card.Number, response.Number)
		assert.Equal(t, testData.card.CardHolder, response.CardHolder)
		assert.Equal(t, testData.card.ExpirationDate, response.ExpirationDate)
		assert.Equal(t, testData.card.CVV, response.CVV)
	})

	t.Run("Изоляция данных между пользователями", func(t *testing.T) {
		// user1 создает карту
		card := createCard(t, router, "user1", testData.card)

		// user2 пытается получить карту user1
		req := createTestRequest("GET", fmt.Sprintf("/api/user/cards/%s", card.ID),
			nil, true, "user2")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// user2 не должен видеть карту user1
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "not found")

		// user1 должен видеть свою карту
		reqUser1 := createTestRequest("GET", fmt.Sprintf("/api/user/cards/%s", card.ID),
			nil, true, "user1")
		wUser1 := httptest.NewRecorder()
		router.ServeHTTP(wUser1, reqUser1)

		assert.Equal(t, http.StatusOK, wUser1.Code)
	})
}

// TestCredentialsCRUD - ТЕСТЫ УЧЕТНЫХ ДАННЫХ
func TestCredentialsCRUD(t *testing.T) {
	router, _ := createTestHandlerAndRouter()
	testData := getTestData()

	// Регистрируем пользователя
	registerTestUser(t, router, "user1", testUsers["user1"])

	t.Run("Создание учетных данных", func(t *testing.T) {
		req := createTestRequest("POST", "/api/user/credentials", testData.credentials, true, "user1")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.Credentials
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testData.credentials.Login, response.Login)
		assert.Equal(t, testData.credentials.Password, response.Password)
		assert.Equal(t, testData.credentials.Metadata, response.Metadata)
	})

	t.Run("Получение всех учетных данных пользователя", func(t *testing.T) {
		// Создаем несколько записей
		for i := 0; i < 2; i++ {
			cred := dtos.NewCredentials{
				Login:           fmt.Sprintf("user%d", i),
				Password:        fmt.Sprintf("pass%d", i),
				NewSecureEntity: dtos.NewSecureEntity{Metadata: fmt.Sprintf("cred %d", i)},
			}
			createCredentials(t, router, "user1", cred)
		}

		// Получаем все
		req := createTestRequest("GET", "/api/user/credentials", nil, true, "user1")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []entities.Credentials
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, len(response) >= 2) // Могут быть данные из предыдущих тестов
	})
}

// TestTextDataCRUD - ТЕСТЫ ТЕКСТОВЫХ ДАННЫХ
func TestTextDataCRUD(t *testing.T) {
	router, _ := createTestHandlerAndRouter()
	testData := getTestData()

	// Регистрируем пользователя
	registerTestUser(t, router, "user1", testUsers["user1"])

	t.Run("Создание текстовых данных", func(t *testing.T) {
		req := createTestRequest("POST", "/api/user/texts", testData.text, true, "user1")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entities.TextData
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testData.text.Data, response.Data)
		assert.Equal(t, testData.text.Metadata, response.Metadata)
	})

	t.Run("Удаление текстовых данных", func(t *testing.T) {
		// Сначала создаем запись
		text := createText(t, router, "user1", testData.text)

		// Удаляем
		req := createTestRequest("DELETE", fmt.Sprintf("/api/user/texts/%s", text.ID),
			nil, true, "user1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusGone, w.Code)

		// Пытаемся получить удаленные данные
		getReq := createTestRequest("GET", fmt.Sprintf("/api/user/texts/%s", text.ID),
			nil, true, "user1")
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		assert.Equal(t, http.StatusNotFound, getW.Code)
	})
}

// TestUnauthorizedAccess - ТЕСТЫ БЕЗ АВТОРИЗАЦИИ
func TestUnauthorizedAccess(t *testing.T) {
	router, _ := createTestHandlerAndRouter()
	testData := getTestData()

	tests := []struct {
		name   string
		method string
		path   string
		body   interface{}
	}{
		{
			name:   "CreateBinary без авторизации",
			method: "POST",
			path:   "/api/user/binaries",
			body:   testData.binary,
		},
		{
			name:   "GetAllBinaries без авторизации",
			method: "GET",
			path:   "/api/user/binaries",
			body:   nil,
		},
		{
			name:   "CreateCard без авторизации",
			method: "POST",
			path:   "/api/user/card",
			body:   testData.card,
		},
		{
			name:   "CreateCredentials без авторизации",
			method: "POST",
			path:   "/api/user/credentials",
			body:   testData.credentials,
		},
		{
			name:   "CreateText без авторизации",
			method: "POST",
			path:   "/api/user/texts",
			body:   testData.text,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createTestRequest(tt.method, tt.path, tt.body, false, "")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
			assert.Contains(t, w.Body.String(), "Authentication required")
		})
	}
}

// TestErrorScenarios - ТЕСТЫ ОШИБОЧНЫХ СЦЕНАРИЕВ
func TestErrorScenarios(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	// Регистрируем пользователя
	registerTestUser(t, router, "user1", testUsers["user1"])

	t.Run("Получение несуществующих данных", func(t *testing.T) {
		req := createTestRequest("GET", "/api/user/binaries/999999",
			nil, true, "user1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "not found")
	})

	t.Run("Некорректный JSON в запросе", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/user/binaries",
			bytes.NewReader([]byte("{invalid json}")))
		req.Header.Set("Content-Type", "application/json")
		ctx := customcontext.WithUserID(req.Context(), "user1")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestMultiUserEnvironment - ТЕСТЫ МНОГОПОЛЬЗОВАТЕЛЬСКОЙ СРЕДЫ
func TestMultiUserEnvironment(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	t.Run("Несколько пользователей создают и получают свои данные", func(t *testing.T) {
		users := []struct {
			login    string
			password string
			count    int
		}{
			{"multiuser1", "password1", 3},
			{"multiuser2", "password2", 2},
			{"multiuser3", "password3", 1},
		}

		// Регистрируем пользователей
		for _, u := range users {
			userDTO := dtos.NewUser{
				Login:    u.login,
				Password: u.password,
			}

			req := createTestRequest("POST", "/register", userDTO, false, "")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code, "User %s should register successfully", u.login)
		}

		// Каждый пользователь создает свои данные
		for _, u := range users {
			for i := 0; i < u.count; i++ {
				binary := dtos.NewBinaryData{
					Data:            []byte(fmt.Sprintf("data from %s #%d", u.login, i+1)),
					NewSecureEntity: dtos.NewSecureEntity{Metadata: fmt.Sprintf("metadata from %s", u.login)},
				}

				req := createTestRequest("POST", "/api/user/binaries", binary, true, u.login)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusCreated, w.Code,
					"User %s should be able to create binary data", u.login)
			}
		}

		// Теперь каждый пользователь получает свои данные
		for _, u := range users {
			req := createTestRequest("GET", "/api/user/binaries", nil, true, u.login)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code,
				"User %s should be able to get all binaries", u.login)

			var response []entities.BinaryData
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Len(t, response, u.count,
				"User %s should have %d binaries, got %d", u.login, u.count, len(response))
		}
	})
}

// TestPasswordHashing - ТЕСТЫ ХЭШИРОВАНИЯ ПАРОЛЕЙ
func TestPasswordHashing(t *testing.T) {
	t.Run("Разные пароли дают разные хэши", func(t *testing.T) {
		hash1, err1 := utils.HashPassword("password1")
		hash2, err2 := utils.HashPassword("password2")

		require.NoError(t, err1)
		require.NoError(t, err2)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Одинаковые пароли дают разные хэши (из-за соли)", func(t *testing.T) {
		hash1, err1 := utils.HashPassword("samepassword")
		hash2, err2 := utils.HashPassword("samepassword")

		require.NoError(t, err1)
		require.NoError(t, err2)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Проверка правильного пароля", func(t *testing.T) {
		password := "testpassword"
		hash, err := utils.HashPassword(password)
		require.NoError(t, err)

		assert.True(t, utils.CheckPasswordHash(password, hash))
	})

	t.Run("Проверка неправильного пароля", func(t *testing.T) {
		password := "testpassword"
		hash, err := utils.HashPassword(password)
		require.NoError(t, err)

		assert.False(t, utils.CheckPasswordHash("wrongpassword", hash))
	})
}

// TestCompleteUserScenario - ПОЛНЫЙ ЦИКЛ ОПЕРАЦИЙ
func TestCompleteUserScenario(t *testing.T) {
	router, _ := createTestHandlerAndRouter()
	testData := getTestData()

	t.Run("Полный сценарий работы пользователя", func(t *testing.T) {
		// 1. Регистрация нового пользователя
		newUser := dtos.NewUser{
			Login:    "newuser",
			Password: "newpassword",
		}

		registerReq := createTestRequest("POST", "/register", newUser, false, "")
		registerW := httptest.NewRecorder()
		router.ServeHTTP(registerW, registerReq)
		assert.Equal(t, http.StatusOK, registerW.Code)

		// 2. Логин с новыми учетными данными
		loginReq := map[string]string{
			"login":    "newuser",
			"password": "newpassword",
		}

		loginHttpReq := createTestRequest("POST", "/login", loginReq, false, "")
		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, loginHttpReq)
		assert.Equal(t, http.StatusOK, loginW.Code)

		// 3. Создание различных данных
		binaryReq := createTestRequest("POST", "/api/user/binaries", testData.binary, true, "newuser")
		binaryW := httptest.NewRecorder()
		router.ServeHTTP(binaryW, binaryReq)
		assert.Equal(t, http.StatusCreated, binaryW.Code)

		var binaryResp entities.BinaryData
		err := json.Unmarshal(binaryW.Body.Bytes(), &binaryResp)
		require.NoError(t, err)

		cardReq := createTestRequest("POST", "/api/user/card", testData.card, true, "newuser")
		cardW := httptest.NewRecorder()
		router.ServeHTTP(cardW, cardReq)
		assert.Equal(t, http.StatusCreated, cardW.Code)

		var cardResp entities.CardInformation
		err = json.Unmarshal(cardW.Body.Bytes(), &cardResp)
		require.NoError(t, err)

		// 4. Получение созданных данных
		getBinaryReq := createTestRequest("GET", fmt.Sprintf("/api/user/binaries/%s", binaryResp.ID),
			nil, true, "newuser")
		getBinaryW := httptest.NewRecorder()
		router.ServeHTTP(getBinaryW, getBinaryReq)
		assert.Equal(t, http.StatusOK, getBinaryW.Code)

		getCardReq := createTestRequest("GET", fmt.Sprintf("/api/user/cards/%s", cardResp.ID),
			nil, true, "newuser")
		getCardW := httptest.NewRecorder()
		router.ServeHTTP(getCardW, getCardReq)
		assert.Equal(t, http.StatusOK, getCardW.Code)

		// 5. Проверка, что данные созданы
		assert.NotEmpty(t, binaryResp.ID)
		assert.NotEmpty(t, cardResp.ID)
	})
}

// TestMiddleware - тестирование middleware
func TestMiddleware(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	t.Run("Проверка работы middleware авторизации", func(t *testing.T) {
		// Пытаемся получить доступ без авторизации
		req := createTestRequest("GET", "/api/user/binaries", nil, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Middleware должен вернуть 401
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
