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
	rds ds.RulesDataSource
}

func NewSessionRepo(sds ds.SessionDataSource, rds ds.RulesDataSource) *SessionRepo {
	return &SessionRepo{
		sds: sds,
		rds: rds,
	}
}

func (r *SessionRepo) Get(ctx context.Context, opts ...port.GetOpt) ([]*model.Session, error) {
	cfg := &port.GetConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	var dsOpts []ds.FetchOption
	if cfg.State != nil {
		dsOpts = append(dsOpts, ds.WithState(*cfg.State))
	}
	if cfg.UUID != nil {
		dsOpts = append(dsOpts, ds.WithUUID(*cfg.UUID))
	}

	dsrecords, err := r.sds.Fetch(ctx, dsOpts...)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	var dsessions []*model.Session

	for _, srecord := range dsrecords {
		dsession, err := mapper.ToSessionDomain(srecord)
		if err != nil {
			return nil, fmt.Errorf("get session: %w", err)
		}

		rrecord, err := r.rds.Fetch(ctx, srecord.RulesUUID)
		if err != nil {
			return nil, fmt.Errorf("get session: %w", err)
		}
		drules := mapper.ToRulesDomain(rrecord)
		dsession.Rules = drules

		dsessions = append(dsessions, dsession)
	}

	return dsessions, nil
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
