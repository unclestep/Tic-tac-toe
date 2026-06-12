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
	repo.On("Get", mock.Anything, "?").Return(nil, port.ErrSessionNotFound)

	cmd := &port.ConnectCommand{SessionUUID: "?"}
	res, err := usecase.NewConnect(repo).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrSessionNotFound)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestConnectSessionFull(t *testing.T) {
	session := newSessionWithPlayers(t, newPlayer("1", model.MarkX), newPlayer("2", model.MarkO))
	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)

	cmd := &port.ConnectCommand{SessionUUID: "1"}
	res, err := usecase.NewConnect(repo).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, model.ErrSessionFull)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestConnectSessionRepoErr(t *testing.T) {
	errRepo := errors.New("session db error")
	session := newSessionWithPlayers(t, newPlayer("1", model.MarkX))
	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(errRepo)

	ctx := context.WithValue(context.Background(), port.UserUUIDKey, "1")
	cmd := &port.ConnectCommand{SessionUUID: "1"}
	res, err := usecase.NewConnect(repo).Execute(ctx, cmd)

	assert.ErrorIs(t, err, errRepo)
	assert.Equal(t, 2, len(session.Players))
	assert.Nil(t, res)
	repo.AssertExpectations(t)
}

func TestConnectPlayers(t *testing.T) {
	session := newSessionWithPlayers(t)
	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	ctx := context.WithValue(context.Background(), port.UserUUIDKey, "1")
	cmd := &port.ConnectCommand{SessionUUID: "1"}

	res, err := usecase.NewConnect(repo).Execute(ctx, cmd)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(session.Players))
	assert.NotNil(t, res)

	res, err = usecase.NewConnect(repo).Execute(ctx, cmd)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(session.Players))
	assert.NotNil(t, res)

	res, err = usecase.NewConnect(repo).Execute(ctx, cmd)
	assert.ErrorIs(t, err, model.ErrSessionFull)
	assert.Equal(t, 2, len(session.Players))
	assert.Nil(t, res)

	repo.AssertExpectations(t)
}
