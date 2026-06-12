package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"tictactoe/internal/application/port"
	"tictactoe/internal/application/usecase"
	"tictactoe/internal/domain/model"
)

func TestCreateSessionRepoErr(t *testing.T) {
	errRepo := errors.New("session db error")

	sessionRepo := &mockSessionRepo{}
	sessionRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Session")).Return(errRepo)

	cmd := &port.CreateCommand{
		Rules:         model.NewDefaultRules(),
		SessionParams: &model.SessionParams{},
	}
	res, err := usecase.NewCreate(sessionRepo).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errRepo)
	assert.Nil(t, res)
	sessionRepo.AssertExpectations(t)
}

func TestCreateSuccess(t *testing.T) {
	sessionRepo := &mockSessionRepo{}
	sessionRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Session")).Return(nil)

	cmd := &port.CreateCommand{
		Rules:         model.NewDefaultRules(),
		SessionParams: &model.SessionParams{},
	}
	res, err := usecase.NewCreate(sessionRepo).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	sessionRepo.AssertExpectations(t)
}
