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

func TestConnectSessionNotFound(t *testing.T) {
	repo := &mockSessionRepo{}

	repo.On("Get", mock.Anything, "Nonexistent").Return(nil, port.ErrSessionNotFound)
	cmd := &port.ConnectCommand{SessionUUID: "Nonexistent"}
	res, err := usecase.NewConnect(repo).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrSessionNotFound)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestConnectSessionFull(t *testing.T) {
	session := newSessionWithPlayers(
		t,
		model.NewPlayer("p1", "p1", model.MarkX),
		model.NewPlayer("p2", "p2", model.MarkO),
	)
	repo := &mockSessionRepo{}

	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	cmd := &port.ConnectCommand{SessionUUID: "session-1"}
	res, err := usecase.NewConnect(repo).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, model.ErrSessionFull)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestConnectSessionRepoErr(t *testing.T) {
	errRepo := errors.New("session db error")
	session := newSessionWithPlayers(
		t,
		model.NewPlayer("p1", "p1", model.MarkX),
	)
	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(errRepo)
	cmd := &port.ConnectCommand{SessionUUID: "session-1"}
	res, err := usecase.NewConnect(repo).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errRepo)
	repo.AssertExpectations(t)
	assert.Equal(t, 2, len(session.Players))
	assert.Nil(t, res)
}

func TestConnectPlayers(t *testing.T) {
	session := newSessionWithPlayers(t)
	repo := &mockSessionRepo{}

	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	cmd := &port.ConnectCommand{SessionUUID: "session-1"}

	res, err := usecase.NewConnect(repo).Execute(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(session.Players))
	assert.NotNil(t, res)

	res, err = usecase.NewConnect(repo).Execute(context.Background(), cmd)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(session.Players))
	assert.NotNil(t, res)

	res, err = usecase.NewConnect(repo).Execute(context.Background(), cmd)
	assert.ErrorIs(t, err, model.ErrSessionFull)
	assert.Equal(t, 2, len(session.Players))
	assert.Nil(t, res)

	repo.AssertExpectations(t)
}
