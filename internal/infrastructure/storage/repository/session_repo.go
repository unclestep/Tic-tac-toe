package repository

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/infrastructure/storage/mapper"
	storagePort "tictactoe/internal/infrastructure/storage/port"
)

type SessionRepo struct {
	ds storagePort.SessionDataSource
}

func NewSessionRepo(ds storagePort.SessionDataSource) *SessionRepo {
	return &SessionRepo{
		ds: ds,
	}
}

func (r *SessionRepo) Get(ctx context.Context, id string) (*model.Session, error) {
	record, err := r.ds.Fetch(ctx, id)
	if err != nil {
		if errors.Is(err, port.ErrSessionNotFound) {
			return nil, fmt.Errorf("get session: %w", port.ErrSessionNotFound)
		}
		return nil, fmt.Errorf("get session: %w", err)
	}
	dsession, err := mapper.ToSessionDomain(record)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}
	return dsession, nil
}

func (r *SessionRepo) Save(ctx context.Context, session *model.Session) error {
	if session == nil {
		return fmt.Errorf("save session: session is nil")
	}
	if session.UUID == "" {
		return fmt.Errorf("save session: empty uuid")
	}
	if err := r.ds.Store(ctx, mapper.ToSessionStorage(session)); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func (r *SessionRepo) Delete(ctx context.Context, id string) error {
	if err := r.ds.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
