package usecase_test

import (
	"context"
	"testing"

	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSessionRepo struct {
	mock.Mock
}

func (m *mockSessionRepo) Get(ctx context.Context, uuid string) (*model.Session, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Session), args.Error(1)
}

func (m *mockSessionRepo) Save(ctx context.Context, session *model.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

type mockRulesRepo struct {
	mock.Mock
}

func (m *mockRulesRepo) Get(ctx context.Context, uuid string) (*model.Rules, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Rules), args.Error(1)
}

func (m *mockRulesRepo) Save(ctx context.Context, rules *model.Rules) error {
	args := m.Called(ctx, rules)
	return args.Error(0)
}

type mockBotMover struct {
	mock.Mock
}

func (m *mockBotMover) MakeMove(session *model.Session, botMark model.Mark) error {
	args := m.Called(session, botMark)
	return args.Error(0)
}

type mockHumanMover struct {
	mock.Mock
}

func (m *mockHumanMover) MakeMove(session *model.Session, player *model.Player, p geometry.Point) error {
	args := m.Called(session, player, p)
	return args.Error(0)
}

func newSessionWithPlayers(t *testing.T, players ...*model.Player) *model.Session {
	t.Helper()
	rules := &model.Rules{UUID: "rules-1", BoardWidth: 3, BoardHeight: 3, WinLength: 3}
	board := model.NewBoard(3, 3)
	s := model.NewSession("session-1", &model.SessionParams{}, rules, board)
	for _, p := range players {
		require.NoError(t, s.AddPlayer(p))
	}
	return s
}
