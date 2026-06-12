package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"tictactoe/internal/application/port"
	"tictactoe/internal/infrastructure/storage/ds"
	dsmodel "tictactoe/internal/infrastructure/storage/model"

	"github.com/jackc/pgx/v5"
)

type SessionDataSource struct {
	dbtx DBTX
}

func NewSessionDataSource(dbtx DBTX) *SessionDataSource {
	return &SessionDataSource{
		dbtx: dbtx,
	}
}

func (s *SessionDataSource) Fetch(parent context.Context, opts ...ds.SessionOpt) ([]*dsmodel.SessionRecord, error) {
	cfg := &ds.SessionConfig{}

	for _, opt := range opts {
		opt.ApplyToSession(cfg)
	}

	sessionSQL := `
		SELECT s.uuid,
			   r.board_width, r.board_height, r.win_length,
			   p.uuid, p.name, p.mark,
			   s.board, s.turn, s.winner, s.state, s.seed
		FROM sessions s
		LEFT JOIN players p ON p.session_uuid = s.uuid
		JOIN rules r ON r.session_uuid = s.uuid
		WHERE 1=1
	`

	var args []any

	if cfg.State != nil {
		args = append(args, *cfg.State)
		sessionSQL += fmt.Sprintf(" AND s.state = $%d", len(args))
	}

	if cfg.UUIDs != nil {
		args = append(args, cfg.UUIDs)
		sessionSQL += fmt.Sprintf(" AND s.uuid = ANY($%d)", len(args))
	}

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	rows, err := s.dbtx.Query(ctx, sessionSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("fetch: query: %w", err)
	}
	defer rows.Close()

	recordMap := make(map[string]*dsmodel.SessionRecord)
	var order []string

	for rows.Next() {
		var s dsmodel.SessionRecord
		var r dsmodel.RulesRecord
		var pUUID, pName, pMark *string
		var winnerUUID *string
		var binBoard []byte

		if err := rows.Scan(&s.UUID,
			&r.BoardWidth, &r.BoardHeight, &r.WinLength,
			&pUUID, &pName, &pMark,
			&binBoard, &s.Turn, &winnerUUID, &s.State, &s.Seed); err != nil {
			return nil, fmt.Errorf("fetch: %w", err)
		}

		if err := json.Unmarshal(binBoard, &s.Board); err != nil {
			return nil, fmt.Errorf("fetch: unmarshal board: %w", err)
		}

		if pUUID != nil && pName != nil && pMark != nil {
			s.Players = append(s.Players, &dsmodel.PlayerRecord{
				UUID: *pUUID,
				Name: *pName,
				Mark: *pMark,
			})
		}

		if prev, ok := recordMap[s.UUID]; ok {
			prev.Players = append(prev.Players, s.Players...)
		} else {
			s.Rules = &r
			recordMap[s.UUID] = &s
			order = append(order, s.UUID)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fetch: after srows close: %w", err)
	}

	records := make([]*dsmodel.SessionRecord, len(order))

	for i, uuid := range order {
		records[i] = recordMap[uuid]
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("fetch: %w", port.ErrSessionNotFound)
	}

	return records, nil
}

func (s *SessionDataSource) Store(parent context.Context, session *dsmodel.SessionRecord) error {
	if session.Board == nil || len(session.Board.Cells) == 0 {
		return fmt.Errorf("store: board is empty")
	}

	binBoard, err := json.Marshal(session.Board)
	if err != nil {
		return fmt.Errorf("store: marshal board: %w", err)
	}

	var winner *string
	if session.Winner != nil {
		winner = &session.Winner.UUID
	}

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	tx, err := s.dbtx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("store: begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	batch.Queue(`
		INSERT INTO sessions(uuid, board, turn, winner, state, seed)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT(uuid) DO UPDATE
		SET board = EXCLUDED.board,
			turn = EXCLUDED.turn,
			winner = EXCLUDED.winner,
			state = EXCLUDED.state,
			seed = EXCLUDED.seed;
	`, session.UUID, binBoard, session.Turn, winner, session.State, session.Seed)

	batch.Queue(`
		INSERT INTO rules(session_uuid, board_width, board_height, win_length)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(session_uuid) DO NOTHING
	`, session.UUID, session.Rules.BoardWidth, session.Rules.BoardHeight, session.Rules.WinLength)

	for _, p := range session.Players {
		batch.Queue(`INSERT INTO players(uuid, session_uuid, user_uuid, name, mark)
					VALUES ($1, $2, $3, $4, $5)
					ON CONFLICT(uuid) DO NOTHING
		`, p.UUID, session.UUID, p.UserUUID, p.Name, p.Mark)
	}

	br := tx.SendBatch(ctx, batch)

	for range batch.Len() {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return fmt.Errorf("store: batch exec: %w", err)
		}
	}

	if err := br.Close(); err != nil {
		return fmt.Errorf("store: close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		_ = br.Close()
		return fmt.Errorf("store: commit: %w", err)
	}

	return nil
}

func (s *SessionDataSource) Delete(parent context.Context, UUID string) error {
	sql := `
		DELETE FROM sessions WHERE uuid = $1
	`

	ctx, cancel := context.WithTimeout(parent, 15*time.Second)
	defer cancel()

	_, err := s.dbtx.Exec(ctx, sql, UUID)
	if err != nil {
		return fmt.Errorf("delete: exec: %w", err)
	}

	return nil
}
