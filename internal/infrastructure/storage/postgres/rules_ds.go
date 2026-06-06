package postgres

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/internal/application/port"
	dsmodel "tictactoe/internal/infrastructure/storage/model"

	"github.com/jackc/pgx/v5"
	"time"
)

type RulesDataSource struct {
	dbtx DBTX
}

func NewRulesDataSource(dbtx DBTX) *RulesDataSource {
	return &RulesDataSource{
		dbtx: dbtx,
	}
}

func (ds *RulesDataSource) Fetch(parent context.Context, uuid string) (*dsmodel.RulesRecord, error) {
	sql := `
		SELECT uuid, board_width, board_height, win_length
		FROM rules
		WHERE uuid = $1
	`
	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	var rulesRecord dsmodel.RulesRecord
	err := ds.dbtx.QueryRow(ctx, sql, uuid).Scan(&rulesRecord.UUID, &rulesRecord.BoardWidth, &rulesRecord.BoardHeight, &rulesRecord.WinLength)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("fetch: %w", port.ErrRulesNotFound)
		}
		return nil, fmt.Errorf("fetch: query row: %w", err)
	}

	return &rulesRecord, nil
}

func (ds *RulesDataSource) Store(parent context.Context, record *dsmodel.RulesRecord) error {
	sql := `
		INSERT INTO rules (uuid, board_width, board_height, win_length)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(uuid) DO NOTHING
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	_, err := ds.dbtx.Exec(ctx, sql, record.UUID, record.BoardWidth, record.BoardHeight, record.WinLength)
	if err != nil {
		return fmt.Errorf("store: exec: %w", err)
	}
	return nil
}

func (ds *RulesDataSource) Delete(parent context.Context, uuid string) error {
	sql := `
		DELETE FROM rules WHERE uuid = $1
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	_, err := ds.dbtx.Exec(ctx, sql, uuid)
	if err != nil {
		return fmt.Errorf("delete: exec: %w", err)
	}

	return nil
}
