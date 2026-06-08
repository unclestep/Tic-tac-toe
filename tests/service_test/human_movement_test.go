package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tictactoe/internal/domain/model"
	"tictactoe/internal/domain/service"
	"tictactoe/pkg/geometry"
)

func newMovementSetup(t *testing.T, width, height, winLen int) (*service.HumanMovement, *model.Session) {
	t.Helper()
	rules := &model.Rules{UUID: "rules-1", BoardWidth: width, BoardHeight: height, WinLength: winLen}
	checker := service.NewWinChecker()
	movement := service.NewHumanMovement(checker)
	board := model.NewBoard(width, height)
	session := model.NewSession("session-1", &model.SessionParams{Seed: 0}, rules, board)
	return movement, session
}

func newPlayerWithMark(uuid string, mark model.Mark) *model.Player {
	p := model.NewPlayer(uuid, "player", model.MarkX)
	p.Mark = mark
	return p
}

func TestMakeOutOfBounds(t *testing.T) {
	movement, session := newMovementSetup(t, 3, 3, 3)
	playerX := newPlayerWithMark("p1", model.MarkX)

	err := movement.MakeMove(session, playerX, geometry.NewPoint(10, 10))

	assert.ErrorIs(t, err, model.ErrOutOfBounds)
	assert.Equal(t, 0, session.Turn)
}

func TestMakeCellOccupied(t *testing.T) {
	movement, session := newMovementSetup(t, 3, 3, 3)
	playerX := newPlayerWithMark("p1", model.MarkX)
	playerO := newPlayerWithMark("p2", model.MarkO)

	require.NoError(t, movement.MakeMove(session, playerX, geometry.NewPoint(0, 0)))
	err := movement.MakeMove(session, playerO, geometry.NewPoint(0, 0))

	assert.ErrorIs(t, err, model.ErrCellNotEmpty)
	assert.Equal(t, 1, session.Turn)
}

func TestMakeValidMoveGameContinues(t *testing.T) {
	movement, session := newMovementSetup(t, 3, 3, 3)
	playerX := newPlayerWithMark("p1", model.MarkX)

	err := movement.MakeMove(session, playerX, geometry.NewPoint(0, 0))

	require.NoError(t, err)
	assert.Empty(t, session.Winner)
	assert.Equal(t, 1, session.Turn)
}

func TestMakeTurnIncrements(t *testing.T) {
	movement, session := newMovementSetup(t, 3, 3, 3)
	playerX := newPlayerWithMark("p1", model.MarkX)
	playerO := newPlayerWithMark("p2", model.MarkO)

	require.NoError(t, movement.MakeMove(session, playerX, geometry.NewPoint(0, 0)))
	require.NoError(t, movement.MakeMove(session, playerO, geometry.NewPoint(1, 0)))
	require.NoError(t, movement.MakeMove(session, playerX, geometry.NewPoint(2, 0)))

	assert.Equal(t, 3, session.Turn)
}

func TestMakeXWins(t *testing.T) {
	movement, session := newMovementSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 1)))
	playerX := newPlayerWithMark("p1", model.MarkX)

	err := movement.MakeMove(session, playerX, geometry.NewPoint(0, 2))

	require.NoError(t, err)
	assert.Equal(t, model.StateGameOver, session.State)
	assert.Equal(t, model.MarkX.String(), session.Winner)
	assert.Equal(t, 1, session.Turn)
}

func TestMakeOWins(t *testing.T) {
	movement, session := newMovementSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(1, 0)))
	playerO := newPlayerWithMark("p2", model.MarkO)

	err := movement.MakeMove(session, playerO, geometry.NewPoint(2, 0))

	require.NoError(t, err)
	assert.Equal(t, model.StateGameOver, session.State)
	assert.Equal(t, model.MarkO.String(), session.Winner)
	assert.Equal(t, 1, session.Turn)
}

func TestMakeDraw(t *testing.T) {
	movement, session := newMovementSetup(t, 2, 2, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(1, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 1)))
	playerO := newPlayerWithMark("p2", model.MarkO)

	err := movement.MakeMove(session, playerO, geometry.NewPoint(1, 1))

	require.NoError(t, err)
	assert.Equal(t, model.StateGameOver, session.State)
	assert.Equal(t, model.MarkEmpty.String(), session.Winner)
	assert.Equal(t, 1, session.Turn)
}

func TestMakeAfterGameOver(t *testing.T) {
	movement, session := newMovementSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 1)))
	playerX := newPlayerWithMark("p1", model.MarkX)

	require.NoError(t, movement.MakeMove(session, playerX, geometry.NewPoint(0, 2)))
	require.Equal(t, model.StateGameOver, session.State)
	turnAfterWin := session.Turn

	err := movement.MakeMove(session, playerX, geometry.NewPoint(1, 0))

	assert.ErrorIs(t, err, model.ErrGameAlreadyOver)
	assert.Equal(t, turnAfterWin, session.Turn)
}
