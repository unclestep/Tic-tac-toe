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

func TestCreateRulesRepoErr(t *testing.T) {
	errRepo := errors.New("rules db error")

	rulesRepo := &mockRulesRepo{}
	rulesRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Rules")).Return(errRepo)

	sessionRepo := &mockSessionRepo{}

	cmd := &port.CreateCommand{
		Rules:         model.NewDefaultRules(),
		SessionParams: &model.SessionParams{},
	}
	res, err := usecase.NewCreate(sessionRepo, rulesRepo).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errRepo)
	assert.Nil(t, res)
	sessionRepo.AssertNotCalled(t, "Save")
	rulesRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestCreateSessionRepoErr(t *testing.T) {
	errRepo := errors.New("session db error")

	rulesRepo := &mockRulesRepo{}
	rulesRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Rules")).Return(nil)

	sessionRepo := &mockSessionRepo{}
	sessionRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Session")).Return(errRepo)

	cmd := &port.CreateCommand{
		Rules:         model.NewDefaultRules(),
		SessionParams: &model.SessionParams{},
	}
	res, err := usecase.NewCreate(sessionRepo, rulesRepo).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errRepo)
	assert.Nil(t, res)
	rulesRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestCreateSuccess(t *testing.T) {
	rulesRepo := &mockRulesRepo{}
	rulesRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Rules")).Return(nil)

	sessionRepo := &mockSessionRepo{}
	sessionRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.Session")).Return(nil)

	cmd := &port.CreateCommand{
		Rules:         model.NewDefaultRules(),
		SessionParams: &model.SessionParams{},
	}
	res, err := usecase.NewCreate(sessionRepo, rulesRepo).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Rules.UUID)
	rulesRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}
