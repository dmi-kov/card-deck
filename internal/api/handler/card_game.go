package handler

import (
	"net/http"

	apimodels "github.com/card-deck/internal/api/models"
	"github.com/card-deck/internal/models"
	deckhelper "github.com/card-deck/pkg/deck"
	"github.com/card-deck/pkg/errors"
	httphelper "github.com/card-deck/pkg/http"

	"github.com/card-deck/internal/repository"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// CardGameHandler represents type to handle income HTTP requests for card game
type CardGameHandler struct {
	repo   *repository.Repository
	logger *zap.Logger
}

// NewCardGameHandler creates new instance of CardGameHandler
func NewCardGameHandler(repo *repository.Repository, logger *zap.Logger) *CardGameHandler {
	return &CardGameHandler{
		repo:   repo,
		logger: logger,
	}
}

// MountRoutes mounts the endpoint routes to the router instance
func (h *CardGameHandler) MountRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Method(http.MethodPost, "/v1/deck", httphelper.Handler(h.CreateDeck))
		r.Method(http.MethodGet, "/v1/deck/{deckID}", httphelper.Handler(h.OpenDeck))
		r.Method(http.MethodPatch, "/v1/deck/{deckID}/cards", httphelper.Handler(h.DrawCards))
	})
}

// CreateDeck Creates new Deck
// Route /v1/deck [post]
func (h *CardGameHandler) CreateDeck(w http.ResponseWriter, r *http.Request) error {
	var payload apimodels.CreateNewDeckRequest

	if err := httphelper.ReadJSON(r, &payload); err != nil {
		return errors.Wrap(err, errors.InvalidInput, "failed parse payload")
	}
	if !deckhelper.IsValidCodes(payload.Cards) {
		return errors.New(errors.InvalidInput, "given cards is not valid")
	}

	deck := &models.Deck{
		IsShuffled: payload.IsShuffled,
		CardCodes:  payload.Cards,
	}

	if len(deck.CardCodes) == 0 {
		deck.CardCodes = deckhelper.CreateDefaultCodes()
	}

	if deck.IsShuffled {
		deckhelper.ShuffleDeck(deck.CardCodes)
	}

	if err := h.repo.CreateDeck(deck); err != nil {
		return errors.Wrap(err, errors.Internal, "failed store new deck")
	}

	return httphelper.WriteSuccessResponse(w, http.StatusOK, deck)
}

// OpenDeck returns all cards into deck by it's ID
// Route /v1/deck/{deckID} [get]
func (h *CardGameHandler) OpenDeck(w http.ResponseWriter, r *http.Request) error {
	deckID := chi.URLParam(r, "deckID")
	if deckID == "" {
		return errors.New(errors.InvalidInput, "deckID is required")
	}

	deck, err := h.repo.GetDeckByID(deckID)
	if err != nil {
		return errors.New(errors.NotFound, "deck not found")
	}

	cards, err := models.BuildCardsFromCodes(deck.CardCodes)
	if err != nil {
		return errors.New(errors.Internal, "failed map codes to cards")
	}
	deck.Cards = cards

	return httphelper.WriteSuccessResponse(w, http.StatusOK, deck)
}

// DrawCards draws [N] cards from deck by it's ID
// Route /v1/deck/{deckID}/cards [patch]
func (h *CardGameHandler) DrawCards(w http.ResponseWriter, r *http.Request) error {
	var payload apimodels.DrawCardsRequest

	deckID := chi.URLParam(r, "deckID")
	if deckID == "" {
		return errors.New(errors.InvalidInput, "deckID is required")
	}

	if err := httphelper.ReadJSON(r, &payload); err != nil {
		return errors.Wrap(err, errors.InvalidInput, "failed parse payload")
	}

	if payload.Count > 52 || payload.Count == 0 {
		return errors.New(errors.InvalidInput, "count cannot be more than 52 or 0")
	}

	deck, err := h.repo.GetDeckByID(deckID)
	if err != nil {
		return errors.New(errors.NotFound, "deck not found")
	}

	if deck.Remaining == 0 {
		return errors.New(errors.InvalidInput, "deck remaining 0 cards")
	}
	if deck.Remaining < payload.Count {
		payload.Count = deck.Remaining
	}

	drawnCodes, err := models.DrawRandomNCars(payload.Count, deck.CardCodes)
	if err != nil {
		return errors.New(errors.Internal, "failed draw cards from deck")
	}

	cards, err := models.BuildCardsFromCodes(drawnCodes)
	if err != nil {
		return errors.New(errors.Internal, "failed map codes to cards")
	}

	leftCodes := models.RemoveDrawnCodes(drawnCodes, deck.CardCodes)

	deck.CardCodes = leftCodes

	if err = h.repo.UpdateDeck(deck); err != nil {
		return errors.New(errors.Internal, "failed draw cards")
	}

	return httphelper.WriteSuccessResponse(w, http.StatusOK, cards)
}
