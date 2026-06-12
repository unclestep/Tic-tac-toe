package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tictactoe/internal/domain/model"
	"tictactoe/internal/domain/service"
	"tictactoe/pkg/geometry"
)

func newBotLogicSetup(t *testing.T, width, height, winLen int) (*service.BotMovement, *model.Session) {
	t.Helper()
	rules := model.NewDefaultRules()
	board := model.NewBoard(width, height)
	checker := service.NewWinChecker()
	heuristic := service.NewHeuristic()
	bot := service.NewBotMovement(checker, heuristic)
	session := model.NewSession("1", &model.SessionParams{}, rules, board)
	return bot, session
}

func TestMakeMoveStateGameOver(t *testing.T) {
	bot, session := newBotLogicSetup(t, 3, 3, 3)
	session.State = model.StateGameOver

	err := bot.MakeMove(session, model.MarkX)

	assert.ErrorIs(t, err, model.ErrGameAlreadyOver)
}

func TestMakeMoveFullBoard(t *testing.T) {
	bot, session := newBotLogicSetup(t, 2, 2, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(1, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(0, 1)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(1, 1)))

	err := bot.MakeMove(session, model.MarkX)

	assert.ErrorIs(t, err, model.ErrGameAlreadyOver)
}

func TestMakeMoveBotXWinsImmediately(t *testing.T) {
	// XX.
	// O..
	// ...
	bot, session := newBotLogicSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(1, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(0, 1)))

	err := bot.MakeMove(session, model.MarkX)

	require.NoError(t, err)
	assert.Equal(t, model.StateGameOver, session.State)
	assert.Equal(t, *model.NewBot(model.MarkX), *session.Winner)
}

func TestMakeMoveBotOWinsImmediately(t *testing.T) {
	// X..
	// OO.
	// X..
	bot, session := newBotLogicSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 2)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(0, 1)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(1, 1)))

	err := bot.MakeMove(session, model.MarkO)

	require.NoError(t, err)
	assert.Equal(t, model.StateGameOver, session.State)
	assert.Equal(t, *model.NewBot(model.MarkO), *session.Winner)
}

func TestMakeMovePrefersWinOverBlock(t *testing.T) {
	// XX.  <- X wins at (2,0)
	// OO.  <- O threatens at (2,1)
	// ...
	bot, session := newBotLogicSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(1, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(0, 1)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(1, 1)))

	err := bot.MakeMove(session, model.MarkX)

	require.NoError(t, err)
	assert.Equal(t, model.StateGameOver, session.State)
	assert.Equal(t, *model.NewBot(model.MarkX), *session.Winner)
}

func TestMakeMoveBotXBlocksO(t *testing.T) {
	// .OO  <- O wins at (0,0) if unblocked
	// X..
	// ...
	bot, session := newBotLogicSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(1, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(2, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 1)))

	err := bot.MakeMove(session, model.MarkX)

	require.NoError(t, err)
	mark, getErr := session.Board.GetMark(geometry.NewPoint(0, 0))
	require.NoError(t, getErr)
	assert.Equal(t, model.MarkX, mark)
}

func TestMakeMoveBotOBlocksX(t *testing.T) {
	// XX.  <- X wins at (2,0) if unblocked
	// O..
	// ...
	bot, session := newBotLogicSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(1, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(0, 1)))

	err := bot.MakeMove(session, model.MarkO)

	require.NoError(t, err)
	mark, getErr := session.Board.GetMark(geometry.NewPoint(2, 0))
	require.NoError(t, getErr)
	assert.Equal(t, model.MarkO, mark)
}

func TestMakeMoveDraw(t *testing.T) {
	// 2x2 board, WinLength=3: winning is impossible
	// XO
	// X.  <- bot O plays last cell (1,1)
	bot, session := newBotLogicSetup(t, 2, 2, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(1, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 1)))

	err := bot.MakeMove(session, model.MarkO)

	require.NoError(t, err)
	assert.Equal(t, model.StateGameOver, session.State)
	assert.Nil(t, session.Winner)
}

func TestMakeMoveTurnIncrements(t *testing.T) {
	bot, session := newBotLogicSetup(t, 3, 3, 3)

	require.NoError(t, bot.MakeMove(session, model.MarkX))

	assert.Equal(t, 1, session.Turn)
}

func TestMakeMoveStateRemainsPlayingWhileNoWinner(t *testing.T) {
	// X..
	// O..
	// ...
	bot, session := newBotLogicSetup(t, 3, 3, 3)
	require.NoError(t, session.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))
	require.NoError(t, session.Board.SetMark(model.MarkO, geometry.NewPoint(0, 1)))

	err := bot.MakeMove(session, model.MarkX)

	require.NoError(t, err)
	assert.Empty(t, session.Winner)
}
