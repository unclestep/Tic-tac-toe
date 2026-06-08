package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"tictactoe/internal/domain/model"
)

func TestNewPlayer(t *testing.T) {
	p := model.NewPlayer("a", "a", model.MarkX)
	assert.NotNil(t, p)
}

func TestPlayerClone(t *testing.T) {
	original := model.NewPlayer("a", "a", model.MarkX)
	assert.NotNil(t, original)

	clone := original.Clone()
	assert.NotNil(t, clone)

	assert.Equal(t, *original, *clone)

	clone.UUID = "b"
	assert.NotEqual(t, clone.UUID, original.UUID)
}
