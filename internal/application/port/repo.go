package port

import (
	"context"
	"errors"

	"tictactoe/internal/domain/model"
)

type GetConfig struct {
	State *model.State
	UUID  *string
}

type GetOpt func(*GetConfig)

func WithState(s model.State) GetOpt {
	return func(cfg *GetConfig) {
		cfg.State = &s
	}
}

func WithUUID(UUID string) GetOpt {
	return func(cfg *GetConfig) {
		cfg.UUID = &UUID
	}
}

type SessionRepo interface {
	Get(ctx context.Context, opts ...GetOpt) ([]*model.Session, error)
	Save(ctx context.Context, session *model.Session) error
}

type RulesRepo interface {
	Get(ctx context.Context, rulesUUID string) (*model.Rules, error)
	Save(ctx context.Context, rules *model.Rules) error
}

type UserRepo interface {
	Get(ctx context.Context, login string) (*model.User, error)
	Save(ctx context.Context, user *model.User) error
}

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrRulesNotFound   = errors.New("rules not found")
	ErrUserNotFound    = errors.New("user not found")
)
