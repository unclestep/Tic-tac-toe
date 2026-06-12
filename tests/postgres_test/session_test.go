package postgres_test

import (
	"context"
	"fmt"

	"tictactoe/internal/domain/model"
	"tictactoe/internal/infrastructure/storage/ds"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
	pg "tictactoe/internal/infrastructure/storage/postgres"
)

func (s *PostgresSuite) savepoint(name string, fn func()) {
	ctx := context.Background()
	_, err := s.tx.Exec(ctx, fmt.Sprintf("SAVEPOINT %s", name))
	s.Require().NoError(err)
	defer func() {
		_, _ = s.tx.Exec(ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT %s", name))
	}()
	fn()
}

func (s *PostgresSuite) prepareSessions() {
	sessions := []*dsmodel.SessionRecord{
		{
			UUID: "A",
			Rules: &dsmodel.RulesRecord{
				BoardWidth: 3, BoardHeight: 3, WinLength: 3,
			},
			Board: &dsmodel.BoardRecord{
				Width:  3,
				Height: 3,
				Cells:  []int8{1, 0, 0, 0, 0, 0, 0, 0, 1},
			},
			Players: []*dsmodel.PlayerRecord{
				{UUID: "A", UserUUID: "BOT", Name: "A", Mark: "X"},
				{UUID: "B", UserUUID: "BOT", Name: "B", Mark: "O"},
			},
			Turn:  0,
			State: "Lobby",
			Seed:  int64(0),
		},
		{
			UUID: "B",
			Rules: &dsmodel.RulesRecord{
				BoardWidth: 2, BoardHeight: 2, WinLength: 2,
			},
			Board: &dsmodel.BoardRecord{
				Width:  2,
				Height: 2,
				Cells:  []int8{0, 0, 0, 0},
			},
			Players: []*dsmodel.PlayerRecord{
				{UUID: "Z", UserUUID: "BOT", Name: "Z", Mark: "X"},
			},
			Turn:  0,
			State: "Playing",
			Seed:  int64(0),
		},
	}

	pds := pg.NewSessionDataSource(s.tx)
	for _, session := range sessions {
		err := pds.Store(context.Background(), session)
		s.Require().NoError(err)
	}
}

func (s *PostgresSuite) TestSessionStore() {
	s.prepareSessions()
	pds := pg.NewSessionDataSource(s.tx)

	s.Run("Default", func() {
		rs, err := pds.Fetch(context.Background())
		s.Require().NoError(err)
		s.Len(rs, 2)
		total := 0
		for _, r := range rs {
			total += len(r.Players)
		}
		s.Equal(3, total)
	})

	s.Run("ExistentSessionFields", func() {
		s.savepoint("existent_session_fields", func() {
			updated := &dsmodel.SessionRecord{
				UUID:  "A",
				Rules: &dsmodel.RulesRecord{BoardWidth: 3, BoardHeight: 3, WinLength: 3},
				Board: &dsmodel.BoardRecord{Width: 1, Height: 1, Cells: []int8{1}},
				Players: []*dsmodel.PlayerRecord{
					{UUID: "A", UserUUID: "BOT", Name: "A", Mark: "X"},
					{UUID: "B", UserUUID: "BOT", Name: "B", Mark: "O"},
				},
				Turn:  -1,
				State: "GameOver",
				Seed:  int64(-1),
			}
			s.Require().NoError(pds.Store(context.Background(), updated))

			rs, err := pds.Fetch(context.Background(), ds.WithUUID("A"))
			s.Require().NoError(err)
			s.Require().Len(rs, 1)
			s.Equal(-1, rs[0].Turn)
			s.Equal("GameOver", rs[0].State)
			s.Equal(int64(-1), rs[0].Seed)
		})
	})

	s.Run("ExistentPlayerSlots", func() {
		s.savepoint("existent_player_slots", func() {
			updated := &dsmodel.SessionRecord{
				UUID:  "A",
				Rules: &dsmodel.RulesRecord{BoardWidth: 3, BoardHeight: 3, WinLength: 3},
				Board: &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1, 0, 0, 0, 0, 0, 0, 0, 1}},
				Players: []*dsmodel.PlayerRecord{
					{UUID: "C", UserUUID: "BOT", Name: "C", Mark: "X"},
					{UUID: "D", UserUUID: "BOT", Name: "D", Mark: "O"},
				},
				Turn: 0, State: "Lobby", Seed: int64(0),
			}
			s.Require().NoError(pds.Store(context.Background(), updated))

			rs, err := pds.Fetch(context.Background(), ds.WithUUID("A"))
			s.Require().NoError(err)
			s.Require().Len(rs, 1)
			s.Require().Len(rs[0].Players, 2)
			for _, p := range rs[0].Players {
				s.NotEqual("A", p.UUID)
				s.NotEqual("B", p.UUID)
			}

			all, err := pds.Fetch(context.Background())
			s.Require().NoError(err)
			total := 0
			for _, r := range all {
				total += len(r.Players)
			}
			s.Equal(3, total)
		})
	})

	s.Run("NilBoard", func() {
		err := pds.Store(context.Background(), &dsmodel.SessionRecord{UUID: "new"})
		s.Error(err)
		rs, err := pds.Fetch(context.Background())
		s.Require().NoError(err)
		s.Len(rs, 2)
	})

	s.Run("EmptyCells", func() {
		s.savepoint("empty_cells", func() {
			err := pds.Store(context.Background(), &dsmodel.SessionRecord{
				UUID:  "new",
				Board: &dsmodel.BoardRecord{Cells: []int8{}},
			})
			s.Error(err)
			rs, err := pds.Fetch(context.Background())
			s.Require().NoError(err)
			s.Len(rs, 2)
		})
	})

	s.Run("BoardRoundTrip", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithUUID("A"))
		s.Require().NoError(err)
		s.Require().Len(rs, 1)
		s.Require().NotNil(rs[0].Board)
		s.Equal(3, rs[0].Board.Width)
		s.Equal(3, rs[0].Board.Height)
		s.Equal([]int8{1, 0, 0, 0, 0, 0, 0, 0, 1}, rs[0].Board.Cells)
	})

	s.Run("NoPlayers", func() {
		s.savepoint("no_players", func() {
			session := &dsmodel.SessionRecord{
				UUID:  "no-players",
				Rules: &dsmodel.RulesRecord{BoardWidth: 3, BoardHeight: 3, WinLength: 3},
				Board: &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1}},
				State: "Lobby",
			}
			s.Require().NoError(pds.Store(context.Background(), session))

			rs, err := pds.Fetch(context.Background(), ds.WithUUID("no-players"))
			s.Require().NoError(err)
			s.Require().Len(rs, 1)
			s.Nil(rs[0].Players)
		})
	})
}

func (s *PostgresSuite) TestSessionFetch() {
	s.prepareSessions()
	pds := pg.NewSessionDataSource(s.tx)

	s.Run("Default", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithUUID("A"))
		s.Require().NoError(err)
		s.Require().Len(rs, 1)
		r := rs[0]
		s.Equal("A", r.UUID)
		s.Require().NotNil(r.Rules)
		s.Equal(3, r.Rules.BoardWidth)
		s.Equal(3, r.Rules.BoardHeight)
		s.Equal(3, r.Rules.WinLength)
		s.Equal(0, r.Turn)
		s.Nil(r.Winner)
		s.Equal("Lobby", r.State)
		s.Equal(int64(0), r.Seed)
		s.Require().NotNil(r.Board)
		s.Equal(3, r.Board.Width)
		s.Equal(3, r.Board.Height)
		s.Equal([]int8{1, 0, 0, 0, 0, 0, 0, 0, 1}, r.Board.Cells)
		s.Require().Equal(2, len(r.Players))
		for _, p := range r.Players {
			switch p.Mark {
			case "X":
				s.Equal("A", p.UUID)
				s.Equal("A", p.Name)
			case "O":
				s.Equal("B", p.UUID)
				s.Equal("B", p.Name)
			default:
				s.Fail("unexpected mark", p.Mark)
			}
		}
	})

	s.Run("Nonexistent", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithUUID("Z"))
		s.NoError(err)
		s.Empty(rs)
	})

	s.Run("NoPlayers", func() {
		s.savepoint("fetch_no_players", func() {
			session := &dsmodel.SessionRecord{
				UUID:  "no-players",
				Rules: &dsmodel.RulesRecord{BoardWidth: 3, BoardHeight: 3, WinLength: 3},
				Board: &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1}},
				State: "Lobby",
			}
			s.Require().NoError(pds.Store(context.Background(), session))

			rs, err := pds.Fetch(context.Background(), ds.WithUUID("no-players"))
			s.Require().NoError(err)
			s.Require().Len(rs, 1)
			s.Nil(rs[0].Players)
		})
	})

	s.Run("CellsRoundTrip", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithUUID("A"))
		s.Require().NoError(err)
		s.Require().Len(rs, 1)
		s.Require().NotNil(rs[0].Board)
		s.Equal([]int8{1, 0, 0, 0, 0, 0, 0, 0, 1}, rs[0].Board.Cells)
	})

	s.Run("NilCellsRoundTrip", func() {
		s.savepoint("fetch_nil_cells", func() {
			_, err := s.tx.Exec(context.Background(), `
				INSERT INTO sessions (uuid, board, turn, winner, state, seed)
				VALUES ('nil-cells', '{"width":3,"height":3,"cells":null}'::jsonb, 0, NULL, 'Lobby', 0)
			`)
			s.Require().NoError(err)
			_, err = s.tx.Exec(context.Background(), `
				INSERT INTO rules (session_uuid, board_width, board_height, win_length)
				VALUES ('nil-cells', 3, 3, 3)
			`)
			s.Require().NoError(err)

			rs, err := pds.Fetch(context.Background(), ds.WithUUID("nil-cells"))
			s.Require().NoError(err)
			s.Require().Len(rs, 1)
			s.Require().NotNil(rs[0].Board)
			s.Nil(rs[0].Board.Cells)
		})
	})

	s.Run("NilBoardRoundTrip", func() {
		s.savepoint("fetch_nil_board", func() {
			_, err := s.tx.Exec(context.Background(), `
				INSERT INTO sessions (uuid, board, turn, winner, state, seed)
				VALUES ('nil-board', 'null'::jsonb, 0, NULL, 'Lobby', 0)
			`)
			s.Require().NoError(err)
			_, err = s.tx.Exec(context.Background(), `
				INSERT INTO rules (session_uuid, board_width, board_height, win_length)
				VALUES ('nil-board', 3, 3, 3)
			`)
			s.Require().NoError(err)

			rs, err := pds.Fetch(context.Background(), ds.WithUUID("nil-board"))
			s.Require().NoError(err)
			s.Require().Len(rs, 1)
			s.Nil(rs[0].Board)
		})
	})

	s.Run("PlayersIsolation", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithUUID("A"))
		s.Require().NoError(err)
		s.Require().Len(rs, 1)
		s.Require().Equal(2, len(rs[0].Players))
		for _, p := range rs[0].Players {
			s.NotEqual("Z", p.UUID)
		}
	})

	s.Run("NoOpts", func() {
		rs, err := pds.Fetch(context.Background())
		s.Require().NoError(err)
		s.Len(rs, 2)
	})

	s.Run("MultipleUUIDs", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithUUID("A", "B"))
		s.Require().NoError(err)
		s.Len(rs, 2)
	})

	s.Run("ByState", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithState(model.StatePlaying))
		s.Require().NoError(err)
		s.Len(rs, 1)
	})

	s.Run("ByStateAndMultipleUUIDs", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithState(model.StateLobby), ds.WithUUID("A", "B"))
		s.Require().NoError(err)
		s.Len(rs, 1)
	})
}

func (s *PostgresSuite) TestSessionDelete() {
	s.prepareSessions()
	pds := pg.NewSessionDataSource(s.tx)

	s.Run("Default", func() {
		s.savepoint("delete_default", func() {
			s.Require().NoError(pds.Delete(context.Background(), "A"))

			rs, err := pds.Fetch(context.Background())
			s.Require().NoError(err)
			s.Require().Len(rs, 1)
			s.NotEqual("A", rs[0].UUID)
		})
	})

	s.Run("Nonexistent", func() {
		s.Require().NoError(pds.Delete(context.Background(), "Z"))
		rs, err := pds.Fetch(context.Background())
		s.Require().NoError(err)
		s.Len(rs, 2)
	})

	s.Run("CascadePlayers", func() {
		s.savepoint("delete_cascade_players", func() {
			s.Require().NoError(pds.Delete(context.Background(), "A"))

			rs, err := pds.Fetch(context.Background())
			s.Require().NoError(err)
			s.Require().Len(rs, 1)
			s.Require().Len(rs[0].Players, 1)
			s.Equal("Z", rs[0].Players[0].UUID)
		})
	})
}
