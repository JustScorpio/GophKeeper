package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// ApiClient - клиент для взаимодействия с апи сервера
type APIClient struct {
	baseURL    string
	httpClient *http.Client
	jar        http.CookieJar
}

// NewAPIClient - создать клиент для взаимодействия с апи сервера
func NewAPIClient(baseURL string) *APIClient {
	jar, _ := cookiejar.New(nil)
	return &APIClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Jar: jar},
		jar:        jar,
	}
}

// Register - регистрация пользователя
func (s *APIClient) Register(ctx context.Context, login, password string) error {
	reqBody := map[string]string{
		"login":    login,
		"password": password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/api/user/register", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Login - аутентификация пользователя
func (s *APIClient) Login(ctx context.Context, login, password string) error {
	reqBody := map[string]string{
		"login":    login,
		"password": password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/api/user/login", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	return nil
}

// CreateBinary - создать бинарные данные
func (c *APIClient) CreateBinary(ctx context.Context, dto *dtos.NewBinaryData) (*entities.BinaryData, error) {
	jsonData, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/user/binaries", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create binary failed with status: %d", resp.StatusCode)
	}

	var binary entities.BinaryData
	if err := json.NewDecoder(resp.Body).Decode(&binary); err != nil {
		return nil, err
	}

	return &binary, nil
}

// GetAllBinaries - получить все бинарные данные
func (c *APIClient) GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/user/binaries", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get binaries failed with status: %d", resp.StatusCode)
	}

	var binaries []entities.BinaryData
	if err := json.NewDecoder(resp.Body).Decode(&binaries); err != nil {
		return nil, err
	}

	return binaries, nil
}

// UpdateBinary - обновить бинарные данные
func (c *APIClient) UpdateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error) {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+"/api/user/binaries", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update binary failed with status: %d", resp.StatusCode)
	}

	var updatedBinary entities.BinaryData
	if err := json.NewDecoder(resp.Body).Decode(&updatedBinary); err != nil {
		return nil, err
	}

	return &updatedBinary, nil
}

// DeleteBinary - удалить бинарные данные
func (c *APIClient) DeleteBinary(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+"/api/user/binaries/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusGone {
		return fmt.Errorf("delete binary failed with status: %d", resp.StatusCode)
	}

	return nil
}

// CreateCard - создать данные карты
func (c *APIClient) CreateCard(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error) {
	jsonData, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/user/cards", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create card failed with status: %d", resp.StatusCode)
	}

	var card entities.CardInformation
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, err
	}

	return &card, nil
}

// GetAllCards - получить данные всех карт
func (c *APIClient) GetAllCards(ctx context.Context) ([]entities.CardInformation, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/user/cards", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get cards failed with status: %d", resp.StatusCode)
	}

	var cards []entities.CardInformation
	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		return nil, err
	}

	return cards, nil
}

// UpdateCard - обновить данные карты
func (c *APIClient) UpdateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+"/api/user/cards", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update card failed with status: %d", resp.StatusCode)
	}

	var updatedCard entities.CardInformation
	if err := json.NewDecoder(resp.Body).Decode(&updatedCard); err != nil {
		return nil, err
	}

	return &updatedCard, nil
}

// DeleteCard - удалить данные карты
func (c *APIClient) DeleteCard(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+"/api/user/cards/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusGone {
		return fmt.Errorf("delete card failed with status: %d", resp.StatusCode)
	}

	return nil
}

// CreateCredentials - создать учётные данные
func (c *APIClient) CreateCredentials(ctx context.Context, dto *dtos.NewCredentials) (*entities.Credentials, error) {
	jsonData, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/user/credentials", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create credentials failed with status: %d", resp.StatusCode)
	}

	var credentials entities.Credentials
	if err := json.NewDecoder(resp.Body).Decode(&credentials); err != nil {
		return nil, err
	}

	return &credentials, nil
}

// GetAllCredentials - получить все учётные данные
func (c *APIClient) GetAllCredentials(ctx context.Context) ([]entities.Credentials, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/user/credentials", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get credentials failed with status: %d", resp.StatusCode)
	}

	var credentials []entities.Credentials
	if err := json.NewDecoder(resp.Body).Decode(&credentials); err != nil {
		return nil, err
	}

	return credentials, nil
}

// UpdateCredentials - обновить учётные данные
func (c *APIClient) UpdateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error) {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+"/api/user/credentials", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update credentials failed with status: %d", resp.StatusCode)
	}

	var updatedCredentials entities.Credentials
	if err := json.NewDecoder(resp.Body).Decode(&updatedCredentials); err != nil {
		return nil, err
	}

	return &updatedCredentials, nil
}

// DeleteCredentials - удалить учётные данные
func (c *APIClient) DeleteCredentials(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+"/api/user/credentials/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusGone {
		return fmt.Errorf("delete credentials failed with status: %d", resp.StatusCode)
	}

	return nil
}

// CreateText - создать текстовые данные
func (c *APIClient) CreateText(ctx context.Context, dto *dtos.NewTextData) (*entities.TextData, error) {
	jsonData, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/user/texts", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create text failed with status: %d", resp.StatusCode)
	}

	var text entities.TextData
	if err := json.NewDecoder(resp.Body).Decode(&text); err != nil {
		return nil, err
	}

	return &text, nil
}

// GetAllTexts - получить все текстовые данные
func (c *APIClient) GetAllTexts(ctx context.Context) ([]entities.TextData, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/user/texts", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get texts failed with status: %d", resp.StatusCode)
	}

	var texts []entities.TextData
	if err := json.NewDecoder(resp.Body).Decode(&texts); err != nil {
		return nil, err
	}

	return texts, nil
}

// UpdateText - обновить текстовые данные
func (c *APIClient) UpdateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error) {
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+"/api/user/texts", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update text failed with status: %d", resp.StatusCode)
	}

	var updatedTextData entities.TextData
	if err := json.NewDecoder(resp.Body).Decode(&updatedTextData); err != nil {
		return nil, err
	}

	return &updatedTextData, nil
}

// DeleteText - удалить текстовые данные
func (c *APIClient) DeleteText(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+"/api/user/texts/"+id, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusGone {
		return fmt.Errorf("delete text failed with status: %d", resp.StatusCode)
	}

	return nil
}
