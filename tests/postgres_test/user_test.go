package postgres_test

import (
	"context"

	"tictactoe/internal/infrastructure/storage/ds"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
	pg "tictactoe/internal/infrastructure/storage/postgres"
)

func (s *PostgresSuite) prepareUsers() {
	users := []*dsmodel.UserRecord{
		{UUID: "1", Login: "user1", Password: "pass1"},
		{UUID: "2", Login: "user2", Password: "pass2"},
	}

	pds := pg.NewUserDataSource(s.tx)
	for _, u := range users {
		s.Require().NoError(pds.Store(context.Background(), u))
	}
}

func (s *PostgresSuite) TestUserFetch() {
	s.prepareUsers()
	pds := pg.NewUserDataSource(s.tx)

	s.Run("NoOpts", func() {
		rs, err := pds.Fetch(context.Background())
		s.Require().NoError(err)
		s.Len(rs, 3) // BOT + 2 test users
	})

	s.Run("ByUUID", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithUUID("1"))
		s.Require().NoError(err)
		s.Len(rs, 1)
	})

	s.Run("MultipleUUIDs", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithUUID("1", "2"))
		s.Require().NoError(err)
		s.Len(rs, 2)
	})

	s.Run("ByLogin", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithLogin("user1"))
		s.Require().NoError(err)
		s.Len(rs, 1)
	})

	s.Run("ByLoginAndMultipleUUIDs", func() {
		rs, err := pds.Fetch(context.Background(), ds.WithLogin("user1", "user2"), ds.WithUUID("1"))
		s.Require().NoError(err)
		s.Len(rs, 1)
	})
}
