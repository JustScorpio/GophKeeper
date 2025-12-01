// Сервис для работы с аутентификацией
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

type AuthService struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthService(baseURL string) *AuthService {
	return &AuthService{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// Register - регистрация пользователя
func (s *AuthService) Register(ctx context.Context, login, password string) (*entities.User, error) {
	reqBody := map[string]string{
		"login":    login,
		"password": password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/api/user/register", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	var user entities.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Login - аутентификация пользователя
func (s *AuthService) Login(ctx context.Context, login, password string) (*entities.User, error) {
	reqBody := map[string]string{
		"login":    login,
		"password": password,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/api/user/login", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var user entities.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
