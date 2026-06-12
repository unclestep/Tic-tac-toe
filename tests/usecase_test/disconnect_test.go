package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"tictactoe/internal/application/port"
	"tictactoe/internal/application/usecase"
	"tictactoe/internal/domain/model"
)

func TestDisconnectSessionNotFound(t *testing.T) {
	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(nil, port.ErrSessionNotFound)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "1"}
	_, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrSessionNotFound)
	repo.AssertExpectations(t)
	bot.AssertNotCalled(t, "MakeMove")
}

func TestDisconnectGetRepoError(t *testing.T) {
	repoErr := errors.New("db connection failed")
	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(nil, repoErr)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "1"}
	_, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
	bot.AssertNotCalled(t, "MakeMove")
}

func TestDisconnectEmptySessionGetTurnPlayerFails(t *testing.T) {
	session := newSessionWithPlayers(t)
	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "1"}
	_, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, model.ErrPlayerNotFound)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Save")
	bot.AssertNotCalled(t, "MakeMove")
}

func TestDisconnectPlayerNotInSession(t *testing.T) {
	session := newSessionWithPlayers(t, newPlayer("1", model.MarkX))
	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "?"}
	_, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, model.ErrPlayerNotFound)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Save")
	bot.AssertNotCalled(t, "MakeMove")
}

func TestDisconnectBotDontMoveInLobby(t *testing.T) {
	session := newSessionWithPlayers(t, newPlayer("1", model.MarkX))
	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "1"}
	_, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	bot.AssertNotCalled(t, "MakeMove")
}

func TestDisconnectBotMoveError(t *testing.T) {
	botErr := errors.New("bot make move error")
	session := newSessionWithPlayers(t, newPlayer("1", model.MarkX), newPlayer("2", model.MarkO))
	session.State = model.StatePlaying

	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)
	bot.On("MakeMove", session, model.MarkX).Return(botErr)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "1"}
	_, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, botErr)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "Save")
	bot.AssertExpectations(t)
}

func TestDisconnectSaveSessionError(t *testing.T) {
	repoErr := errors.New("session save err")
	session := newSessionWithPlayers(t, newPlayer("1", model.MarkX), newPlayer("2", model.MarkO))

	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(repoErr)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "1"}
	_, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, repoErr)
	repo.AssertExpectations(t)
}

func TestDisconnectTurnPlayerTriggersBotMove(t *testing.T) {
	session := newSessionWithPlayers(t, newPlayer("1", model.MarkX), newPlayer("2", model.MarkO))
	session.State = model.StatePlaying

	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)
	bot.On("MakeMove", session, model.MarkX).Return(nil)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "1"}
	result, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.IsPlayerExist("1"))
	repo.AssertExpectations(t)
	bot.AssertExpectations(t)
}

func TestDisconnectNotPlayerTurnBotDontMove(t *testing.T) {
	session := newSessionWithPlayers(t, newPlayer("1", model.MarkX), newPlayer("2", model.MarkO))
	session.State = model.StatePlaying

	repo := &mockSessionRepo{}
	bot := &mockBotMover{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	cmd := &port.DisconnectCommand{SessionUUID: "1", PlayerUUID: "2"}
	result, err := usecase.NewDisconnect(repo, bot).Execute(context.Background(), cmd)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.IsPlayerExist("1"))
	assert.False(t, result.IsPlayerExist("2"))
	repo.AssertExpectations(t)
	bot.AssertNotCalled(t, "MakeMove")
}
