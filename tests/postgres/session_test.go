package postgres_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	appPort "tictactoe/internal/application/port"
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
			UUID:      "A",
			RulesUUID: "A",
			Board: &dsmodel.BoardRecord{
				Width:  3,
				Height: 3,
				Cells:  []int8{1, 0, 0, 0, 0, 0, 0, 0, 1},
			},
			Players: []*dsmodel.PlayerRecord{
				{UUID: "A", Name: "A", Mark: "X", Bot: false},
				{UUID: "B", Name: "B", Mark: "O", Bot: false},
			},
			Bots:   0,
			Turn:   0,
			Winner: "",
			State:  "Lobby",
			Seed:   int64(0),
		},
		{
			UUID:      "B",
			RulesUUID: "B",
			Board: &dsmodel.BoardRecord{
				Width:  2,
				Height: 2,
				Cells:  []int8{0, 0, 0, 0},
			},
			Players: []*dsmodel.PlayerRecord{
				{UUID: "Z", Name: "Z", Mark: "X", Bot: false},
			},
			Bots:   1,
			Turn:   0,
			Winner: "A",
			State:  "",
			Seed:   int64(0),
		},
	}

	ds := pg.NewSessionDataSource(s.tx)
	for _, session := range sessions {
		err := ds.Store(context.Background(), session)
		s.Require().NoError(err)
	}
}

func (s *PostgresSuite) getAllSessions() []*dsmodel.SessionRecord {
	sql := `
		SELECT uuid, rules_uuid, board, bots, turn, winner, state, seed
		FROM sessions
	`

	rows, err := s.tx.Query(context.Background(), sql)
	s.Require().NoError(err)
	defer rows.Close()

	var ss []*dsmodel.SessionRecord

	for rows.Next() {
		var session dsmodel.SessionRecord
		var binBoard []byte
		err := rows.Scan(
			&session.UUID, &session.RulesUUID,
			&binBoard, &session.Bots, &session.Turn, &session.Winner,
			&session.State, &session.Seed,
		)
		s.Require().NoError(err)

		err = json.Unmarshal(binBoard, &session.Board)
		s.Require().NoError(err)
		ss = append(ss, &session)
	}

	if err := rows.Err(); err != nil {
		s.Require().NoError(err)
	}

	return ss
}

func (s *PostgresSuite) getAllPlayers() []*dsmodel.PlayerRecord {
	sql := `
		SELECT uuid, session_uuid, name, mark, bot
		FROM players
	`

	rows, err := s.tx.Query(context.Background(), sql)
	s.Require().NoError(err)
	defer rows.Close()

	var players []*dsmodel.PlayerRecord

	for rows.Next() {
		var p dsmodel.PlayerRecord
		err := rows.Scan(&p.UUID, nil, &p.Name, &p.Mark, &p.Bot)
		s.Require().NoError(err)
		players = append(players, &p)
	}

	if err := rows.Err(); err != nil {
		s.Require().NoError(err)
	}

	return players
}

func (s *PostgresSuite) getPlayersBySession(sessionUUID string) []*dsmodel.PlayerRecord {
	sql := `
		SELECT uuid, name, mark, bot
		FROM players
		WHERE session_uuid = $1
	`

	rows, err := s.tx.Query(context.Background(), sql, sessionUUID)
	s.Require().NoError(err)
	defer rows.Close()

	var players []*dsmodel.PlayerRecord

	for rows.Next() {
		var p dsmodel.PlayerRecord
		err := rows.Scan(&p.UUID, &p.Name, &p.Mark, &p.Bot)
		s.Require().NoError(err)
		players = append(players, &p)
	}

	return players
}

func (s *PostgresSuite) TestSessionStore() {
	s.prepareRules()
	s.prepareSessions()
	ds := pg.NewSessionDataSource(s.tx)

	s.Run("Default", func() {
		sessions := s.getAllSessions()
		players := s.getAllPlayers()
		s.Equal(2, len(sessions))
		s.Equal(3, len(players))
	})

	s.Run("ExistentSessionFields", func() {
		s.savepoint("existent_session_fields", func() {
			updated := &dsmodel.SessionRecord{
				UUID:      "A",
				RulesUUID: "A",
				Board:     &dsmodel.BoardRecord{Width: 1, Height: 1, Cells: []int8{1}},
				Players: []*dsmodel.PlayerRecord{
					{UUID: "A", Name: "A", Mark: "X", Bot: false},
					{UUID: "B", Name: "B", Mark: "O", Bot: false},
				},
				Bots:   -1,
				Turn:   -1,
				Winner: "Someone",
				State:  "GameOver",
				Seed:   int64(-1),
			}
			err := ds.Store(context.Background(), updated)
			s.Require().NoError(err)

			sessions := s.getAllSessions()
			for _, session := range sessions {
				if session.UUID == "A" {
					s.Equal(-1, session.Bots)
					s.Equal(-1, session.Turn)
					s.Equal("Someone", session.Winner)
					s.Equal("GameOver", session.State)
					s.Equal(int64(-1), session.Seed)
				}
			}
		})
	})

	s.Run("ExistentPlayerSlots", func() {
		s.savepoint("existent_player_slots", func() {
			updated := &dsmodel.SessionRecord{
				UUID:      "A",
				RulesUUID: "A",
				Board:     &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1, 0, 0, 0, 0, 0, 0, 0, 1}},
				Players: []*dsmodel.PlayerRecord{
					{UUID: "C", Name: "C", Mark: "X", Bot: false},
					{UUID: "D", Name: "D", Mark: "O", Bot: false},
				},
				Bots: 0, Turn: 0, Winner: "", State: "Lobby", Seed: int64(0),
			}
			err := ds.Store(context.Background(), updated)
			s.Require().NoError(err)

			players := s.getPlayersBySession("A")
			s.Equal(2, len(players))
			for _, p := range players {
				s.NotEqual("A", p.UUID)
				s.NotEqual("B", p.UUID)
			}

			s.Equal(3, len(s.getAllPlayers()))
		})
	})

	s.Run("NilBoard", func() {
		err := ds.Store(context.Background(), &dsmodel.SessionRecord{UUID: "new"})
		s.Error(err)
		s.Equal(2, len(s.getAllSessions()))
	})

	s.Run("EmptyCells", func() {
		s.savepoint("empty_cells", func() {
			err := ds.Store(context.Background(), &dsmodel.SessionRecord{
				UUID:      "new",
				RulesUUID: "A",
				Board: &dsmodel.BoardRecord{
					Cells: []int8{},
				},
			})
			s.Error(err)
			s.Equal(2, len(s.getAllSessions()))
		})
	})

	s.Run("BoardRoundTrip", func() {
		sessions := s.getAllSessions()
		for _, session := range sessions {
			if session.UUID == "A" {
				s.Require().NotNil(session.Board)
				s.Equal(3, session.Board.Width)
				s.Equal(3, session.Board.Height)
				s.Equal([]int8{1, 0, 0, 0, 0, 0, 0, 0, 1}, session.Board.Cells)
			}
		}
	})

	s.Run("BotThenHuman", func() {
		s.savepoint("bot_then_human", func() {
			withBot := &dsmodel.SessionRecord{
				UUID:      "A",
				RulesUUID: "A",
				Board:     &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1, 0, 0, 0, 0, 0, 0, 0, 1}},
				Players: []*dsmodel.PlayerRecord{
					{UUID: "BOT_1", Name: "Bot", Mark: "X", Bot: true},
					{UUID: "B", Name: "B", Mark: "O", Bot: false},
				},
				Bots: 1, Turn: 0, Winner: "", State: "Lobby", Seed: 0,
			}
			s.Require().NoError(ds.Store(context.Background(), withBot))

			withHuman := &dsmodel.SessionRecord{
				UUID:      "A",
				RulesUUID: "A",
				Board:     &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1, 0, 0, 0, 0, 0, 0, 0, 1}},
				Players: []*dsmodel.PlayerRecord{
					{UUID: "HUMAN_2", Name: "Human", Mark: "X", Bot: false},
					{UUID: "B", Name: "B", Mark: "O", Bot: false},
				},
				Bots: 0, Turn: 0, Winner: "", State: "Lobby", Seed: 0,
			}
			s.Require().NoError(ds.Store(context.Background(), withHuman))

			players := s.getPlayersBySession("A")
			s.Equal(2, len(players))
			for _, p := range players {
				if p.Mark == "X" {
					s.Equal("HUMAN_2", p.UUID)
					s.Equal(false, p.Bot)
				}
			}
		})
	})

	s.Run("InvalidRulesReference", func() {
		s.savepoint("invalid_rules_reference", func() {
			session := &dsmodel.SessionRecord{
				UUID:      "orphan",
				RulesUUID: "nonexistent",
				Board:     &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1}},
				Players:   []*dsmodel.PlayerRecord{},
				State:     "Lobby",
			}
			err := ds.Store(context.Background(), session)
			s.Require().Error(err)
			s.True(errors.Is(err, appPort.ErrRulesNotFound))
		})
	})

	s.Run("NoPlayers", func() {
		s.savepoint("no_players", func() {
			session := &dsmodel.SessionRecord{
				UUID:      "no-players",
				RulesUUID: "A",
				Board:     &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1}},
				Players:   nil,
				State:     "Lobby",
			}
			err := ds.Store(context.Background(), session)
			s.Require().NoError(err)
			s.Equal(0, len(s.getPlayersBySession("no-players")))
		})
	})
}

func (s *PostgresSuite) TestSessionFetch() {
	s.prepareRules()
	s.prepareSessions()
	ds := pg.NewSessionDataSource(s.tx)

	s.Run("Default", func() {
		r, err := ds.Fetch(context.Background(), "A")
		s.Require().NoError(err)
		s.Equal("A", r.UUID)
		s.Equal("A", r.RulesUUID)
		s.Equal(0, r.Bots)
		s.Equal(0, r.Turn)
		s.Equal("", r.Winner)
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
				s.Equal(false, p.Bot)
			case "O":
				s.Equal("B", p.UUID)
				s.Equal("B", p.Name)
				s.Equal(false, p.Bot)
			default:
				s.Fail("unexpected mark", p.Mark)
			}
		}
	})

	s.Run("Nonexistent", func() {
		r, err := ds.Fetch(context.Background(), "Z")
		s.Error(err)
		s.True(errors.Is(err, appPort.ErrSessionNotFound))
		s.Nil(r)
	})

	s.Run("NoPlayers", func() {
		s.savepoint("fetch_no_players", func() {
			session := &dsmodel.SessionRecord{
				UUID:      "no-players",
				RulesUUID: "A",
				Board:     &dsmodel.BoardRecord{Width: 3, Height: 3, Cells: []int8{1}},
				Players:   nil,
				State:     "Lobby",
			}
			s.Require().NoError(ds.Store(context.Background(), session))

			r, err := ds.Fetch(context.Background(), "no-players")
			s.Require().NoError(err)
			s.Nil(r.Players)
		})
	})

	s.Run("CellsRoundTrip", func() {
		r, err := ds.Fetch(context.Background(), "A")
		s.Require().NoError(err)
		s.Require().NotNil(r.Board)
		s.Equal([]int8{1, 0, 0, 0, 0, 0, 0, 0, 1}, r.Board.Cells)
	})

	s.Run("NilCellsRoundTrip", func() {
		s.savepoint("fetch_nil_cells", func() {
			_, err := s.tx.Exec(context.Background(), `
            INSERT INTO sessions (uuid, rules_uuid, board, bots, turn, winner, state, seed)
            VALUES ('nil-cells', 'A', '{"width":3,"height":3,"cells":null}'::jsonb, 0, 0, '', 'Lobby', 0)
        `)
			s.Require().NoError(err)

			r, err := ds.Fetch(context.Background(), "nil-cells")
			s.Require().NoError(err)
			s.Require().NotNil(r.Board)
			s.Nil(r.Board.Cells)
		})
	})

	s.Run("NilBoardRoundTrip", func() {
		s.savepoint("fetch_nil_board", func() {
			_, err := s.tx.Exec(context.Background(), `
				INSERT INTO sessions (uuid, rules_uuid, board, bots, turn, winner, state, seed)
				VALUES ('nil-board', 'A', 'null'::jsonb, 0, 0, '', 'Lobby', 0)
			`)
			s.Require().NoError(err)

			r, err := ds.Fetch(context.Background(), "nil-board")
			s.Require().NoError(err)
			s.Nil(r.Board)
		})
	})

	s.Run("PlayersIsolation", func() {
		r, err := ds.Fetch(context.Background(), "A")
		s.Require().NoError(err)
		s.Require().Equal(2, len(r.Players))
		for _, p := range r.Players {
			s.NotEqual("Z", p.UUID)
		}
	})
}

func (s *PostgresSuite) TestSessionDelete() {
	s.prepareRules()
	s.prepareSessions()
	ds := pg.NewSessionDataSource(s.tx)

	s.Run("Default", func() {
		s.savepoint("delete_default", func() {
			err := ds.Delete(context.Background(), "A")
			s.Require().NoError(err)
			sessions := s.getAllSessions()
			s.Equal(1, len(sessions))
			for _, session := range sessions {
				s.NotEqual("A", session.UUID)
			}
		})
	})

	s.Run("Nonexistent", func() {
		err := ds.Delete(context.Background(), "Z")
		s.Require().NoError(err)
		s.Equal(2, len(s.getAllSessions()))
	})

	s.Run("CascadePlayers", func() {
		s.savepoint("delete_cascade_players", func() {
			err := ds.Delete(context.Background(), "A")
			s.Require().NoError(err)

			players := s.getAllPlayers()
			s.Equal(1, len(players))
			s.Equal("Z", players[0].UUID)
		})
	})
}
