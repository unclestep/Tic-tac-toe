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

// Experimentaly chosen
const (
	seedPlayerFirst = 0
	seedBotFirst    = 1
)

func TestStartSessionNotFound(t *testing.T) {
	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "nonexistent").Return(nil, port.ErrSessionNotFound)

	cmd := &port.StartCommand{SessionUUID: "nonexistent", PlayerUUID: "p1"}
	res, err := usecase.NewStart(repo, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrSessionNotFound)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestStartGameAlreadyStarted(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Start()

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.StartCommand{SessionUUID: "session-1", PlayerUUID: "p1"}
	res, err := usecase.NewStart(repo, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrGameAlreadyStarted)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestStartPlayerNotBelongToSession(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.StartCommand{SessionUUID: "session-1", PlayerUUID: "nonexistent"}
	res, err := usecase.NewStart(repo, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrPlayerNotBelongToSession)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestStartFullSession(t *testing.T) {
	session := newSessionWithPlayers(
		t,
		model.NewPlayer("p1", "p1", model.MarkX),
		model.NewPlayer("p2", "p2", model.MarkO),
	)

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	botMover := &mockBotMover{}

	cmd := &port.StartCommand{SessionUUID: "session-1", PlayerUUID: "p1"}
	res, err := usecase.NewStart(repo, botMover).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, model.MarkX, session.Players[0].Mark)
	assert.Equal(t, model.MarkO, session.Players[1].Mark)
	botMover.AssertNotCalled(t, "MakeMove")
	repo.AssertExpectations(t)
	botMover.AssertExpectations(t)
}

func TestStartPlayerGoesFirst(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Params.Seed = seedPlayerFirst

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	botMover := &mockBotMover{}

	cmd := &port.StartCommand{SessionUUID: "session-1", PlayerUUID: "p1"}
	res, err := usecase.NewStart(repo, botMover).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, model.MarkX, session.Players[0].Mark)
	botMover.AssertNotCalled(t, "MakeMove")
	repo.AssertExpectations(t)
	botMover.AssertExpectations(t)
}

func TestStartBotGoesFirst(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Params.Seed = seedBotFirst

	botMover := &mockBotMover{}
	botMover.On("MakeMove", session, model.MarkX).Return(nil)

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	cmd := &port.StartCommand{SessionUUID: "session-1", PlayerUUID: "p1"}
	res, err := usecase.NewStart(repo, botMover).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, model.MarkO, session.Players[0].Mark)
	botMover.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestStartBotMoverErr(t *testing.T) {
	errBot := errors.New("bot move error")

	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Params.Seed = seedBotFirst

	botMover := &mockBotMover{}
	botMover.On("MakeMove", session, model.MarkX).Return(errBot)

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.StartCommand{SessionUUID: "session-1", PlayerUUID: "p1"}
	res, err := usecase.NewStart(repo, botMover).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errBot)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	botMover.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestStartSaveErr(t *testing.T) {
	errRepo := errors.New("session db error")
	session := newSessionWithPlayers(
		t,
		model.NewPlayer("p1", "p1", model.MarkX),
		model.NewPlayer("p2", "p2", model.MarkO),
	)

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(errRepo)

	cmd := &port.StartCommand{SessionUUID: "session-1", PlayerUUID: "p1"}
	res, err := usecase.NewStart(repo, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errRepo)
	assert.Nil(t, res)
	repo.AssertExpectations(t)
}
