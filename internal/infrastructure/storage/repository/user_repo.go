package repository

import (
	"context"
	"fmt"

	"tictactoe/internal/domain/model"
	"tictactoe/internal/infrastructure/storage/ds"
	"tictactoe/internal/infrastructure/storage/mapper"
)

type UserRepo struct {
	uds ds.UserDataSource
}

func NewUserRepo(uds ds.UserDataSource) *UserRepo {
	return &UserRepo{
		uds: uds,
	}
}

func (r *UserRepo) Get(ctx context.Context, login string) (*model.User, error) {
	record, err := r.uds.Fetch(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return mapper.ToUserDomain(record), nil
}

func (r *UserRepo) Save(ctx context.Context, user *model.User) error {
	if user == nil {
		return fmt.Errorf("save user: nil user")
	}
	if user.UUID == "" {
		return fmt.Errorf("save user: empty uuid")
	}
	if err := r.uds.Store(ctx, mapper.ToUserStorage(user)); err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}
