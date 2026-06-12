package postgres

import (
	"context"
	"fmt"
	"time"

	"tictactoe/internal/application/port"
	"tictactoe/internal/infrastructure/storage/ds"
	"tictactoe/internal/infrastructure/storage/model"

	"github.com/jackc/pgx/v5"
)

type UserDataSource struct {
	dbtx DBTX
}

func NewUserDataSource(dbtx DBTX) *UserDataSource {
	return &UserDataSource{
		dbtx: dbtx,
	}
}

func (u *UserDataSource) Fetch(parent context.Context, opts ...ds.UserOpt) ([]*model.UserRecord, error) {
	sql := `
		SELECT uuid, login, password
		FROM users
		WHERE 1=1
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	cfg := &ds.UserConfig{}
	for _, opt := range opts {
		opt.ApplyToUser(cfg)
	}

	var args []any

	if cfg.UUIDs != nil {
		args = append(args, cfg.UUIDs)
		sql += fmt.Sprintf(" AND uuid = ANY($%d)", len(args))
	}

	if cfg.Logins != nil {
		args = append(args, cfg.Logins)
		sql += fmt.Sprintf(" AND login = ANY($%d)", len(args))
	}

	rows, err := u.dbtx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("fetch users: %w", err)
	}

	records, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[model.UserRecord])
	if err != nil {
		return nil, fmt.Errorf("fetch: collect rows: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("fetch: %w", port.ErrUserNotFound)
	}

	return records, nil
}

func (u *UserDataSource) Store(parent context.Context, user *model.UserRecord) error {
	sql := `
		INSERT INTO users (uuid, login, password)
		VALUES ($1, $2, $3)
		ON CONFLICT(uuid) DO UPDATE
		SET login = EXCLUDED.login,
			password = EXCLUDED.password
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	_, err := u.dbtx.Exec(ctx, sql, user.UUID, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("store: user insert: %w", err)
	}

	return nil
}
