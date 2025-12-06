package clients

import (
	"context"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// IAPIClient - интерфейс для API клиента
type IAPIClient interface {
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) error

	// Binary methods
	CreateBinary(ctx context.Context, dto *dtos.NewBinaryData) (*entities.BinaryData, error)
	GetAllBinaries(ctx context.Context) ([]entities.BinaryData, error)
	UpdateBinary(ctx context.Context, entity *entities.BinaryData) (*entities.BinaryData, error)
	DeleteBinary(ctx context.Context, id string) error

	// Card methods
	CreateCard(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error)
	GetAllCards(ctx context.Context) ([]entities.CardInformation, error)
	UpdateCard(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error)
	DeleteCard(ctx context.Context, id string) error

	// Credentials methods
	CreateCredentials(ctx context.Context, dto *dtos.NewCredentials) (*entities.Credentials, error)
	GetAllCredentials(ctx context.Context) ([]entities.Credentials, error)
	UpdateCredentials(ctx context.Context, entity *entities.Credentials) (*entities.Credentials, error)
	DeleteCredentials(ctx context.Context, id string) error

	// Text methods
	CreateText(ctx context.Context, dto *dtos.NewTextData) (*entities.TextData, error)
	GetAllTexts(ctx context.Context) ([]entities.TextData, error)
	UpdateText(ctx context.Context, entity *entities.TextData) (*entities.TextData, error)
	DeleteText(ctx context.Context, id string) error
}
