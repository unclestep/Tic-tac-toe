package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tictactoe/internal/domain/model"
)

func newTestSession(t *testing.T) *model.Session {
	t.Helper()
	board := model.NewBoard(3, 3)
	rules := model.NewDefaultRules()
	return model.NewSession("1", &model.SessionParams{Seed: 0}, rules, board)
}

func newTestPlayer1(t *testing.T) *model.Player {
	t.Helper()
	return model.NewPlayer("1", "1", "1", model.MarkX)
}

func newTestPlayer2(t *testing.T) *model.Player {
	t.Helper()
	return model.NewPlayer("2", "2", "2", model.MarkO)
}

func newTestPlayer3(t *testing.T) *model.Player {
	t.Helper()
	return model.NewPlayer("3", "3", "3", model.MarkO)
}

func TestCloneSessionFieldsMatch(t *testing.T) {
	s := newTestSession(t)
	clone := s.Clone()
	assert.Equal(t, s, clone)
}

func TestCloneSessionParamsAreIndependent(t *testing.T) {
	s := newTestSession(t)
	clone := s.Clone()

	clone.Board.Cells[0] = 1
	clone.Winner = newTestPlayer1(t)

	assert.NotEqual(t, *clone, *s)
}

func TestIsPlayerExistEmptySession(t *testing.T) {
	s := newTestSession(t)
	assert.False(t, s.IsPlayerExist("1"))
}

func TestIsPlayerExistFound(t *testing.T) {
	s := newTestSession(t)
	err := s.AddPlayer(newTestPlayer1(t))
	require.NoError(t, err)
	assert.True(t, s.IsPlayerExist("1"))
}

func TestIsPlayerExistUnknownUUID(t *testing.T) {
	s := newTestSession(t)
	err := s.AddPlayer(newTestPlayer1(t))
	require.NoError(t, err)
	assert.False(t, s.IsPlayerExist("?"))
}

func TestIsFullEmpty(t *testing.T) {
	s := newTestSession(t)
	assert.False(t, s.IsFull())
}

func TestIsFullOnePlayer(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	assert.False(t, s.IsFull())
}

func TestIsFullTwoPlayers(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))
	assert.True(t, s.IsFull())
}

func TestGetAvailableMarksEmptySession(t *testing.T) {
	s := newTestSession(t)
	assert.ElementsMatch(t, []model.Mark{model.MarkX, model.MarkO}, s.GetAvailableMarks())
}

func TestGetAvailableMarksOnePlayerX(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	assert.Equal(t, []model.Mark{model.MarkO}, s.GetAvailableMarks())
}

func TestGetAvailableMarksOnePlayerO(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))
	assert.Equal(t, []model.Mark{model.MarkX}, s.GetAvailableMarks())
}

func TestGetAvailableMarksFull(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))
	assert.Empty(t, s.GetAvailableMarks())
}

func TestAddPlayerSuccess(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	assert.Len(t, s.Players, 1)
	assert.True(t, s.IsPlayerExist("1"))
}

func TestAddPlayerTwoPlayersSuccess(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))
	assert.Len(t, s.Players, 2)
}

func TestAddPlayerSessionFull(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))
	assert.ErrorIs(t, s.AddPlayer(newTestPlayer3(t)), model.ErrSessionFull)
}

func TestAddPlayerDuplicateUUID(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	assert.ErrorIs(t, s.AddPlayer(newTestPlayer1(t)), model.ErrPlayerExists)
}

func TestAddPlayerMarkTaken(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))
	assert.ErrorIs(t, s.AddPlayer(newTestPlayer3(t)), model.ErrMarkTaken)
}

func TestRemovePlayerFirst(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))
	require.NoError(t, s.RemovePlayer("1"))

	assert.Len(t, s.Players, 1)
	assert.False(t, s.IsPlayerExist("1"))
	assert.True(t, s.IsPlayerExist("2"))
}

func TestRemovePlayerSecond(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))
	require.NoError(t, s.RemovePlayer("2"))

	assert.Len(t, s.Players, 1)
	assert.True(t, s.IsPlayerExist("1"))
	assert.False(t, s.IsPlayerExist("2"))
}

func TestRemovePlayerOnly(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	require.NoError(t, s.RemovePlayer("1"))
	assert.Empty(t, s.Players)
}

func TestRemovePlayerNotFound(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	err := s.RemovePlayer("?")
	assert.ErrorIs(t, err, model.ErrPlayerNotFound)
	assert.Equal(t, 1, len(s.Players))
}

func TestRemovePlayerEmptySession(t *testing.T) {
	s := newTestSession(t)
	err := s.RemovePlayer("1")
	assert.ErrorIs(t, err, model.ErrPlayerNotFound)
}

func TestGetTurnPlayerEmptySession(t *testing.T) {
	s := newTestSession(t)
	_, err := s.GetTurnPlayer()
	assert.ErrorIs(t, err, model.ErrPlayerNotFound)
}

func TestGetTurnPlayerRotation(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer1(t)))
	require.NoError(t, s.AddPlayer(newTestPlayer2(t)))

	cases := []struct {
		turn     int
		expected string
	}{
		{0, "1"},
		{1, "2"},
		{2, "1"},
		{3, "2"},
	}

	for _, tc := range cases {
		s.Turn = tc.turn
		player, err := s.GetTurnPlayer()
		require.NoError(t, err)
		assert.Equal(t, tc.expected, player.UUID, "turn %d", tc.turn)
	}
}
