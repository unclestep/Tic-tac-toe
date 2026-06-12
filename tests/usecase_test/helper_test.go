package usecase_test

import (
	"context"
	"testing"

	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	"tictactoe/pkg/geometry"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockSessionRepo struct {
	mock.Mock
}

func (m *mockSessionRepo) Get(ctx context.Context, opts ...port.SessionGetOpt) ([]*model.Session, error) {
	cfg := &port.SessionGetConfig{}
	for _, o := range opts {
		o.ApplyToSession(cfg)
	}
	uuid := ""
	if len(cfg.UUIDs) > 0 {
		uuid = cfg.UUIDs[0]
	}
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return []*model.Session{args.Get(0).(*model.Session)}, args.Error(1)
}

func (m *mockSessionRepo) Save(ctx context.Context, session *model.Session) error {
	args := m.Called(ctx, session)
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

func newPlayer(uuid string, mark model.Mark) *model.Player {
	return model.NewPlayer(uuid, uuid, uuid, mark)
}

func newSessionWithPlayers(t *testing.T, players ...*model.Player) *model.Session {
	t.Helper()
	board := model.NewBoard(3, 3)
	s := model.NewSession("1", &model.SessionParams{}, model.NewDefaultRules(), board)
	for _, p := range players {
		require.NoError(t, s.AddPlayer(p))
	}
	return s
}
