package port

import (
	"context"
	"errors"

	"tictactoe/internal/domain/model"
)

type SessionRepo interface {
	Get(ctx context.Context, opts ...SessionGetOpt) ([]*model.Session, error)
	Save(ctx context.Context, session *model.Session) error
}

type UserRepo interface {
	Get(ctx context.Context, opts ...UserGetOpt) ([]*model.User, error)
	Save(ctx context.Context, user *model.User) error
}

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrReturnedNotOne  = errors.New("returned not one")
)
