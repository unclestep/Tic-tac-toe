package port

import (
	"context"
	"tictactoe/internal/infrastructure/storage/model"
)

type SessionDataSource interface {
	Store(ctx context.Context, session model.SessionRecord) error
	Fetch(ctx context.Context, id string) (model.SessionRecord, error)
	Delete(ctx context.Context, id string) error
}
