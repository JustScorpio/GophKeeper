// handlers_test - тесты хэндлеров
package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/handlers"
	"github.com/JustScorpio/GophKeeper/backend/internal/hash"
	"github.com/JustScorpio/GophKeeper/backend/internal/middleware/auth"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/backend/internal/repositories/inmemory"
	"github.com/JustScorpio/GophKeeper/backend/internal/services"
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

	//Инициализация аутентификатора
	auth.Init("649f7b24-76ea-4f15-bd32-31fab91d63f6")

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
		assert.True(t, hash.CheckPasswordHash(newUser.Password, user.Password))
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
		hash1, err1 := hash.HashPassword("password1")
		hash2, err2 := hash.HashPassword("password2")

		require.NoError(t, err1)
		require.NoError(t, err2)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Одинаковые пароли дают разные хэши (из-за соли)", func(t *testing.T) {
		hash1, err1 := hash.HashPassword("samepassword")
		hash2, err2 := hash.HashPassword("samepassword")

		require.NoError(t, err1)
		require.NoError(t, err2)

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Проверка правильного пароля", func(t *testing.T) {
		password := "testpassword"
		hashpass, err := hash.HashPassword(password)
		require.NoError(t, err)

		assert.True(t, hash.CheckPasswordHash(password, hashpass))
	})

	t.Run("Проверка неправильного пароля", func(t *testing.T) {
		password := "testpassword"
		hashpass, err := hash.HashPassword(password)
		require.NoError(t, err)

		assert.False(t, hash.CheckPasswordHash("wrongpassword", hashpass))
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

// TestContentTypeValidation - тесты валидации Content-Type
func TestContentTypeValidation(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	// Регистрируем пользователя
	registerTestUser(t, router, "ctuser", "password123")

	testData := getTestData()

	tests := []struct {
		name        string
		method      string
		path        string
		contentType string
		body        interface{}
		auth        bool
		userLogin   string
		expected    int
	}{
		{
			name:        "CreateBinary без Content-Type",
			method:      "POST",
			path:        "/api/user/binaries",
			contentType: "",
			body:        testData.binary,
			auth:        true,
			userLogin:   "ctuser",
			expected:    http.StatusBadRequest,
		},
		{
			name:        "CreateBinary с неправильным Content-Type",
			method:      "POST",
			path:        "/api/user/binaries",
			contentType: "text/plain",
			body:        testData.binary,
			auth:        true,
			userLogin:   "ctuser",
			expected:    http.StatusBadRequest,
		},
		{
			name:        "Register без Content-Type",
			method:      "POST",
			path:        "/register",
			contentType: "",
			body:        dtos.NewUser{Login: "test", Password: "test"},
			auth:        false,
			userLogin:   "",
			expected:    http.StatusBadRequest,
		},
		{
			name:        "Login с неправильным Content-Type",
			method:      "POST",
			path:        "/login",
			contentType: "application/xml",
			body:        map[string]string{"login": "test", "password": "test"},
			auth:        false,
			userLogin:   "",
			expected:    http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader(jsonBody))
				if tt.contentType != "" {
					req.Header.Set("Content-Type", tt.contentType)
				}
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
				if tt.contentType != "" {
					req.Header.Set("Content-Type", tt.contentType)
				}
			}

			if tt.auth {
				ctx := customcontext.WithUserID(req.Context(), tt.userLogin)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expected, w.Code,
				"Expected status %d, got %d for %s", tt.expected, w.Code, tt.name)
		})
	}
}

// TestMethodNotAllowed - тесты неподдерживаемых методов HTTP
func TestMethodNotAllowed(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	// Регистрируем пользователя
	registerTestUser(t, router, "methoduser", "password123")

	testData := getTestData()

	tests := []struct {
		name      string
		method    string
		path      string
		auth      bool
		userLogin string
	}{
		{
			name:      "Register с методом GET",
			method:    "GET",
			path:      "/register",
			auth:      false,
			userLogin: "",
		},
		{
			name:      "Login с методом PUT",
			method:    "PUT",
			path:      "/login",
			auth:      false,
			userLogin: "",
		},
		{
			name:      "GetBinary с методом POST",
			method:    "POST",
			path:      "/api/user/binaries/123",
			auth:      true,
			userLogin: "methoduser",
		},
		{
			name:      "UpdateBinary с методом PATCH",
			method:    "PATCH",
			path:      "/api/user/binaries",
			auth:      true,
			userLogin: "methoduser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.method == "POST" || tt.method == "PUT" || tt.method == "PATCH" {
				// Для методов с телом добавляем JSON
				jsonBody, _ := json.Marshal(testData.binary)
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			if tt.auth {
				ctx := customcontext.WithUserID(req.Context(), tt.userLogin)
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusMethodNotAllowed, w.Code,
				"Expected Method Not Allowed for %s %s", tt.method, tt.path)
		})
	}
}

// TestValidationErrors - тесты валидации входных данных
func TestValidationErrors(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	// Регистрируем пользователя
	registerTestUser(t, router, "validuser", "password123")

	tests := []struct {
		name     string
		method   string
		path     string
		body     interface{}
		expected int
		errorMsg string
	}{
		{
			name:     "Register без логина",
			method:   "POST",
			path:     "/register",
			body:     dtos.NewUser{Login: "", Password: "password"},
			expected: http.StatusBadRequest,
			errorMsg: "Login and password are required",
		},
		{
			name:     "Register без пароля",
			method:   "POST",
			path:     "/register",
			body:     dtos.NewUser{Login: "user", Password: ""},
			expected: http.StatusBadRequest,
			errorMsg: "Login and password are required",
		},
		{
			name:     "Login без логина",
			method:   "POST",
			path:     "/login",
			body:     map[string]string{"login": "", "password": "password"},
			expected: http.StatusBadRequest,
			errorMsg: "Login and password are required",
		},
		{
			name:     "CreateBinary с пустыми данными",
			method:   "POST",
			path:     "/api/user/binaries",
			body:     dtos.NewBinaryData{Data: []byte(""), NewSecureEntity: dtos.NewSecureEntity{Metadata: "test"}},
			expected: http.StatusBadRequest,
			errorMsg: "Data cannot be empty",
		},
		{
			name:   "CreateCard без номера",
			method: "POST",
			path:   "/api/user/card",
			body: dtos.NewCardInformation{
				Number:          "",
				CardHolder:      "John Doe",
				ExpirationDate:  "12/25",
				CVV:             "123",
				NewSecureEntity: dtos.NewSecureEntity{Metadata: "test"},
			},
			expected: http.StatusBadRequest,
			errorMsg: "All card fields are required",
		},
		{
			name:   "CreateCredentials без пароля",
			method: "POST",
			path:   "/api/user/credentials",
			body: dtos.NewCredentials{
				Login:           "test",
				Password:        "",
				NewSecureEntity: dtos.NewSecureEntity{Metadata: "test"},
			},
			expected: http.StatusBadRequest,
			errorMsg: "Login and password are required",
		},
		{
			name:     "CreateText с пустым текстом",
			method:   "POST",
			path:     "/api/user/texts",
			body:     dtos.NewTextData{Data: "", NewSecureEntity: dtos.NewSecureEntity{Metadata: "test"}},
			expected: http.StatusBadRequest,
			errorMsg: "Text data cannot be empty",
		},
		{
			name:   "UpdateBinary без ID",
			method: "PUT",
			path:   "/api/user/binaries",
			body: entities.BinaryData{
				Data:         []byte("data"),
				SecureEntity: entities.SecureEntity{ID: "", Metadata: "test"},
			},
			expected: http.StatusBadRequest,
			errorMsg: "ID is required",
		},
		{
			name:   "UpdateCard без ID",
			method: "PUT",
			path:   "/api/user/card",
			body: entities.CardInformation{
				Number:         "4111111111111111",
				CardHolder:     "John Doe",
				ExpirationDate: "12/25",
				CVV:            "123",
				SecureEntity:   entities.SecureEntity{ID: "", Metadata: "test"},
			},
			expected: http.StatusBadRequest,
			errorMsg: "ID is required",
		},
		{
			name:   "UpdateCredentials без ID",
			method: "PUT",
			path:   "/api/user/credentials",
			body: entities.Credentials{
				Login:        "test",
				Password:     "password",
				SecureEntity: entities.SecureEntity{ID: "", Metadata: "test"},
			},
			expected: http.StatusBadRequest,
			errorMsg: "ID is required",
		},
		{
			name:   "UpdateText без ID",
			method: "PUT",
			path:   "/api/user/texts",
			body: entities.TextData{
				Data:         "test text",
				SecureEntity: entities.SecureEntity{ID: "", Metadata: "test"},
			},
			expected: http.StatusBadRequest,
			errorMsg: "ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Для защищенных маршрутов добавляем контекст
			if strings.HasPrefix(tt.path, "/api/user") {
				ctx := customcontext.WithUserID(req.Context(), "validuser")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expected, w.Code,
				"Expected status %d, got %d for %s", tt.expected, w.Code, tt.name)

			if tt.errorMsg != "" {
				assert.Contains(t, w.Body.String(), tt.errorMsg)
			}
		})
	}
}

// TestDataSizeLimits - тесты ограничений размера данных
func TestDataSizeLimits(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	// Регистрируем пользователя
	registerTestUser(t, router, "sizeuser", "password123")

	t.Run("CreateBinary с данными больше 10MB", func(t *testing.T) {
		// Создаем данные размером 11MB
		largeData := make([]byte, 11*1024*1024) // 11MB
		for i := range largeData {
			largeData[i] = byte(i % 256)
		}

		binaryData := dtos.NewBinaryData{
			Data:            largeData,
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "large data"},
		}

		req := createTestRequest("POST", "/api/user/binaries", binaryData, true, "sizeuser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
		assert.Contains(t, w.Body.String(), "Data too large")
	})

	t.Run("CreateText с текстом больше 1MB", func(t *testing.T) {
		// Создаем текст размером 1.1MB
		largeText := strings.Repeat("A", 1024*1024+1024)

		textData := dtos.NewTextData{
			Data:            largeText,
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "large text"},
		}

		req := createTestRequest("POST", "/api/user/texts", textData, true, "sizeuser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
		assert.Contains(t, w.Body.String(), "Text too large")
	})

	t.Run("UpdateBinary с данными больше 10MB", func(t *testing.T) {
		// Сначала создаем нормальную запись
		binary := createBinary(t, router, "sizeuser", dtos.NewBinaryData{
			Data:            []byte("initial data"),
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "test"},
		})

		// Пытаемся обновить с большими данными
		largeData := make([]byte, 11*1024*1024)
		updateData := entities.BinaryData{
			Data: largeData,
			SecureEntity: entities.SecureEntity{
				ID:       binary.ID,
				Metadata: "updated",
			},
		}

		req := createTestRequest("PUT", "/api/user/binaries", updateData, true, "sizeuser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	})
}

// TestCookieHandling - тесты работы с cookies
func TestCookieHandling(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	t.Run("Установка JWT cookie при регистрации", func(t *testing.T) {
		newUser := dtos.NewUser{
			Login:    "cookieuser",
			Password: "password123",
		}

		req := createTestRequest("POST", "/register", newUser, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Проверяем, что cookie установлена
		cookies := w.Result().Cookies()
		var jwtCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "jwt_token" {
				jwtCookie = cookie
				break
			}
		}

		assert.NotNil(t, jwtCookie, "JWT cookie should be set")
		assert.NotEmpty(t, jwtCookie.Value, "JWT token should not be empty")
		assert.True(t, jwtCookie.HttpOnly, "JWT cookie should be HttpOnly")
	})

	t.Run("Установка JWT cookie при логине", func(t *testing.T) {
		// Сначала регистрируем пользователя
		registerTestUser(t, router, "loginuser", "password123")

		// Логинимся
		loginReq := map[string]string{
			"login":    "loginuser",
			"password": "password123",
		}

		req := createTestRequest("POST", "/login", loginReq, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Проверяем cookie
		cookies := w.Result().Cookies()
		var jwtCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "jwt_token" {
				jwtCookie = cookie
				break
			}
		}

		assert.NotNil(t, jwtCookie, "JWT cookie should be set after login")
	})
}

// TestEmptyArraysResponse - тесты возврата пустых массивов
func TestEmptyArraysResponse(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	// Регистрируем пользователя
	registerTestUser(t, router, "emptyuser", "password123")

	tests := []struct {
		name string
		path string
	}{
		{"GetAllBinaries с пустым результатом", "/api/user/binaries"},
		{"GetAllCards с пустым результатом", "/api/user/cards"},
		{"GetAllCredentials с пустым результатом", "/api/user/credentials"},
		{"GetAllTexts с пустым результатом", "/api/user/texts"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := createTestRequest("GET", tt.path, nil, true, "emptyuser")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			// Проверяем, что возвращается пустой массив
			var response interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// В зависимости от типа ответа проверяем структуру
			switch tt.path {
			case "/api/user/binaries":
				var arr []entities.BinaryData
				err = json.Unmarshal(w.Body.Bytes(), &arr)
				require.NoError(t, err)
				assert.NotNil(t, arr)
				assert.Equal(t, 0, len(arr))
			case "/api/user/cards":
				var arr []entities.CardInformation
				err = json.Unmarshal(w.Body.Bytes(), &arr)
				require.NoError(t, err)
				assert.NotNil(t, arr)
				assert.Equal(t, 0, len(arr))
			}
		})
	}
}

// TestUpdateScenarios - тесты сценариев обновления
func TestUpdateScenarios(t *testing.T) {
	router, _ := createTestHandlerAndRouter()
	testData := getTestData()

	// Регистрируем пользователя
	registerTestUser(t, router, "updateuser", "password123")

	t.Run("Обновление несуществующей записи", func(t *testing.T) {
		updateData := entities.BinaryData{
			Data: []byte("updated data"),
			SecureEntity: entities.SecureEntity{
				ID:       "non-existent-id",
				Metadata: "updated",
			},
		}

		req := createTestRequest("PUT", "/api/user/binaries", updateData, true, "updateuser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Сервис должен вернуть 404
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Обновление записи другого пользователя", func(t *testing.T) {
		// Создаем второго пользователя
		registerTestUser(t, router, "updateuser2", "password456")

		// user1 создает запись
		binary := createBinary(t, router, "updateuser", testData.binary)

		// user2 пытается обновить запись user1
		updateData := entities.BinaryData{
			Data: []byte("hacked data"),
			SecureEntity: entities.SecureEntity{
				ID:       binary.ID,
				Metadata: "hacked",
			},
		}

		req := createTestRequest("PUT", "/api/user/binaries", updateData, true, "updateuser2")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// user2 не должен иметь доступа к записи user1
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestEdgeCases - тесты граничных случаев
func TestEdgeCases(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	t.Run("Регистрация с максимально допустимыми данными", func(t *testing.T) {
		// Создаем длинный логин и пароль
		longLogin := strings.Repeat("a", 255)
		longPassword := strings.Repeat("b", 71) //Максимальная длина bcrypt - 72 символа

		newUser := dtos.NewUser{
			Login:    longLogin,
			Password: longPassword,
		}

		req := createTestRequest("POST", "/register", newUser, false, "")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Проверяем, что регистрация прошла успешно
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Создание бинарных данных с минимальным размером", func(t *testing.T) {
		registerTestUser(t, router, "edgeuser", "password123")

		binaryData := dtos.NewBinaryData{
			Data:            []byte("a"), // Минимальный размер
			NewSecureEntity: dtos.NewSecureEntity{Metadata: "minimal"},
		}

		req := createTestRequest("POST", "/api/user/binaries", binaryData, true, "edgeuser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

// TestResponseHeaders - тесты заголовков ответов
func TestResponseHeaders(t *testing.T) {
	router, _ := createTestHandlerAndRouter()
	testData := getTestData()

	registerTestUser(t, router, "headeruser", "password123")

	t.Run("Content-Type в ответе при создании", func(t *testing.T) {
		req := createTestRequest("POST", "/api/user/binaries", testData.binary, true, "headeruser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("Content-Type в ответе при получении", func(t *testing.T) {
		binary := createBinary(t, router, "headeruser", testData.binary)

		req := createTestRequest("GET", fmt.Sprintf("/api/user/binaries/%s", binary.ID),
			nil, true, "headeruser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	})

	t.Run("Статус 201 при создании", func(t *testing.T) {
		req := createTestRequest("POST", "/api/user/texts", testData.text, true, "headeruser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Статус 410 при удалении", func(t *testing.T) {
		text := createText(t, router, "headeruser", testData.text)

		req := createTestRequest("DELETE", fmt.Sprintf("/api/user/texts/%s", text.ID),
			nil, true, "headeruser")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusGone, w.Code)
	})
}

// TestRequestBodyClosure - тесты закрытия тела запроса
func TestRequestBodyClosure(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	t.Run("Деferred закрытие тела запроса", func(t *testing.T) {
		// Создаем запрос с телом
		newUser := dtos.NewUser{
			Login:    "closeuser",
			Password: "password123",
		}

		jsonBody, _ := json.Marshal(newUser)
		req := httptest.NewRequest("POST", "/register", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		// Здесь можно было бы проверить, что r.Body.Close() вызывается,
		// но в тестах это сложно сделать без моков
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestInvalidJSON - тесты невалидного JSON
func TestInvalidJSON(t *testing.T) {
	router, _ := createTestHandlerAndRouter()

	registerTestUser(t, router, "jsonuser", "password123")

	tests := []struct {
		name     string
		method   string
		path     string
		body     string
		auth     bool
		expected int
	}{
		{
			name:     "Invalid JSON в Register",
			method:   "POST",
			path:     "/register",
			body:     "{invalid json",
			auth:     false,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid JSON в Login",
			method:   "POST",
			path:     "/login",
			body:     "{login: test, password: test}",
			auth:     false,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Invalid JSON в CreateBinary",
			method:   "POST",
			path:     "/api/user/binaries",
			body:     "{data: not base64}",
			auth:     true,
			expected: http.StatusBadRequest,
		},
		{
			name:     "Malformed JSON в CreateCard",
			method:   "POST",
			path:     "/api/user/card",
			body:     `{"number": "1234", "cardHolder": "John",`,
			auth:     true,
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader([]byte(tt.body)))
			req.Header.Set("Content-Type", "application/json")

			if tt.auth {
				ctx := customcontext.WithUserID(req.Context(), "jsonuser")
				req = req.WithContext(ctx)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expected, w.Code)
			assert.Contains(t, w.Body.String(), "Invalid request body")
		})
	}
}
