package repository

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/card-deck/internal/models"
	"github.com/card-deck/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const decksTable = "decks"

// CreateDeck creates new deck
func (r *Repository) CreateDeck(deck *models.Deck) error {
	now := time.Now().UTC()
	deck.CreatedAt = now
	deck.UpdatedAt = now
	deck.DeckID = uuid.NewV4().String()
	deck.Remaining = uint(len(deck.CardCodes))

	_, err := sb.Insert(decksTable).
		SetMap(map[string]interface{}{
			"deck_id":     deck.DeckID,
			"is_shuffled": deck.IsShuffled,
			"remaining":   deck.Remaining,
			"card_codes":  deck.CardCodes,
			"created_at":  deck.CreatedAt,
			"updated_at":  deck.UpdatedAt,
		}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

// GetDeckByID returns deck by it's ID
func (r *Repository) GetDeckByID(deckID string) (*models.Deck, error) {
	query, args, err := sb.Select(
		"deck_id",
		"updated_at",
		"is_shuffled",
		"remaining",
		"card_codes",
		"created_at",
		"updated_at",
	).From(decksTable).
		Where(sq.Eq{"deck_id": deckID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var deck models.Deck
	if err = r.db.Get(&deck, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.NotFound, "deck not found")
		}
		return nil, err
	}

	return &deck, nil
}

// UpdateDeck updates deck by it's ID
func (r *Repository) UpdateDeck(deck *models.Deck) error {
	now := time.Now().UTC()
	deck.UpdatedAt = now
	deck.Remaining = uint(len(deck.CardCodes))

	_, err := sb.Update(decksTable).
		Where(sq.Eq{
			"deck_id": deck.DeckID,
		}).
		SetMap(map[string]interface{}{
			"remaining":  deck.Remaining,
			"card_codes": deck.CardCodes,
			"updated_at": deck.UpdatedAt,
		}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}
