package services

import (
	"context"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/clients"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

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
	serverBinaries, err := s.apiClient.GetAllBinaries(ctx)
	if err != nil {
		return err
	}

	localBinaries, err := s.localStorage.GetAllBinaries(ctx)
	if err != nil {
		return err
	}

	for _, local := range localBinaries {
		if err := s.localStorage.DeleteBinary(ctx, local.ID); err != nil {
			return err
		}
	}

	for _, binary := range serverBinaries {
		dto := &entities.BinaryData{
			Data:         binary.Data,
			SecureEntity: entities.SecureEntity{ID: binary.ID, Metadata: binary.Metadata},
		}
		if _, err := s.localStorage.CreateBinary(ctx, dto); err != nil {
			return err
		}
	}

	return nil
}

// syncCards - синхронизировать данные карт
func (s *SyncService) syncCards(ctx context.Context) error {
	serverCards, err := s.apiClient.GetAllCards(ctx)
	if err != nil {
		return err
	}

	localCards, err := s.localStorage.GetAllCards(ctx)
	if err != nil {
		return err
	}

	for _, local := range localCards {
		if err := s.localStorage.DeleteCard(ctx, local.ID); err != nil {
			return err
		}
	}

	for _, card := range serverCards {
		dto := &entities.CardInformation{
			Number:         card.Number,
			CardHolder:     card.CardHolder,
			ExpirationDate: card.ExpirationDate,
			CVV:            card.CVV,
			SecureEntity:   entities.SecureEntity{ID: card.ID, Metadata: card.Metadata},
		}
		if _, err := s.localStorage.CreateCard(ctx, dto); err != nil {
			return err
		}
	}

	return nil
}

// syncCredentials - синхронизировать учётные данные
func (s *SyncService) syncCredentials(ctx context.Context) error {
	serverCredentials, err := s.apiClient.GetAllCredentials(ctx)
	if err != nil {
		return err
	}

	localCredentials, err := s.localStorage.GetAllCredentials(ctx)
	if err != nil {
		return err
	}

	for _, local := range localCredentials {
		if err := s.localStorage.DeleteCredentials(ctx, local.ID); err != nil {
			return err
		}
	}

	for _, cred := range serverCredentials {
		dto := &entities.Credentials{
			Login:        cred.Login,
			Password:     cred.Password,
			SecureEntity: entities.SecureEntity{ID: cred.ID, Metadata: cred.Metadata},
		}
		if _, err := s.localStorage.CreateCredentials(ctx, dto); err != nil {
			return err
		}
	}

	return nil
}

// syncTexts - синхронизировать текстовые данные
func (s *SyncService) syncTexts(ctx context.Context) error {
	serverTexts, err := s.apiClient.GetAllTexts(ctx)
	if err != nil {
		return err
	}

	localTexts, err := s.localStorage.GetAllTexts(ctx)
	if err != nil {
		return err
	}

	for _, local := range localTexts {
		if err := s.localStorage.DeleteText(ctx, local.ID); err != nil {
			return err
		}
	}

	for _, text := range serverTexts {
		dto := &entities.TextData{
			Data:         text.Data,
			SecureEntity: entities.SecureEntity{ID: text.ID, Metadata: text.Metadata},
		}
		if _, err := s.localStorage.CreateText(ctx, dto); err != nil {
			return err
		}
	}

	return nil
}
