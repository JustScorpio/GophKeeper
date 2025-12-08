package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/customerrors"
	"github.com/JustScorpio/GophKeeper/backend/internal/middleware/auth" //В файле c middleware не только middleware, но и ауфные функции
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/JustScorpio/GophKeeper/backend/internal/services"
	"github.com/JustScorpio/GophKeeper/backend/internal/utils"
	"github.com/go-chi/chi"
)

type GophkeeperHandler struct {
	service *services.StorageService
}

func NewGophkeeperHandler(service *services.StorageService) *GophkeeperHandler {
	return &GophkeeperHandler{
		service: service,
	}
}

// Регистрация пользователя
func (h *GophkeeperHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Только Content-Type: JSON
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req dtos.NewUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Валидация
	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	// Хэшируем пароль
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	req.Password = hashedPassword

	// Создаём пользователя
	user, err := h.service.CreateUser(r.Context(), req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	// Устанавливаем JWT с логином
	if err := auth.SetJWTCookie(w, user.Login); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Аутентификация пользователя
func (h *GophkeeperHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Только Content-Type: JSON
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	//Валидация
	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	// Проверяем пользователя
	user, err := h.service.GetUser(r.Context(), req.Login)
	if err != nil || user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Устанавливаем JWT токен с логином
	if err := auth.SetJWTCookie(w, user.Login); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CreateBinary - создать бинарные данные
func (h *GophkeeperHandler) CreateBinary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Только Content-Type: JSON
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req dtos.NewBinaryData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация
	if len(req.Data) == 0 {
		http.Error(w, "Data cannot be empty", http.StatusBadRequest)
		return
	}

	if len(req.Data) > 10*1024*1024 { // 10MB limit
		http.Error(w, "Data too large", http.StatusRequestEntityTooLarge)
		return
	}

	binary, err := h.service.CreateBinary(r.Context(), &req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(binary)
}

// GetBinary - получить бинарные данные
func (h *GophkeeperHandler) GetBinary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	binary, err := h.service.GetBinary(r.Context(), id)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if binary == nil {
		http.Error(w, "Binary data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(binary)
}

// GetAllBinaries - получить все бинарные данные пользователя
func (h *GophkeeperHandler) GetAllBinaries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	binaries, err := h.service.GetAllBinaries(r.Context())
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	// Всегда возвращаем массив, даже если он пустой
	if binaries == nil {
		binaries = []entities.BinaryData{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(binaries)
}

// UpdateBinary - обновить бинарные данные
func (h *GophkeeperHandler) UpdateBinary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//Только Content-Type: JSON
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req entities.BinaryData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация
	if req.ID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if len(req.Data) == 0 {
		http.Error(w, "Data cannot be empty", http.StatusBadRequest)
		return
	}

	if len(req.Data) > 10*1024*1024 { // 10MB limit
		http.Error(w, "Data too large", http.StatusRequestEntityTooLarge)
		return
	}

	updatedBinary, err := h.service.UpdateBinary(r.Context(), &req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBinary)
}

// DeleteBinary - удалить бинарные данные
func (h *GophkeeperHandler) DeleteBinary(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteBinary(r.Context(), id)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusGone)
}

// CreateCard - создать данные банковской карты
func (h *GophkeeperHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req dtos.NewCardInformation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация
	if req.Number == "" || req.CardHolder == "" || req.ExpirationDate == "" || req.CVV == "" {
		http.Error(w, "All card fields are required", http.StatusBadRequest)
		return
	}

	card, err := h.service.CreateCard(r.Context(), &req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

// GetCard - получить данные банковской карты
func (h *GophkeeperHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	card, err := h.service.GetCard(r.Context(), id)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if card == nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(card)
}

// GetAllCards - получить все данные банковских карт пользователя
func (h *GophkeeperHandler) GetAllCards(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	cards, err := h.service.GetAllCards(r.Context())
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if cards == nil {
		cards = []entities.CardInformation{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}

// UpdateCard - обновить данные банковской карты
func (h *GophkeeperHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req entities.CardInformation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация
	if req.ID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if req.Number == "" || req.CardHolder == "" || req.ExpirationDate == "" || req.CVV == "" {
		http.Error(w, "All card fields are required", http.StatusBadRequest)
		return
	}

	updatedCard, err := h.service.UpdateCard(r.Context(), &req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCard)
}

// DeleteCard - удалить данные банковской карты
func (h *GophkeeperHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteCard(r.Context(), id)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusGone)
}

// CreateCredentials - создать учётные данные
func (h *GophkeeperHandler) CreateCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req dtos.NewCredentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация
	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	credentials, err := h.service.CreateCredentials(r.Context(), &req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credentials)
}

// GetCredentials - получить учётные данные
func (h *GophkeeperHandler) GetCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	credentials, err := h.service.GetCredentials(r.Context(), id)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if credentials == nil {
		http.Error(w, "Credentials not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(credentials)
}

// GetAllCredentials - получить все учётные данные пользователя
func (h *GophkeeperHandler) GetAllCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	credentials, err := h.service.GetAllCredentials(r.Context())
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if credentials == nil {
		credentials = []entities.Credentials{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(credentials)
}

// UpdateCredentials - обновить учётные данные
func (h *GophkeeperHandler) UpdateCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req entities.Credentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация
	if req.ID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	updatedCredentials, err := h.service.UpdateCredentials(r.Context(), &req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedCredentials)
}

// DeleteCredentials - удалить учётные данные
func (h *GophkeeperHandler) DeleteCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteCredentials(r.Context(), id)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusGone)
}

// CreateText - создать текстовые данные
func (h *GophkeeperHandler) CreateText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req dtos.NewTextData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация
	if req.Data == "" {
		http.Error(w, "Text data cannot be empty", http.StatusBadRequest)
		return
	}

	if len(req.Data) > 1*1024*1024 { // 1MB limit для текста
		http.Error(w, "Text too large", http.StatusRequestEntityTooLarge)
		return
	}

	text, err := h.service.CreateText(r.Context(), &req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(text)
}

// GetText - получить текстовые данные
func (h *GophkeeperHandler) GetText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	text, err := h.service.GetText(r.Context(), id)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if text == nil {
		http.Error(w, "Text data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(text)
}

// GetAllTexts - получить все текстовые данные пользователя
func (h *GophkeeperHandler) GetAllTexts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	texts, err := h.service.GetAllTexts(r.Context())
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	if texts == nil {
		texts = []entities.TextData{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(texts)
}

// UpdateText - обновить текстовые данные
func (h *GophkeeperHandler) UpdateText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req entities.TextData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Валидация
	if req.ID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	if req.Data == "" {
		http.Error(w, "Text data cannot be empty", http.StatusBadRequest)
		return
	}

	if len(req.Data) > 1*1024*1024 { // 1MB limit для текста
		http.Error(w, "Text too large", http.StatusRequestEntityTooLarge)
		return
	}

	updatedText, err := h.service.UpdateText(r.Context(), &req)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedText)
}

// DeleteText - удалить текстовые данные
func (h *GophkeeperHandler) DeleteText(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	login := customcontext.GetUserID(r.Context())
	if login == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteText(r.Context(), id)
	if err != nil {
		var statusCode = http.StatusInternalServerError

		var httpErr *customerrors.HTTPError
		if errors.As(err, &httpErr) {
			statusCode = httpErr.Code
		}

		http.Error(w, err.Error(), statusCode)
		return
	}

	w.WriteHeader(http.StatusGone)
}
