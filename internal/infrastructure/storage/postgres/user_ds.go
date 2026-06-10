package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tictactoe/internal/application/port"
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

func (ds *UserDataSource) Fetch(parent context.Context, login string) (*model.UserRecord, error) {
	sql := `
		SELECT uuid, login, password
		FROM users
		WHERE login = $1
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	var userRecord model.UserRecord
	err := ds.dbtx.QueryRow(ctx, sql, login).Scan(&userRecord.UUID, &userRecord.Login, &userRecord.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("fetch: %w", port.ErrUserNotFound)
		}
		return nil, fmt.Errorf("fetch: query row: %w", err)
	}

	return &userRecord, nil
}

func (ds *UserDataSource) Store(parent context.Context, user *model.UserRecord) error {
	sql := `
		INSERT INTO users (uuid, login, password)
		VALUES ($1, $2, $3)
		ON CONFLICT(uuid) DO UPDATE
		SET login = EXCLUDED.login,
			password = EXCLUDED.password
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	_, err := ds.dbtx.Exec(ctx, sql, user.UUID, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("store: user insert: %w", err)
	}

	return nil
}
