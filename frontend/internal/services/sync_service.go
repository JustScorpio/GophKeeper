package services

import (
	"context"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
)

type SyncService struct {
	apiClient    *APIClient
	localStorage *StorageService
}

// NewSyncService - создать сервис синхронизации данных
func NewSyncService(apiClient *APIClient, localStorage *StorageService) *SyncService {
	return &SyncService{
		apiClient:    apiClient,
		localStorage: localStorage,
	}
}

// SyncOnLogin - синхронизация при логине
func (s *SyncService) SyncOnLogin(ctx context.Context) error {
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
		dto := &dtos.NewBinaryData{
			Data:            binary.Data,
			NewSecureEntity: dtos.NewSecureEntity{Metadata: binary.Metadata},
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
		dto := &dtos.NewCardInformation{
			Number:          card.Number,
			CardHolder:      card.CardHolder,
			ExpirationDate:  card.ExpirationDate,
			CVV:             card.CVV,
			NewSecureEntity: dtos.NewSecureEntity{Metadata: card.Metadata},
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
		dto := &dtos.NewCredentials{
			Login:           cred.Login,
			Password:        cred.Password,
			NewSecureEntity: dtos.NewSecureEntity{Metadata: cred.Metadata},
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
		dto := &dtos.NewTextData{
			Data:            text.Data,
			NewSecureEntity: dtos.NewSecureEntity{Metadata: text.Metadata},
		}
		if _, err := s.localStorage.CreateText(ctx, dto); err != nil {
			return err
		}
	}

	return nil
}
