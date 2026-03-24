package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	mapper "tictactoe/internal/infrastructure/storage/mapper"
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
		return nil, fmt.Errorf("%w: %s", port.ErrSessionNotFound, id)
	}
	return mapper.ToSessionDomain(record), nil
}

func (r *SessionRepo) Save(ctx context.Context, session *model.Session) error {
	if session == nil {
		return fmt.Errorf("%w: given session is nil", port.ErrSessionNotSaved)
	}

	if session.UUID == "" {
		session.UUID = uuid.New().String()
	}

	return r.ds.Store(ctx, mapper.ToSessionRecord(session))
}

func (r *SessionRepo) Create(ctx context.Context, params *model.SessionParams, rules *model.Rules) (*model.Session, error) {
	uuid := uuid.New().String()
	session, err := model.NewSession(uuid, params, rules)

	if err != nil {
		return nil, err
	}

	if err = r.ds.Store(ctx, mapper.ToSessionRecord(session)); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepo) Delete(ctx context.Context, id string) error {
	return r.ds.Delete(ctx, id)
}
