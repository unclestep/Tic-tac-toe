package postgres_test

import (
	"context"
	"testing"

	pg "tictactoe/internal/infrastructure/storage/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	pool      *pgxpool.Pool
	tx        pgx.Tx
}

func (s *PostgresSuite) SetupSuite() {
	ctx := context.Background()
	container, err := postgres.Run(
		ctx, "postgres:18.3-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)

	s.Require().NoError(err)
	s.container = container

	databaseURL, err := container.ConnectionString(ctx, "sslmode=disable")
	s.Require().NoError(err)

	err = pg.Migrate(ctx, databaseURL)
	s.Require().NoError(err)

	pool, err := pg.NewPool(ctx, databaseURL)
	s.Require().NoError(err)
	s.pool = pool
}

func (s *PostgresSuite) TearDownSuite() {
	s.pool.Close()
	err := s.container.Terminate(context.Background())
	s.Require().NoError(err)
}

func (s *PostgresSuite) SetupTest() {
	tx, err := s.pool.Begin(context.Background())
	s.Require().NoError(err)
	s.tx = tx
}

func (s *PostgresSuite) TearDownTest() {
	err := s.tx.Rollback(context.Background())
	s.Require().NoError(err)
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
