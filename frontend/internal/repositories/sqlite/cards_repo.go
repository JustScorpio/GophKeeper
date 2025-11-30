// SQLite Репозиторий для банковских карт
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/JustScorpio/GophKeeper/frontend/internal/models/dtos"
	"github.com/JustScorpio/GophKeeper/frontend/internal/models/entities"
)

// CardsRepo - репозиторий с данными банковских карт
type CardsRepo struct {
	db *sql.DB
}

// NewPgCardsRepo - инициализация репозитория
func NewCardsRepo(db *sql.DB) (*CardsRepo, error) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cards (
			id TEXT PRIMARY KEY,
			number TEXT NOT NULL,
			card_holder TEXT NOT NULL,
			expiration_date TEXT NOT NULL,
			cvv TEXT NOT NULL,
			metadata TEXT
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create cards table: %w", err)
	}

	return &CardsRepo{db: db}, nil
}

// GetAll - получить все сущности (при наличии прав у текущего пользователя)
func (r *CardsRepo) GetAll(ctx context.Context) ([]entities.CardInformation, error) {
	rows, err := r.db.Query("SELECT id, number, card_holder, expiration_date, cvv, metadata FROM cards")
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	var cards []entities.CardInformation
	for rows.Next() {
		var card entities.CardInformation
		err := rows.Scan(&card.ID, &card.Number, &card.CardHolder, &card.ExpirationDate, &card.CVV, &card.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to scan card: %w", err)
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// Get - получить сущность по ИД (при наличии прав у текущего пользователя)
func (r *CardsRepo) Get(ctx context.Context, id string) (*entities.CardInformation, error) {
	var card entities.CardInformation
	err := r.db.QueryRow("SELECT id, number, card_holder, expiration_date, cvv, metadata FROM cards WHERE id = ?", id).Scan(&card.ID, &card.Number, &card.CardHolder, &card.ExpirationDate, &card.CVV, &card.Metadata)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	return &card, nil
}

// Create - создать сущность
func (r *CardsRepo) Create(ctx context.Context, dto *dtos.NewCardInformation) (*entities.CardInformation, error) {
	var card entities.CardInformation
	err := r.db.QueryRow(
		"INSERT INTO cards (number, card_holder, expiration_date, cvv, metadata) VALUES (?, ?, ?, ?, ?) RETURNING id, number, card_holder, expiration_date, cvv, metadata",
		dto.Number, dto.CardHolder, dto.ExpirationDate, dto.CVV, dto.Metadata,
	).Scan(&card.ID, &card.Number, &card.CardHolder, &card.ExpirationDate, &card.CVV, &card.Metadata)

	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	return &card, nil
}

// Update - изменить сущность
func (r *CardsRepo) Update(ctx context.Context, entity *entities.CardInformation) (*entities.CardInformation, error) {
	var updatedCard entities.CardInformation
	err := r.db.QueryRow("UPDATE cards SET number = ?, card_holder = ?, expiration_date = ?, cvv = ?, metadata = ? WHERE id = ? RETURNING id, number, card_holder, expiration_date, cvv, metadata", entity.Number, entity.CardHolder, entity.ExpirationDate, entity.CVV, entity.Metadata, entity.ID).Scan(&updatedCard.ID, &updatedCard.Number, &updatedCard.CardHolder, &updatedCard.ExpirationDate, &updatedCard.CVV, &updatedCard.Metadata)

	if err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}
	return &updatedCard, nil
}

// Delete - удалить сущность
func (r *CardsRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec("DELETE FROM cards WHERE id = ?", id)
	return err
}
