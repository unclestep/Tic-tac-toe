package ds

import (
	"context"

	"tictactoe/internal/infrastructure/storage/model"
)

type SessionDataSource interface {
	Store(ctx context.Context, session *model.SessionRecord) error
	Fetch(ctx context.Context, opts ...SessionOpt) ([]*model.SessionRecord, error)
	Delete(ctx context.Context, UUID string) error
}
