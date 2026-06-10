package repository

import (
	"context"
	"fmt"

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

func (r *SessionRepo) Get(ctx context.Context, id string) (*model.Session, error) {
	srecord, err := r.sds.Fetch(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

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

	return dsession, nil
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
