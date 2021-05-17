package deckhelper

import (
	"fmt"
	"math/rand"
	"time"
)

// Constants for standard 52-card deck
const (
	ACE   = "A"
	TWO   = "2"
	THREE = "3"
	FOUR  = "4"
	FIVE  = "5"
	SIX   = "6"
	SEVEN = "7"
	EIGHT = "8"
	NINE  = "9"
	TEN   = "10"
	JACK  = "J"
	QUEEN = "Q"
	KING  = "K"

	CLUBS    = "C"
	DIAMONDS = "D"
	HEARTS   = "H"
	SPADES   = "S"
)

// Default sequences of cards
var (
	SuitsSequence = []string{SPADES, DIAMONDS, CLUBS, HEARTS}
	CardsSequence = []string{ACE, TWO, THREE, FOUR, FIVE, SIX, SEVEN, EIGHT, NINE, TEN, JACK, QUEEN, KING}
)

// SuitsMapping maps suit code to suit name
var SuitsMapping = map[string]string{
	SPADES:   "SPADES",
	CLUBS:    "CLUBS",
	DIAMONDS: "DIAMONDS",
	HEARTS:   "HEARTS",
}

// CardsMapping maps card code to card name
var CardsMapping = map[string]string{
	ACE:   "ACE",
	TWO:   TWO,
	THREE: THREE,
	FOUR:  FOUR,
	FIVE:  FIVE,
	SIX:   SIX,
	SEVEN: SEVEN,
	EIGHT: EIGHT,
	NINE:  NINE,
	TEN:   TEN,
	JACK:  "JACK",
	QUEEN: "QUEEN",
	KING:  "KING",
}

// CreateDefaultCodes Creates default cards sequence
func CreateDefaultCodes() (codes []string) {
	for _, suit := range SuitsSequence {
		for _, card := range CardsSequence {
			code := fmt.Sprintf("%s%s", card, suit)
			codes = append(codes, code)
		}
	}
	return codes
}

// ShuffleDeck returns shuffled deck codes
func ShuffleDeck(codes []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(codes), func(i, j int) { codes[i], codes[j] = codes[j], codes[i] })
}

// IsValidCodes checks if given list of codes is valid
// by comparing with default codes and check if duplicates presented
func IsValidCodes(codes []string) bool {
	defaultCodes := CreateDefaultCodes()
	for _, code := range codes {
		if !contains(defaultCodes, code) {
			return false
		}
	}

	keys := make(map[string]bool)
	for _, code := range codes {
		if _, value := keys[code]; value {
			return false
		}
		keys[code] = true

	}

	return true
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
