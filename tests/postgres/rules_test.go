package postgres_test

import (
	"context"
	"errors"

	appPort "tictactoe/internal/application/port"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
	pg "tictactoe/internal/infrastructure/storage/postgres"
)

func (s *PostgresSuite) prepareRules() {
	records := []*dsmodel.RulesRecord{
		{
			UUID:        "A",
			BoardWidth:  3,
			BoardHeight: 3,
			WinLength:   3,
		},
		{
			UUID:        "B",
			BoardWidth:  6,
			BoardHeight: 6,
			WinLength:   6,
		},
		{
			UUID:        "C",
			BoardWidth:  9,
			BoardHeight: 9,
			WinLength:   9,
		},
	}

	ds := pg.NewRulesDataSource(s.tx)
	for _, r := range records {
		err := ds.Store(context.Background(), r)
		s.Require().NoError(err)
	}
}

func (s *PostgresSuite) getAllRules() []*dsmodel.RulesRecord {
	sql := `
		SELECT uuid, board_width, board_height, win_length
		FROM rules
	`

	rows, err := s.tx.Query(context.Background(), sql)
	s.NoError(err)
	defer rows.Close()

	var rr []*dsmodel.RulesRecord
	for rows.Next() {
		var r dsmodel.RulesRecord
		err := rows.Scan(&r.UUID, &r.BoardWidth, &r.BoardHeight, &r.WinLength)
		s.Require().NoError(err)
		rr = append(rr, &r)
	}

	if err := rows.Err(); err != nil {
		s.Require().NoError(err)
	}

	return rr
}

func (s *PostgresSuite) TestRulesStore() {
	s.prepareRules()
	ds := pg.NewRulesDataSource(s.tx)

	s.Run("Default", func() {
		rr := s.getAllRules()
		s.Equal(3, len(rr))
	})
	s.Run("Update existed", func() {
		u := &dsmodel.RulesRecord{
			UUID:       "A",
			BoardWidth: 2,
			WinLength:  -1,
		}
		err := ds.Store(context.Background(), u)
		s.NoError(err)
		rr := s.getAllRules()
		s.Equal(3, len(rr))
		for _, r := range rr {
			if r.UUID == "A" {
				s.Equal(3, r.BoardWidth)
				s.Equal(3, r.BoardHeight)
				s.Equal(3, r.WinLength)
			}
		}
	})
}

func (s *PostgresSuite) TestRulesFetch() {
	s.prepareRules()
	ds := pg.NewRulesDataSource(s.tx)

	s.Run("Default", func() {
		r, err := ds.Fetch(context.Background(), "A")
		s.NoError(err)
		s.Equal("A", r.UUID)
		s.Equal(3, r.BoardWidth)
		s.Equal(3, r.BoardHeight)
		s.Equal(3, r.WinLength)
	})
	s.Run("Nonexistent", func() {
		r, err := ds.Fetch(context.Background(), "Z")
		s.Error(err)
		s.True(errors.Is(err, appPort.ErrRulesNotFound))
		s.Nil(r)
	})
}

func (s *PostgresSuite) TestRulesDelete() {
	s.prepareRules()
	ds := pg.NewRulesDataSource(s.tx)

	s.Run("Default", func() {
		err := ds.Delete(context.Background(), "A")
		s.NoError(err)
		rr := s.getAllRules()
		s.Equal(2, len(rr))
		for _, r := range rr {
			s.NotEqual("A", r.UUID)
		}
	})
	s.Run("Noexistent", func() {
		err := ds.Delete(context.Background(), "Z")
		s.NoError(err)
		rr := s.getAllRules()
		s.Equal(2, len(rr))
	})
}
