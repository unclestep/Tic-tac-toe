package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(parent context.Context, databaseURL string) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("new pool: parse config: %w", err)
	}

	conf.MaxConnLifetime = 30 * time.Minute
	conf.MaxConnLifetimeJitter = 15 * time.Minute
	conf.MaxConnIdleTime = 15 * time.Minute
	conf.MaxConns = 20
	conf.MinConns = 5

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("new pool: connect config: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("new pool: ping: %w", err)
	}

	return pool, nil
}
