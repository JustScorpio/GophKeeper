// services - сервисы
package services

import (
	"context"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/clients"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// SyncService - сервис синхронизации данных
type SyncService struct {
	apiClient    clients.IAPIClient
	localStorage *StorageService
}

// NewSyncService - создать сервис синхронизации данных
func NewSyncService(apiClient clients.IAPIClient, localStorage *StorageService) *SyncService {
	return &SyncService{
		apiClient:    apiClient,
		localStorage: localStorage,
	}
}

// Sync - синхронизация данных с сервером
func (s *SyncService) Sync(ctx context.Context) error {
	if err := s.syncBinaries(ctx); err != nil {
		return fmt.Errorf("failed to sync binaries: %w", err)
	}
	if err := s.syncCards(ctx); err != nil {
		return fmt.Errorf("failed to sync cards: %w", err)
	}
	if err := s.syncCredentials(ctx); err != nil {
		return fmt.Errorf("failed to sync credentials: %w", err)
	}
	if err := s.syncTexts(ctx); err != nil {
		return fmt.Errorf("failed to sync texts: %w", err)
	}

	return nil
}

// syncBinaries - синхронизировать бинарные данные
func (s *SyncService) syncBinaries(ctx context.Context) error {
	// Данные с сервера
	serverBinaries, err := s.apiClient.GetAllBinaries(ctx)
	if err != nil {
		return fmt.Errorf("get server binaries: %w", err)
	}

	// Локальные данные
	localBinaries, err := s.localStorage.GetAllBinaries(ctx)
	if err != nil {
		return fmt.Errorf("get local binaries: %w", err)
	}

	// Создаём мапу серверных сущностей
	serverMap := make(map[string]entities.BinaryData, len(serverBinaries))
	for _, binary := range serverBinaries {
		serverMap[binary.ID] = binary
	}

	// Удаляем локальные сущности, которых нет на сервере или которые изменены
	for _, localBinary := range localBinaries {
		id := localBinary.ID
		serverBinary, existsOnServer := serverMap[id]

		if existsOnServer && entities.Equals(&localBinary, &serverBinary) {
			// Сущность не изменилась - удаляем из serverMap, чтобы не создавать заново
			delete(serverMap, id)
			continue
		}

		// Если сущности нет на сервере либо она изменилась - удаляем локально
		if err := s.localStorage.DeleteBinary(ctx, id); err != nil {
			return fmt.Errorf("delete binary %s: %w", id, err)
		}
	}

	// Создаем сущности которые есть на сервере, но нет локально
	for _, serverBinary := range serverMap {
		if _, err := s.localStorage.CreateBinary(ctx, &serverBinary); err != nil {
			return fmt.Errorf("create binary %s: %w", serverBinary.ID, err)
		}
	}

	return nil
}

// syncCards - синхронизировать данные карт
func (s *SyncService) syncCards(ctx context.Context) error {
	// Данные с сервера
	serverCards, err := s.apiClient.GetAllCards(ctx)
	if err != nil {
		return fmt.Errorf("get server cards: %w", err)
	}

	// Локальные данные
	localCards, err := s.localStorage.GetAllCards(ctx)
	if err != nil {
		return fmt.Errorf("get local cards: %w", err)
	}

	// Создаём мапу серверных сущностей
	serverMap := make(map[string]entities.CardInformation, len(serverCards))
	for _, card := range serverCards {
		serverMap[card.ID] = card
	}

	// Удаляем локальные сущности, которых нет на сервере или которые изменены
	for _, localCard := range localCards {
		id := localCard.ID
		serverCard, existsOnServer := serverMap[id]

		if existsOnServer && entities.Equals(&localCard, &serverCard) {
			// Сущность не изменилась - удаляем из serverMap, чтобы не создавать заново
			delete(serverMap, id)
			continue
		}

		// Если сущности нет на сервере либо она изменилась - удаляем локально
		if err := s.localStorage.DeleteCard(ctx, id); err != nil {
			return fmt.Errorf("delete card %s: %w", id, err)
		}
	}

	// Создаем сущности которые есть на сервере, но нет локально
	for _, serverCard := range serverMap {
		if _, err := s.localStorage.CreateCard(ctx, &serverCard); err != nil {
			return fmt.Errorf("create card %s: %w", serverCard.ID, err)
		}
	}

	return nil
}

// syncCredentials - синхронизировать учётные данные
func (s *SyncService) syncCredentials(ctx context.Context) error {
	// Данные с сервера
	serverCredentials, err := s.apiClient.GetAllCredentials(ctx)
	if err != nil {
		return fmt.Errorf("get server credentials: %w", err)
	}

	// Локальные данные
	localCredentials, err := s.localStorage.GetAllCredentials(ctx)
	if err != nil {
		return fmt.Errorf("get local credentials: %w", err)
	}

	// Создаём мапу серверных сущностей
	serverMap := make(map[string]entities.Credentials, len(serverCredentials))
	for _, cred := range serverCredentials {
		serverMap[cred.ID] = cred
	}

	// Удаляем локальные сущности, которых нет на сервере или которые изменены
	for _, localCred := range localCredentials {
		id := localCred.ID
		serverCred, existsOnServer := serverMap[id]

		if existsOnServer && entities.Equals(&localCred, &serverCred) {
			// Сущность не изменилась - удаляем из serverMap, чтобы не создавать заново
			delete(serverMap, id)
			continue
		}

		// Если сущности нет на сервере либо она изменилась - удаляем локально
		if err := s.localStorage.DeleteCredentials(ctx, id); err != nil {
			return fmt.Errorf("delete credentials %s: %w", id, err)
		}
	}

	// Создаем сущности которые есть на сервере, но нет локально
	for _, serverCred := range serverMap {
		if _, err := s.localStorage.CreateCredentials(ctx, &serverCred); err != nil {
			return fmt.Errorf("create credentials %s: %w", serverCred.ID, err)
		}
	}

	return nil
}

// syncTexts - синхронизировать текстовые данные
func (s *SyncService) syncTexts(ctx context.Context) error {
	// Данные с сервера
	serverTexts, err := s.apiClient.GetAllTexts(ctx)
	if err != nil {
		return fmt.Errorf("get server texts: %w", err)
	}

	// Локальные данные
	localTexts, err := s.localStorage.GetAllTexts(ctx)
	if err != nil {
		return fmt.Errorf("get local texts: %w", err)
	}

	// Создаём мапу серверных сущностей
	serverMap := make(map[string]entities.TextData, len(serverTexts))
	for _, text := range serverTexts {
		serverMap[text.ID] = text
	}

	// Удаляем локальные сущности, которых нет на сервере или которые изменены
	for _, localText := range localTexts {
		id := localText.ID
		serverText, existsOnServer := serverMap[id]

		if existsOnServer && entities.Equals(&localText, &serverText) {
			// Сущность не изменилась - удаляем из serverMap, чтобы не создавать заново
			delete(serverMap, id)
			continue
		}

		// Если сущности нет на сервере либо она изменилась - удаляем локально
		if err := s.localStorage.DeleteText(ctx, id); err != nil {
			return fmt.Errorf("delete text %s: %w", id, err)
		}
	}

	// Создаем сущности которые есть на сервере, но нет локально
	for _, serverText := range serverMap {
		if _, err := s.localStorage.CreateText(ctx, &serverText); err != nil {
			return fmt.Errorf("create text %s: %w", serverText.ID, err)
		}
	}

	return nil
}
