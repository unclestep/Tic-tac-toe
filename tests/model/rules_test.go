package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"tictactoe/internal/domain/model"
)

func TestRulesClone(t *testing.T) {
	original := model.NewDefaultRules()
	assert.NotNil(t, original)

	clone := original.Clone()
	assert.NotNil(t, clone)

	assert.Equal(t, *original, *clone)

	clone.UUID = "b"
	assert.NotEqual(t, clone.UUID, original.UUID)

