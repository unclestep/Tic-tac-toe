package port

import (
	"context"
	"errors"
	"tictactoe/internal/domain/model"
)

type SessionRepo interface {
	Get(ctx context.Context, sessionUUID string) (*model.Session, error)
	Save(ctx context.Context, session *model.Session) error
}

type RulesRepo interface {
	Get(ctx context.Context, rulesUUID string) (*model.Rules, error)
	Save(ctx context.Context, rules *model.Rules) error
}

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrRulesNotFound   = errors.New("rules not found")
)
