package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildCardsFromCodes(t *testing.T) {
	var (
		err   error
		cards Cards
	)

	cards, err = BuildCardsFromCodes([]string{"RR", "RR"})
	require.Nil(t, cards)
	require.EqualError(t, err, "codes is not valid")

	cards, err = BuildCardsFromCodes([]string{})
	require.Empty(t, cards)
	require.NoError(t, err)

	expectedCards := []*Card{
		{
			Value: "ACE",
			Suit:  "SPADES",
			Code:  "AS",
		},
		{
			Value: "KING",
			Suit:  "DIAMONDS",
			Code:  "KD",
		},
		{
			Value: "2",
			Suit:  "CLUBS",
			Code:  "2C",
		},
		{
			Value: "10",
			Suit:  "CLUBS",
			Code:  "10C",
		},
	}

	cards, err = BuildCardsFromCodes([]string{"AS", "KD", "2C", "10C"})
	require.NoError(t, err)
	require.EqualValues(t, expectedCards, cards)
}

func TestDrawRandomNCars(t *testing.T) {
	var (
		result []string
		err    error
	)

	result, err = DrawRandomNCars(3, []string{"AS", "KD", "2C", "10C"})
	require.NoError(t, err)
	require.Len(t, result, 3)

	result, err = DrawRandomNCars(3, []string{"AS", "KD", "2C"})
	require.NoError(t, err)
	require.Len(t, result, 3)

	result, err = DrawRandomNCars(4, []string{"AS", "KD", "2C"})
	require.EqualError(t, err, "count cannot be greater then length of codes slice")
}

func TestRemoveDrawnCodes(t *testing.T) {
	var result []string

	result = RemoveDrawnCodes([]string{"KD"}, []string{"AS", "KD", "2C"})
	require.Len(t, result, 2)
	require.Equal(t, []string{"AS", "2C"}, result)

	result = RemoveDrawnCodes([]string{"KD", "2C"}, []string{"AS", "KD", "2C"})
	require.Len(t, result, 1)
	require.Equal(t, []string{"AS"}, result)

	result = RemoveDrawnCodes([]string{"2C"}, []string{"AS", "KD", "2C"})
	require.Len(t, result, 2)
	require.Equal(t, []string{"AS", "KD"}, result)

	result = RemoveDrawnCodes([]string{"AS"}, []string{"AS", "KD", "2C"})
	require.Len(t, result, 2)
	require.Equal(t, []string{"KD", "2C"}, result)

	result = RemoveDrawnCodes([]string{"AS"}, []string{"AS"})
	require.Len(t, result, 0)
	require.Equal(t, []string{}, result)

	result = RemoveDrawnCodes([]string{"AS"}, []string{})
	require.Len(t, result, 0)
	require.Equal(t, []string{}, result)
}
