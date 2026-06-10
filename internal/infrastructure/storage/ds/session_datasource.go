package ds

import (
	"context"

	model "tictactoe/internal/domain/model"
	"tictactoe/internal/infrastructure/storage/mapper"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
)

type FetchConfig struct {
	State *string
	UUID  *string
}

type FetchOption func(*FetchConfig)

func WithState(s model.State) FetchOption {
	return func(cfg *FetchConfig) {
		str := mapper.StateToString(s)
		cfg.State = &str
	}
}

func WithUUID(UUID string) FetchOption {
	return func(cfg *FetchConfig) {
		cfg.UUID = &UUID
	}
}

type SessionDataSource interface {
	Store(ctx context.Context, session *dsmodel.SessionRecord) error
	Fetch(ctx context.Context, opts ...FetchOption) ([]*dsmodel.SessionRecord, error)
	Delete(ctx context.Context, UUID string) error
}
