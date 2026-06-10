package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"tictactoe/internal/application/port"
	dsmodel "tictactoe/internal/infrastructure/storage/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SessionDataSource struct {
	dbtx DBTX
}

func NewSessionDataSource(dbtx DBTX) *SessionDataSource {
	return &SessionDataSource{
		dbtx: dbtx,
	}
}

func (ds *SessionDataSource) Fetch(parent context.Context, uuid string) (*dsmodel.SessionRecord, error) {
	sessionSQL := `
		SELECT uuid, rules_uuid, board, turn, winner, state, seed
		FROM sessions
		WHERE uuid = $1
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	var sr dsmodel.SessionRecord
	var winnerUUID *string
	var binBoard []byte
	err := ds.dbtx.QueryRow(ctx, sessionSQL, uuid).Scan(&sr.UUID, &sr.RulesUUID, &binBoard, &sr.Turn, &winnerUUID, &sr.State, &sr.Seed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("fetch: %w", port.ErrSessionNotFound)
		}
		return nil, fmt.Errorf("fetch: query row: %w", err)
	}

	if err := json.Unmarshal(binBoard, &sr.Board); err != nil {
		return nil, fmt.Errorf("fetch: unmarshal board: %w", err)
	}

	playersSQL := `
		SELECT uuid, name, mark
		FROM players
		WHERE session_uuid = $1
	`

	rows, err := ds.dbtx.Query(ctx, playersSQL, uuid)
	if err != nil {
		return nil, fmt.Errorf("fetch: query players: %w", err)
	}

	for rows.Next() {
		var p dsmodel.PlayerRecord
		if err := rows.Scan(&p.UUID, &p.Name, &p.Mark); err != nil {
			rows.Close()
			return nil, fmt.Errorf("fetch: scan player: %w", err)
		}
		sr.Players = append(sr.Players, &p)

		if winnerUUID != nil && p.UUID == *winnerUUID {
			sr.Winner = &p
		}
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetch: after rows close: %w", err)
	}

	return &sr, nil
}

func (ds *SessionDataSource) Store(parent context.Context, session *dsmodel.SessionRecord) error {
	sessionSQL := `
		INSERT INTO sessions (uuid, rules_uuid, board, turn, winner, state, seed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT(uuid) DO UPDATE
		SET board = EXCLUDED.board,
			turn = EXCLUDED.turn,
			winner = EXCLUDED.winner,
			state = EXCLUDED.state,
			seed = EXCLUDED.seed
	`

	if session.Board == nil || len(session.Board.Cells) == 0 {
		return fmt.Errorf("store: board is empty")
	}

	binBoard, err := json.Marshal(session.Board)
	if err != nil {
		return fmt.Errorf("store: marshal board: %w", err)
	}

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	var winner *string
	if session.Winner != nil {
		winner = &session.Winner.UUID
	}

	_, err = ds.dbtx.Exec(ctx, sessionSQL,
		session.UUID, session.RulesUUID, binBoard,
		session.Turn, winner, session.State, session.Seed,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return fmt.Errorf("store: %w", port.ErrRulesNotFound)
	}
	if err != nil {
		return fmt.Errorf("store: session insert: %w", err)
	}

	if len(session.Players) == 0 {
		return nil
	}

	playerSQL := `
		INSERT INTO players (uuid, session_uuid, name, mark)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(session_uuid, mark) DO UPDATE
		SET uuid = EXCLUDED.uuid,
			name = EXCLUDED.name
	`

	batch := &pgx.Batch{}
	for _, player := range session.Players {
		batch.Queue(playerSQL, player.UUID, session.UUID, player.Name, player.Mark)
	}

	br := ds.dbtx.SendBatch(ctx, batch)
	defer func() { _ = br.Close() }()

	for range session.Players {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("store: batch exec: %w", err)
		}
	}

	return nil
}

func (ds *SessionDataSource) Delete(parent context.Context, UUID string) error {
	sql := `
		DELETE FROM sessions WHERE uuid = $1
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	_, err := ds.dbtx.Exec(ctx, sql, UUID)
	if err != nil {
		return fmt.Errorf("delete: exec: %w", err)
	}

	return nil
}
