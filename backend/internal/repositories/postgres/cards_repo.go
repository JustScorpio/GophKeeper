// Репозиторий postgres
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/backend/internal/customcontext"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/backend/internal/models/entities"
	"github.com/jackc/pgx/v5"
)

// PgCardsRepo - репозиторий с данными банковских карт
type PgCardsRepo struct {
	db *pgx.Conn
}

// NewPgCardsRepo - инициализация репозитория
func NewPgCardsRepo(db *pgx.Conn) (*PgCardsRepo, error) {
	return &PgCardsRepo{db: db}, nil
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *PgCardsRepo) GetAll(ctx context.Context) ([]entities.CardInformation, error) {
	userID := customcontext.GetUserID((ctx))

	rows, err := r.db.Query(ctx, "SELECT id, number, cardholder, expirationdate, cvv, metadata, ownerid FROM Cards WHERE ownerid = $1", userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var cards []entities.CardInformation
	for rows.Next() {
		var card entities.CardInformation
		err := rows.Scan(&card.ID, &card.Number, &card.CardHolder, &card.ExpirationDate, &card.CVV, &card.Metadata, &card.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *PgCardsRepo) Get(ctx context.Context, id string) (*entities.CardInformation, error) {
	userID := customcontext.GetUserID((ctx))

	var card entities.CardInformation
	err := r.db.QueryRow(ctx, "SELECT id, number, cardholder, expirationdate, cvv, metadata, ownerid FROM Cards WHERE id = $1 AND ownerid = $2", id, userID).Scan(&card.ID, &card.Number, &card.CardHolder, &card.ExpirationDate, &card.CVV, &card.Metadata, &card.OwnerID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Запись не найдена
		}
		return nil, fmt.Errorf("failed to get entity: %w", err)
	}

	return &card, nil
}

// Create - создать сущность
func (r *PgCardsRepo) Create(ctx context.Context, card *dtos.NewCardInformation) (*entities.CardInformation, error) {
	userID := customcontext.GetUserID(ctx)

	var entity entities.CardInformation
	err := r.db.QueryRow(ctx, "INSERT INTO Cards (number, cardholder, expirationdate, cvv, metadata, ownerid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, number, cardholder, expirationdate, cvv, metadata, ownerid", card.Number, card.CardHolder, card.ExpirationDate, card.CVV, card.Metadata, userID).Scan(&entity.ID, &entity.Number, &entity.CardHolder, &entity.ExpirationDate, &entity.CVV, &entity.Metadata, &entity.OwnerID)

	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Update - изменить сущность
func (r *PgCardsRepo) Update(ctx context.Context, card *entities.CardInformation) (*entities.CardInformation, error) {
	userID := customcontext.GetUserID(ctx)

	var updatedEntity entities.CardInformation
	err := r.db.QueryRow(ctx, "UPDATE Cards SET number = $2, cardholder = $3, expirationdate = $4, cvv = $5, metadata = $6 WHERE id = $1 AND ownerid = $7 RETURNING id, number, cardholder, expirationdate, cvv, metadata, ownerid", card.ID, card.Number, card.CardHolder, card.ExpirationDate, card.CVV, card.Metadata, userID).Scan(&updatedEntity.ID, &updatedEntity.Number, &updatedEntity.CardHolder, &updatedEntity.ExpirationDate, &updatedEntity.CVV, &updatedEntity.Metadata, &updatedEntity.OwnerID)

	if err != nil {
		return nil, err
	}
	return &updatedEntity, nil
}

// Delete - удалить сущность
func (r *PgCardsRepo) Delete(ctx context.Context, id string) error {
	userID := customcontext.GetUserID((ctx))

	_, err := r.db.Exec(ctx, "DELETE FROM Cards WHERE id = $1 AND ownerid = $2", id, userID)
	return err
}
