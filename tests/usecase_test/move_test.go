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
	"tictactoe/pkg/geometry"
)

func TestMakeMoveSessionNotFound(t *testing.T) {
	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "nonexistent").Return(nil, port.ErrSessionNotFound)

	cmd := &port.MakeMoveCommand{SessionUUID: "nonexistent", PlayerUUID: "p1"}
	res, err := usecase.NewMakeMove(repo, nil, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrSessionNotFound)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestMakeMoveGameNotStarted(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p1"}
	res, err := usecase.NewMakeMove(repo, nil, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrGameNotStarted)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestMakeMoveGameAlreadyOver(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.State = model.StateGameOver

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p1"}
	res, err := usecase.NewMakeMove(repo, nil, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, model.ErrGameAlreadyOver)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestMakeMovePlayerNotBelongToSession(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Start()

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "nonexistent"}
	res, err := usecase.NewMakeMove(repo, nil, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrPlayerNotBelongToSession)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestMakeMovePlayerCantMakeMove(t *testing.T) {
	session := newSessionWithPlayers(
		t,
		model.NewPlayer("p1", "p1", model.MarkX),
		model.NewPlayer("p2", "p2", model.MarkO),
	)
	session.Start()

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p2"}
	res, err := usecase.NewMakeMove(repo, nil, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, port.ErrPlayerCantMakeMove)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	repo.AssertExpectations(t)
}

func TestMakeMoveHumanMoverErr(t *testing.T) {
	errHuman := errors.New("human move error")

	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Start()
	turnPlayer, _ := session.GetTurnPlayer()

	humanMover := &mockHumanMover{}
	humanMover.On("MakeMove", session, turnPlayer, geometry.Point{}).Return(errHuman)

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p1", MarkPos: geometry.Point{}}
	res, err := usecase.NewMakeMove(repo, humanMover, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errHuman)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	humanMover.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestMakeMoveBotMoverErr(t *testing.T) {
	errBot := errors.New("bot move error")

	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Start()
	turnPlayer, _ := session.GetTurnPlayer()

	humanMover := &mockHumanMover{}
	humanMover.On("MakeMove", session, turnPlayer, geometry.Point{}).Return(nil)

	botMover := &mockBotMover{}
	botMover.On("MakeMove", session, model.MarkO).Return(errBot)

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p1", MarkPos: geometry.Point{}}
	res, err := usecase.NewMakeMove(repo, humanMover, botMover).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errBot)
	assert.Nil(t, res)
	repo.AssertNotCalled(t, "Save")
	humanMover.AssertExpectations(t)
	botMover.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestMakeMoveSessionRepoErr(t *testing.T) {
	errRepo := errors.New("session db error")

	session := newSessionWithPlayers(
		t,
		model.NewPlayer("p1", "p1", model.MarkX),
		model.NewPlayer("p2", "p2", model.MarkO),
	)
	session.Start()
	turnPlayer, _ := session.GetTurnPlayer()

	humanMover := &mockHumanMover{}
	humanMover.On("MakeMove", session, turnPlayer, geometry.Point{}).Return(nil)

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(errRepo)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p1", MarkPos: geometry.Point{}}
	res, err := usecase.NewMakeMove(repo, humanMover, nil).Execute(context.Background(), cmd)

	assert.ErrorIs(t, err, errRepo)
	assert.Nil(t, res)
	humanMover.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestMakeMoveSuccessFullSession(t *testing.T) {
	session := newSessionWithPlayers(
		t,
		model.NewPlayer("p1", "p1", model.MarkX),
		model.NewPlayer("p2", "p2", model.MarkO),
	)
	session.Start()
	turnPlayer, _ := session.GetTurnPlayer()

	humanMover := &mockHumanMover{}
	humanMover.On("MakeMove", session, turnPlayer, geometry.Point{}).Return(nil)

	botMover := &mockBotMover{}

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p1", MarkPos: geometry.Point{}}
	res, err := usecase.NewMakeMove(repo, humanMover, botMover).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	botMover.AssertNotCalled(t, "MakeMove")
	humanMover.AssertExpectations(t)
	botMover.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestMakeMoveSuccessSoloSession(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Start()
	turnPlayer, _ := session.GetTurnPlayer()

	humanMover := &mockHumanMover{}
	humanMover.On("MakeMove", session, turnPlayer, geometry.Point{}).Return(nil)

	botMover := &mockBotMover{}
	botMover.On("MakeMove", session, model.MarkO).Return(nil)

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p1", MarkPos: geometry.Point{}}
	res, err := usecase.NewMakeMove(repo, humanMover, botMover).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	humanMover.AssertExpectations(t)
	botMover.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestMakeMoveGameOverAfterHumanMove(t *testing.T) {
	session := newSessionWithPlayers(t, model.NewPlayer("p1", "p1", model.MarkX))
	session.Start()
	turnPlayer, _ := session.GetTurnPlayer()

	humanMover := &mockHumanMover{}
	humanMover.On("MakeMove", session, turnPlayer, geometry.Point{}).
		Run(func(args mock.Arguments) {
			args.Get(0).(*model.Session).State = model.StateGameOver
		}).
		Return(nil)

	botMover := &mockBotMover{}

	repo := &mockSessionRepo{}
	repo.On("Get", mock.Anything, "session-1").Return(session, nil)
	repo.On("Save", mock.Anything, session).Return(nil)

	cmd := &port.MakeMoveCommand{SessionUUID: "session-1", PlayerUUID: "p1", MarkPos: geometry.Point{}}
	res, err := usecase.NewMakeMove(repo, humanMover, botMover).Execute(context.Background(), cmd)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	botMover.AssertNotCalled(t, "MakeMove")
	humanMover.AssertExpectations(t)
	repo.AssertExpectations(t)
}
