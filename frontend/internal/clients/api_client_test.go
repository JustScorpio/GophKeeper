// clients_test - тесты для клиента
package clients_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JustScorpio/GophKeeper/frontend/internal/clients"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testServer - создает тестовый HTTP сервер с заданными обработчиками
func testServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
}

// TestAPIClient_Register - тесты регистрации
func TestAPIClient_Register(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		login          string
		password       string
		serverStatus   int
		serverResponse string
		wantErr        bool
		errorContains  string
	}{
		{
			name:         "Successful registration",
			login:        "testuser",
			password:     "testpass",
			serverStatus: http.StatusCreated,
			wantErr:      false,
		},
		{
			name:         "Successful registration with OK status",
			login:        "testuser",
			password:     "testpass",
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "Failed registration - conflict",
			login:          "existinguser",
			password:       "testpass",
			serverStatus:   http.StatusConflict,
			serverResponse: "user already exists",
			wantErr:        true,
			errorContains:  "registration failed with status: 409",
		},
		{
			name:           "Failed registration - bad request",
			login:          "",
			password:       "",
			serverStatus:   http.StatusBadRequest,
			serverResponse: "invalid credentials",
			wantErr:        true,
			errorContains:  "registration failed with status: 400",
		},
		{
			name:           "Failed registration - server error",
			login:          "testuser",
			password:       "testpass",
			serverStatus:   http.StatusInternalServerError,
			serverResponse: "internal server error",
			wantErr:        true,
			errorContains:  "registration failed with status: 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый сервер
			server := testServer(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/user/register", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var body map[string]string
				err := json.NewDecoder(r.Body).Decode(&body)
				assert.NoError(t, err)
				assert.Equal(t, tt.login, body["login"])
				assert.Equal(t, tt.password, body["password"])

				w.WriteHeader(tt.serverStatus)
				if tt.serverResponse != "" {
					fmt.Fprint(w, tt.serverResponse)
				}
			})
			defer server.Close()

			// Создаем клиент
			client := clients.NewAPIClient(server.URL)

			// Выполняем тест
			err := client.Register(ctx, tt.login, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestAPIClient_Login - тесты аутентификации
func TestAPIClient_Login(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name           string
		login          string
		password       string
		serverStatus   int
		serverResponse string
		wantErr        bool
		errorContains  string
	}{
		{
			name:         "Successful login",
			login:        "testuser",
			password:     "testpass",
			serverStatus: http.StatusOK,
			wantErr:      false,
		},
		{
			name:           "Failed login - unauthorized",
			login:          "wronguser",
			password:       "wrongpass",
			serverStatus:   http.StatusUnauthorized,
			serverResponse: "invalid credentials",
			wantErr:        true,
			errorContains:  "login failed with status: 401",
		},
		{
			name:           "Failed login - not found",
			login:          "nonexistent",
			password:       "pass",
			serverStatus:   http.StatusNotFound,
			serverResponse: "user not found",
			wantErr:        true,
			errorContains:  "login failed with status: 404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testServer(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/api/user/login", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				var body map[string]string
				err := json.NewDecoder(r.Body).Decode(&body)
				assert.NoError(t, err)
				assert.Equal(t, tt.login, body["login"])
				assert.Equal(t, tt.password, body["password"])

				w.WriteHeader(tt.serverStatus)
				if tt.serverResponse != "" {
					fmt.Fprint(w, tt.serverResponse)
				}

				// Устанавливаем cookie для успешного логина
				if tt.serverStatus == http.StatusOK {
					http.SetCookie(w, &http.Cookie{
						Name:  "session_token",
						Value: "test-session-token",
					})
				}
			})
			defer server.Close()

			client := clients.NewAPIClient(server.URL)

			err := client.Login(ctx, tt.login, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestAPIClient_CreateBinary - тест создания бинарных данных
func TestAPIClient_CreateBinary(t *testing.T) {
	ctx := context.Background()

	testData := &dtos.NewBinaryData{
		Data:            []byte("test binary data"),
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "Test binary"},
	}

	expectedResponse := &entities.BinaryData{
		SecureEntity: entities.SecureEntity{
			ID:       "test-id-123",
			Metadata: "Test binary",
		},
		Data: []byte("test binary data"),
	}

	t.Run("Successful creation", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/binaries", r.URL.Path)
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var dto dtos.NewBinaryData
			err := json.NewDecoder(r.Body).Decode(&dto)
			assert.NoError(t, err)
			assert.Equal(t, testData.Data, dto.Data)
			assert.Equal(t, testData.Metadata, dto.Metadata)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(expectedResponse)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.CreateBinary(ctx, testData)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, expectedResponse.Data, result.Data)
		assert.Equal(t, expectedResponse.Metadata, result.Metadata)
	})

	t.Run("Failed creation - bad request", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid data")
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.CreateBinary(ctx, testData)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create binary failed with status: 400")
	})

	t.Run("Failed creation - server error", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.CreateBinary(ctx, testData)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "create binary failed with status: 500")
	})

	t.Run("Invalid JSON response", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "invalid json {")
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.CreateBinary(ctx, testData)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid character")
	})
}

// TestAPIClient_GetAllBinaries - тест получения бинарных данных
func TestAPIClient_GetAllBinaries(t *testing.T) {
	ctx := context.Background()

	expectedBinaries := []entities.BinaryData{
		{
			SecureEntity: entities.SecureEntity{
				ID:       "id1",
				Metadata: "First binary",
			},
			Data: []byte("data1"),
		},
		{
			SecureEntity: entities.SecureEntity{
				ID:       "id2",
				Metadata: "Second binary",
			},
			Data: []byte("data2"),
		},
	}

	t.Run("Successful retrieval", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/binaries", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expectedBinaries)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllBinaries(ctx)
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, expectedBinaries[0].ID, result[0].ID)
		assert.Equal(t, expectedBinaries[1].ID, result[1].ID)
	})

	t.Run("Empty list", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]entities.BinaryData{})
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllBinaries(ctx)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllBinaries(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get binaries failed with status: 401")
	})
}

// TestAPIClient_UpdateBinary - тест обновления бинарных данных
func TestAPIClient_UpdateBinary(t *testing.T) {
	ctx := context.Background()

	entity := &entities.BinaryData{
		SecureEntity: entities.SecureEntity{
			ID:       "test-id",
			Metadata: "Updated metadata",
		},
		Data: []byte("updated data"),
	}

	t.Run("Successful update", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/binaries", r.URL.Path)
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var received entities.BinaryData
			err := json.NewDecoder(r.Body).Decode(&received)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, received.ID)
			assert.Equal(t, entity.Data, received.Data)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(entity)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.UpdateBinary(ctx, entity)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, entity.ID, result.ID)
	})

	t.Run("Update non-existent", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.UpdateBinary(ctx, entity)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "update binary failed with status: 404")
	})
}

// TestAPIClient_DeleteBinary - тест удаления бинарных данных
func TestAPIClient_DeleteBinary(t *testing.T) {
	ctx := context.Background()
	id := "test-id-123"

	t.Run("Successful deletion", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/binaries/"+id, r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)

			w.WriteHeader(http.StatusGone)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		err := client.DeleteBinary(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Delete non-existent", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		err := client.DeleteBinary(ctx, id)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete binary failed with status: 404")
	})
}

// TestAPIClient_CreateCard - тест создания карты
func TestAPIClient_CreateCard(t *testing.T) {
	ctx := context.Background()

	testData := &dtos.NewCardInformation{
		Number:          "4111111111111111",
		CardHolder:      "John Doe",
		ExpirationDate:  "12/25",
		CVV:             "123",
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "Test card"},
	}

	expectedResponse := &entities.CardInformation{
		SecureEntity: entities.SecureEntity{
			ID:       "card-id-123",
			Metadata: "Test card",
		},
		Number:         "4111111111111111",
		CardHolder:     "John Doe",
		ExpirationDate: "12/25",
		CVV:            "123",
	}

	t.Run("Successful creation", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/cards", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(expectedResponse)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.CreateCard(ctx, testData)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, expectedResponse.Number, result.Number)
	})
}

// TestAPIClient_GetAllCards - тест получения карт
func TestAPIClient_GetAllCards(t *testing.T) {
	ctx := context.Background()

	expectedCards := []entities.CardInformation{
		{
			SecureEntity: entities.SecureEntity{
				ID:       "card1",
				Metadata: "Visa",
			},
			Number:         "4111111111111111",
			CardHolder:     "John Doe",
			ExpirationDate: "12/25",
			CVV:            "123",
		},
		{
			SecureEntity: entities.SecureEntity{
				ID:       "card2",
				Metadata: "MasterCard",
			},
			Number:         "5555555555554444",
			CardHolder:     "Jane Smith",
			ExpirationDate: "11/26",
			CVV:            "456",
		},
	}

	t.Run("Successful retrieval", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expectedCards)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllCards(ctx)
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, expectedCards[0].Number, result[0].Number)
		assert.Equal(t, expectedCards[1].Number, result[1].Number)
	})

	t.Run("Empty list", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]entities.Credentials{})
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllCards(ctx)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllCards(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get cards failed with status: 401")
	})
}

// TestAPIClient_UpdateCard - тест обновления карты
func TestAPIClient_UpdateCard(t *testing.T) {
	ctx := context.Background()

	entity := &entities.CardInformation{
		SecureEntity: entities.SecureEntity{
			ID:       "card-id-123",
			Metadata: "Updated Visa",
		},
		Number:         "5555555555554444",
		CardHolder:     "Jane Doe",
		ExpirationDate: "12/27",
		CVV:            "456",
	}

	t.Run("Successful update", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/cards", r.URL.Path)
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var received entities.CardInformation
			err := json.NewDecoder(r.Body).Decode(&received)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, received.ID)
			assert.Equal(t, entity.Number, received.Number)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(entity)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.UpdateCard(ctx, entity)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, entity.ID, result.ID)
		assert.Equal(t, entity.Number, result.Number)
	})

	t.Run("Update non-existent card", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.UpdateCard(ctx, entity)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "update card failed with status: 404")
	})
}

// TestAPIClient_DeleteCard - тест удаления карты
func TestAPIClient_DeleteCard(t *testing.T) {
	ctx := context.Background()
	id := "card-id-123"

	t.Run("Successful deletion", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/cards/"+id, r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)

			w.WriteHeader(http.StatusGone)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		err := client.DeleteCard(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Delete non-existent card", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		err := client.DeleteCard(ctx, "non-existent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete card failed with status: 404")
	})
}

// TestAPIClient_CreateCredentials - тест создания карт
func TestAPIClient_CreateCredentials(t *testing.T) {
	ctx := context.Background()

	testData := &dtos.NewCredentials{
		Login:           "testuser",
		Password:        "testpass",
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "Test credentials"},
	}

	expectedResponse := &entities.Credentials{
		SecureEntity: entities.SecureEntity{
			ID:       "cred-id-123",
			Metadata: "Test credentials",
		},
		Login:    "testuser",
		Password: "testpass",
	}

	t.Run("Successful creation", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/credentials", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(expectedResponse)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.CreateCredentials(ctx, testData)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, expectedResponse.Login, result.Login)
	})
}

// TestAPIClient_GetAllCredentials - тест получения всех учётных данных
func TestAPIClient_GetAllCredentials(t *testing.T) {
	ctx := context.Background()

	expectedCredentials := []entities.Credentials{
		{
			SecureEntity: entities.SecureEntity{
				ID:       "cred1",
				Metadata: "Work Email",
			},
			Login:    "user1@test.com",
			Password: "password1",
		},
		{
			SecureEntity: entities.SecureEntity{
				ID:       "cred2",
				Metadata: "Personal Account",
			},
			Login:    "user2@test.com",
			Password: "password2",
		},
	}

	t.Run("Successful retrieval", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/credentials", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expectedCredentials)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllCredentials(ctx)
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, expectedCredentials[0].Login, result[0].Login)
		assert.Equal(t, expectedCredentials[1].Login, result[1].Login)
	})

	t.Run("Empty list", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]entities.Credentials{})
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllCredentials(ctx)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllCredentials(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get credentials failed with status: 401")
	})
}

// TestAPIClient_UpdateCredentials - тест обновления учетных данных
func TestAPIClient_UpdateCredentials(t *testing.T) {
	ctx := context.Background()

	entity := &entities.Credentials{
		SecureEntity: entities.SecureEntity{
			ID:       "cred-id-123",
			Metadata: "Updated Account",
		},
		Login:    "updated@example.com",
		Password: "newpassword456",
	}

	t.Run("Successful update", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/credentials", r.URL.Path)
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var received entities.Credentials
			err := json.NewDecoder(r.Body).Decode(&received)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, received.ID)
			assert.Equal(t, entity.Login, received.Login)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(entity)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.UpdateCredentials(ctx, entity)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, entity.ID, result.ID)
		assert.Equal(t, entity.Login, result.Login)
	})

	t.Run("Update non-existent credentials", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.UpdateCredentials(ctx, entity)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "update credentials failed with status: 404")
	})
}

// TestAPIClient_DeleteCredentials - тест удаления учетных данных
func TestAPIClient_DeleteCredentials(t *testing.T) {
	ctx := context.Background()
	id := "cred-id-123"

	t.Run("Successful deletion", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/credentials/"+id, r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)

			w.WriteHeader(http.StatusGone)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		err := client.DeleteCredentials(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Delete non-existent credentials", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		err := client.DeleteCredentials(ctx, "non-existent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete credentials failed with status: 404")
	})
}

func TestAPIClient_CreateText(t *testing.T) {
	ctx := context.Background()

	testData := &dtos.NewTextData{
		Data:            "This is test text data",
		NewSecureEntity: dtos.NewSecureEntity{Metadata: "Test text"},
	}

	expectedResponse := &entities.TextData{
		SecureEntity: entities.SecureEntity{
			ID:       "text-id-123",
			Metadata: "Test text",
		},
		Data: "This is test text data",
	}

	t.Run("Successful creation", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/texts", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(expectedResponse)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.CreateText(ctx, testData)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, expectedResponse.ID, result.ID)
		assert.Equal(t, expectedResponse.Data, result.Data)
	})
}

// TestAPIClient_GetAllTexts - тест получения всех текстовых данных
func TestAPIClient_GetAllTexts(t *testing.T) {
	ctx := context.Background()

	expectedTexts := []entities.TextData{
		{
			SecureEntity: entities.SecureEntity{
				ID:       "text1",
				Metadata: "First note",
			},
			Data: "This is first text",
		},
		{
			SecureEntity: entities.SecureEntity{
				ID:       "text2",
				Metadata: "Second note",
			},
			Data: "This is second text",
		},
	}

	t.Run("Successful retrieval", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/texts", r.URL.Path)
			assert.Equal(t, "GET", r.Method)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(expectedTexts)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllTexts(ctx)
		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.Equal(t, expectedTexts[0].Data, result[0].Data)
		assert.Equal(t, expectedTexts[1].Data, result[1].Data)
	})

	t.Run("Empty list", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]entities.TextData{})
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllTexts(ctx)
		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.GetAllTexts(ctx)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "get texts failed with status: 401")
	})
}

// TestAPIClient_UpdateText - тест обновления текстовых данных
func TestAPIClient_UpdateText(t *testing.T) {
	ctx := context.Background()

	entity := &entities.TextData{
		SecureEntity: entities.SecureEntity{
			ID:       "text-id-123",
			Metadata: "Updated Note",
		},
		Data: "This is updated confidential information",
	}

	t.Run("Successful update", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/texts", r.URL.Path)
			assert.Equal(t, "PUT", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var received entities.TextData
			err := json.NewDecoder(r.Body).Decode(&received)
			assert.NoError(t, err)
			assert.Equal(t, entity.ID, received.ID)
			assert.Equal(t, entity.Data, received.Data)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(entity)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.UpdateText(ctx, entity)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, entity.ID, result.ID)
		assert.Equal(t, entity.Data, result.Data)
	})

	t.Run("Update non-existent text", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		result, err := client.UpdateText(ctx, entity)
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "update text failed with status: 404")
	})
}

// TestAPIClient_DeleteText - тест удаления текстовых данных
func TestAPIClient_DeleteText(t *testing.T) {
	ctx := context.Background()
	id := "text-id-123"

	t.Run("Successful deletion", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/api/user/texts/"+id, r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)

			w.WriteHeader(http.StatusGone)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		err := client.DeleteText(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Delete non-existent text", func(t *testing.T) {
		server := testServer(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
		defer server.Close()

		client := clients.NewAPIClient(server.URL)

		err := client.DeleteText(ctx, "non-existent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "delete text failed with status: 404")
	})
}

// TestAPIClient_CookieJar - тест кук
func TestAPIClient_CookieJar(t *testing.T) {
	ctx := context.Background()

	// Тестируем сохранение cookies между запросами
	loginCookie := ""

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/user/login":
			// Устанавливаем cookie при логине
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: "test-session-value",
			})
			w.WriteHeader(http.StatusOK)
		case "/api/user/binaries":
			// Проверяем, что cookie отправляется
			cookie, err := r.Cookie("session")
			require.NoError(t, err)
			loginCookie = cookie.Value

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]entities.BinaryData{})
		}
	}))
	defer server.Close()

	client := clients.NewAPIClient(server.URL)

	// Логинимся
	err := client.Login(ctx, "test", "pass")
	require.NoError(t, err)

	// Делаем другой запрос
	_, err = client.GetAllBinaries(ctx)
	require.NoError(t, err)

	// Проверяем, что cookie был отправлен
	assert.Equal(t, "test-session-value", loginCookie)
}

// TestAPIClient_InvalidURL - тест с невалидным адресом
func TestAPIClient_InvalidURL(t *testing.T) {
	ctx := context.Background()

	t.Run("Malformed base URL", func(t *testing.T) {
		client := clients.NewAPIClient("://invalid-url")

		err := client.Register(ctx, "test", "pass")
		require.Error(t, err)
	})

	t.Run("Empty base URL", func(t *testing.T) {
		client := clients.NewAPIClient("")

		err := client.Register(ctx, "test", "pass")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported protocol scheme")
	})
}

// TestAPIClient_ConcurrentRequests - тест конкурентных запросов
func TestAPIClient_ConcurrentRequests(t *testing.T) {
	ctx := context.Background()

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]entities.BinaryData{})
	}))
	defer server.Close()

	client := clients.NewAPIClient(server.URL)

	// Запускаем несколько concurrent запросов
	concurrency := 5
	errors := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			_, err := client.GetAllBinaries(ctx)
			errors <- err
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < concurrency; i++ {
		err := <-errors
		assert.NoError(t, err)
	}

	assert.Equal(t, concurrency, requestCount)
}
