package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"
)

func newTestSession(t *testing.T) *model.Session {
	t.Helper()
	board := model.NewBoard(3, 3)
	rules := &model.Rules{UUID: "rules-1", BoardWidth: 3, BoardHeight: 3}
	return model.NewSession("session-1", &model.SessionParams{Seed: 0}, rules, board)
}

func newTestPlayer(uuid string, mark model.Mark) *model.Player {
	p := model.NewPlayer(uuid, "player", mark)
	return p
}

func TestNewSessionInitialState(t *testing.T) {
	rules := &model.Rules{UUID: "rules-1", BoardWidth: 3, BoardHeight: 3}
	params := &model.SessionParams{Seed: 42}

	s := model.NewSession("session-1", params, rules, model.NewBoard(3, 3))

	assert.Equal(t, "session-1", s.UUID)
	assert.Equal(t, "rules-1", s.Rules.UUID)
	assert.Equal(t, model.StateLobby, s.State)
	assert.NotNil(t, s.Board)
	assert.Empty(t, s.Players)
	assert.Equal(t, params, s.Params)
}

func TestCloneSessionFieldsMatch(t *testing.T) {
	s := newTestSession(t)
	s.Turn = 3
	s.Winner = "Alice"
	s.State = model.StatePlaying

	clone := s.Clone()

	assert.Equal(t, s.UUID, clone.UUID)
	assert.Equal(t, s.Rules.UUID, clone.Rules.UUID)
	assert.Equal(t, s.Turn, clone.Turn)
	assert.Equal(t, s.Winner, clone.Winner)
	assert.Equal(t, s.State, clone.State)
	assert.Equal(t, s.Params.Seed, clone.Params.Seed)
}

func TestCloneSessionParamsAreIndependent(t *testing.T) {
	s := newTestSession(t)
	clone := s.Clone()

	clone.Params.Seed = 999

	assert.NotEqual(t, s.Params.Seed, clone.Params.Seed)
}

func TestCloneSessionBoardIsIndependent(t *testing.T) {
	s := newTestSession(t)
	clone := s.Clone()

	require.NoError(t, clone.Board.SetMark(model.MarkX, geometry.NewPoint(0, 0)))

	mark, err := s.Board.GetMark(geometry.NewPoint(0, 0))
	require.NoError(t, err)
	assert.Equal(t, model.MarkEmpty, mark)
}

func TestCloneSessionPlayersAreIndependent(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	clone := s.Clone()
	clone.Players[0].Name = "mutated"

	assert.NotEqual(t, "mutated", s.Players[0].Name)
}

func TestIsPlayerExistFound(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	assert.True(t, s.IsPlayerExist("p1"))
}

func TestIsPlayerExistEmptySession(t *testing.T) {
	s := newTestSession(t)
	assert.False(t, s.IsPlayerExist("p1"))
}

func TestIsPlayerExistUnknownUUID(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	assert.False(t, s.IsPlayerExist("unknown"))
}

func TestIsFullEmpty(t *testing.T) {
	s := newTestSession(t)
	assert.False(t, s.IsFull())
}

func TestIsFullOnePlayer(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	assert.False(t, s.IsFull())
}

func TestIsFullTwoPlayers(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))
	require.NoError(t, s.AddPlayer(newTestPlayer("p2", model.MarkO)))

	assert.True(t, s.IsFull())
}

func TestGetAvailableMarksEmptySession(t *testing.T) {
	s := newTestSession(t)
	assert.ElementsMatch(t, []model.Mark{model.MarkX, model.MarkO}, s.GetAvailableMarks())
}

func TestGetAvailableMarksOnePlayerX(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	assert.Equal(t, []model.Mark{model.MarkO}, s.GetAvailableMarks())
}

func TestGetAvailableMarksOnePlayerO(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkO)))

	assert.Equal(t, []model.Mark{model.MarkX}, s.GetAvailableMarks())
}

func TestGetAvailableMarksFull(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))
	require.NoError(t, s.AddPlayer(newTestPlayer("p2", model.MarkO)))

	assert.Empty(t, s.GetAvailableMarks())
}

func TestAddPlayerSuccess(t *testing.T) {
	s := newTestSession(t)

	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	assert.Len(t, s.Players, 1)
	assert.True(t, s.IsPlayerExist("p1"))
}

func TestAddPlayerTwoPlayersSuccess(t *testing.T) {
	s := newTestSession(t)

	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))
	require.NoError(t, s.AddPlayer(newTestPlayer("p2", model.MarkO)))

	assert.Len(t, s.Players, 2)
}

func TestAddPlayerSessionFull(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))
	require.NoError(t, s.AddPlayer(newTestPlayer("p2", model.MarkO)))

	err := s.AddPlayer(newTestPlayer("p3", model.MarkX))
	assert.ErrorIs(t, err, model.ErrSessionFull)
}

func TestAddPlayerDuplicateUUID(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	err := s.AddPlayer(newTestPlayer("p1", model.MarkO))
	assert.ErrorIs(t, err, model.ErrPlayerExists)
}

func TestAddPlayerMarkTaken(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	err := s.AddPlayer(newTestPlayer("p2", model.MarkX))
	assert.ErrorIs(t, err, model.ErrMarkTaken)
}

func TestRemovePlayerFirst(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))
	require.NoError(t, s.AddPlayer(newTestPlayer("p2", model.MarkO)))

	require.NoError(t, s.RemovePlayer("p1"))

	assert.Len(t, s.Players, 1)
	assert.False(t, s.IsPlayerExist("p1"))
	assert.True(t, s.IsPlayerExist("p2"))
}

func TestRemovePlayerSecond(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))
	require.NoError(t, s.AddPlayer(newTestPlayer("p2", model.MarkO)))

	require.NoError(t, s.RemovePlayer("p2"))

	assert.Len(t, s.Players, 1)
	assert.True(t, s.IsPlayerExist("p1"))
	assert.False(t, s.IsPlayerExist("p2"))
}

func TestRemovePlayerOnly(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	require.NoError(t, s.RemovePlayer("p1"))

	assert.Empty(t, s.Players)
}

func TestRemovePlayerNotFound(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))

	err := s.RemovePlayer("unknown")
	assert.ErrorIs(t, err, model.ErrPlayerNotFound)
}

func TestRemovePlayerEmptySession(t *testing.T) {
	s := newTestSession(t)

	err := s.RemovePlayer("p1")
	assert.ErrorIs(t, err, model.ErrPlayerNotFound)
}

func TestGetTurnPlayerEmptySession(t *testing.T) {
	s := newTestSession(t)

	_, err := s.GetTurnPlayer()
	assert.ErrorIs(t, err, model.ErrPlayerNotFound)
}

func TestGetTurnPlayerRotation(t *testing.T) {
	s := newTestSession(t)
	require.NoError(t, s.AddPlayer(newTestPlayer("p1", model.MarkX)))
	require.NoError(t, s.AddPlayer(newTestPlayer("p2", model.MarkO)))

	cases := []struct {
		turn     int
		expected string
	}{
		{0, "p1"},
		{1, "p2"},
		{2, "p1"},
		{3, "p2"},
	}

	for _, tc := range cases {
		s.Turn = tc.turn
		player, err := s.GetTurnPlayer()
		require.NoError(t, err)
		assert.Equal(t, tc.expected, player.UUID, "turn %d", tc.turn)
	}
}
