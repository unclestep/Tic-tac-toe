package repository

import (
	"context"
	"fmt"

	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/infrastructure/storage/ds"
	"tictactoe/internal/infrastructure/storage/mapper"
)

type SessionRepo struct {
	sds ds.SessionDataSource
}

func NewSessionRepo(sds ds.SessionDataSource) *SessionRepo {
	return &SessionRepo{
		sds: sds,
	}
}

func SessionDomainOptsToDatasourceOpts(c *port.SessionGetConfig) []ds.SessionOpt {
	var opts []ds.SessionOpt
	if c.UUIDs != nil {
		opts = append(opts, ds.WithUUID(c.UUIDs...))
	}
	if c.State != nil {
		opts = append(opts, ds.WithState(*c.State))
	}
	return opts
}

func (r *SessionRepo) Get(ctx context.Context, opts ...port.SessionGetOpt) ([]*model.Session, error) {
	cfg := &port.SessionGetConfig{}
	for _, opt := range opts {
		opt.ApplyToSession(cfg)
	}

	records, err := r.sds.Fetch(ctx, SessionDomainOptsToDatasourceOpts(cfg)...)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	sessions := make([]*model.Session, len(records))

	for i, record := range records {
		session, err := mapper.ToSessionDomain(record)
		if err != nil {
			return nil, fmt.Errorf("get session: %w", err)
		}
		sessions[i] = session
	}

	return sessions, nil
}

func (r *SessionRepo) Save(ctx context.Context, session *model.Session) error {
	if session == nil {
		return fmt.Errorf("save session: session is nil")
	}
	if session.UUID == "" {
		return fmt.Errorf("save session: empty uuid")
	}
	if err := r.sds.Store(ctx, mapper.ToSessionStorage(session)); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return nil
}

func (r *SessionRepo) Delete(ctx context.Context, id string) error {
	if err := r.sds.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
