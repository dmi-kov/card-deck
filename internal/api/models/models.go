package models

// CreateNewDeckRequest represents type for
// request body on creating new Deck
type CreateNewDeckRequest struct {
	IsShuffled bool     `json:"is_shuffled,omitempty"`
	Cards      []string `json:"cards,omitempty"`
}

// DrawCardsRequest represents type for
// request body on draw card(s)
type DrawCardsRequest struct {
	Count uint `json:"count,omitempty"`
}
