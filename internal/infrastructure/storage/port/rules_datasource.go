package port

import (
	"context"
	"tictactoe/internal/infrastructure/storage/model"
)

type RulesDataSource interface {
	Store(ctx context.Context, record model.RulesRecord) error
	Fetch(ctx context.Context, id string) (model.RulesRecord, error)
	Delete(ctx context.Context, id string) error
}
