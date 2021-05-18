package models

import (
	"fmt"
	"math/rand"
	"time"

	deckhelper "github.com/card-deck/pkg/deck"

	"github.com/lib/pq"
)

type (
	// Deck is a type that represents
	// the model of the decks table.
	Deck struct {
		DeckID     string         `json:"deck_id" db:"deck_id"`
		IsShuffled bool           `json:"is_shuffled" db:"is_shuffled"`
		Remaining  uint           `json:"remaining" db:"remaining"`
		Cards      Cards          `json:"cards,omitempty"`
		CardCodes  pq.StringArray `json:"-" db:"card_codes"`
		CreatedAt  time.Time      `json:"created_at" db:"created_at"`
		UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	}

	// Card is a type that represents card object
	Card struct {
		Value string `json:"value"`
		Suit  string `json:"suit"`
		Code  string `json:"code"`
	}

	// Cards is a type that represents
	// list of cards
	Cards []*Card
)

// BuildCardsFromCodes build card objects from card codes
func BuildCardsFromCodes(codes []string) (Cards, error) {
	if !deckhelper.IsValidCodes(codes) {
		return nil, fmt.Errorf("codes is not valid")
	}
	if len(codes) == 0 {
		return Cards{}, nil
	}
	cards := make(Cards, 0, len(codes))
	for _, code := range codes {
		var (
			cardCode  string
			suiteCode string
		)

		// to get for example 10S card
		// TODO: refactor to make more understandable??
		if len(code) == 3 {
			cardCode = code[0:2]
			suiteCode = code[2:3]
		} else {
			cardCode = code[0:1]
			suiteCode = code[1:2]
		}

		card := &Card{
			Value: deckhelper.CardsMapping[cardCode],
			Suit:  deckhelper.SuitsMapping[suiteCode],
			Code:  code,
		}

		cards = append(cards, card)
	}
	return cards, nil
}

// DrawRandomNCars returns N random codes from slice
func DrawRandomNCars(count uint, codes []string) ([]string, error) {
	if int(count) > len(codes) {
		return nil, fmt.Errorf("count cannot be greater then length of codes slice")
	}

	rand.Seed(time.Now().UnixNano())
	randomize := rand.Perm(len(codes))
	var result []string
	for _, v := range randomize[:count] {
		result = append(result, codes[v])
	}

	return result, nil
}

// RemoveDrawnCodes removes drawn codes from slice keeping order
func RemoveDrawnCodes(drawnCodes, allCodes []string) []string {
	for i := 0; i < len(allCodes); i++ {
		code := allCodes[i]
		for _, rem := range drawnCodes {
			if code == rem {
				allCodes = append(allCodes[:i], allCodes[i+1:]...)
				i--
				break
			}
		}
	}
	return allCodes
}
