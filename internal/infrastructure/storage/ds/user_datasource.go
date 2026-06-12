package ds

import (
	"context"

	"tictactoe/internal/infrastructure/storage/model"
)

type UserDataSource interface {
	Store(ctx context.Context, user *model.UserRecord) error
	Fetch(ctx context.Context, opts ...UserOpt) ([]*model.UserRecord, error)
}
