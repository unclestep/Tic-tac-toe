package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tictactoe/internal/application/port"
	dsmodel "tictactoe/internal/infrastructure/storage/model"

	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionDataSource struct {
	pool *pgxpool.Pool
}

func NewSessionDataSource(pool *pgxpool.Pool) *SessionDataSource {
	return &SessionDataSource{
		pool: pool,
	}
}

func (ds *SessionDataSource) Fetch(parent context.Context, uuid string) (*dsmodel.SessionRecord, error) {
	sessionSQL := `
		SELECT s.uuid, s.rules_id, s.board, s.bots, turn, winner, state, seed
		FROM sessions
		WHERE uuid = $1
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	var sr dsmodel.SessionRecord
	var binBoard []byte
	err := ds.pool.QueryRow(ctx, sessionSQL, uuid).Scan(&sr.UUID, &sr.RulesUUID, &binBoard, &sr.Bots, &sr.Turn, &sr.Winner, &sr.State, &sr.Seed)
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
		SELECT uuid, name, mark, bot
		FROM players
		WHERE session_uuid = $1
	`

	rows, err := ds.pool.Query(ctx, playersSQL, uuid)
	if err != nil {
		return nil, fmt.Errorf("fetch: query players: %w", err)
	}

	for rows.Next() {
		var p dsmodel.PlayerRecord
		if err := rows.Scan(&p.UUID, &p.Name, &p.Mark, &p.Bot); err != nil {
			rows.Close()
			return nil, fmt.Errorf("fetch: scan player: %w", err)
		}
		sr.Players = append(sr.Players, &p)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetch: after rows close: %w", err)
	}

	return &sr, nil
}

func (ds *SessionDataSource) Store(parent context.Context, session *dsmodel.SessionRecord) error {
	sessionSQL := `
		INSERT INTO sessions (uuid, rules_id, board, bots, turn, winner, state, seed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT(uuid) DO UPDATE
		SET board = EXCLUDED.board,
			bots = EXCLUDED.bots,
			turn = EXCLUDED.turn,
			winner = EXCLUDED.winner,
			state = EXCLUDED.state,
			seed = EXCLUDED.seed
	`

	if session.Board == nil {
		return fmt.Errorf("store: board is nil")
	}

	binBoard, err := json.Marshal(session.Board)
	if err != nil {
		return fmt.Errorf("store: marshal board: %w", err)
	}

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	_, err = ds.pool.Exec(ctx, sessionSQL,
		session.UUID, session.RulesUUID, binBoard, session.Bots,
		session.Turn, session.Winner, session.State, session.Seed,
	)
	if err != nil {
		return fmt.Errorf("store: session insert: %w", err)
	}

	playerSQL := `
		INSERT INTO players (uuid, session_uuid, name, mark, bot)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(uuid) DO UPDATE
		SET name = EXCLUDED.name,
			mark = EXCLUDED.mark,
			bot = EXCLUDED.bot
	`

	batch := &pgx.Batch{}
	for _, player := range session.Players {
		batch.Queue(playerSQL, player.UUID, session.UUID, player.Name, player.Mark, player.Bot)
	}

	br := ds.pool.SendBatch(ctx, batch)
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

	_, err := ds.pool.Exec(ctx, sql, UUID)
	if err != nil {
		return fmt.Errorf("delete: exec: %w", err)
	}

	return nil
}
