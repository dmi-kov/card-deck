package deckhelper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDefaultCodes(t *testing.T) {
	gotCodes := CreateDefaultCodes()
	require.Len(t, gotCodes, 52)
}

func TestIsValidCodes(t *testing.T) {
	var isValid bool

	isValid = IsValidCodes([]string{"AC", "KD", "2C"})
	require.Equal(t, true, isValid)

	isValid = IsValidCodes([]string{"AC", "KD", "C2"})
	require.Equal(t, false, isValid)

	isValid = IsValidCodes([]string{"AC", "kD", "2C"})
	require.Equal(t, false, isValid)

	isValid = IsValidCodes([]string{"AC", "AC", "2C"})
	require.Equal(t, false, isValid)
}
