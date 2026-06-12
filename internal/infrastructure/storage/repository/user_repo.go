package repository

import (
	"context"
	"fmt"

	"tictactoe/internal/application/port"
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

func UserDomainOptsToDatasourceOpts(c *port.UserGetConfig) []ds.UserOpt {
	var opts []ds.UserOpt
	if c.UUIDs != nil {
		opts = append(opts, ds.WithUUID(c.UUIDs...))
	}
	if c.Logins != nil {
		opts = append(opts, ds.WithLogin(c.Logins...))
	}
	return opts
}

func (r *UserRepo) Get(ctx context.Context, opts ...port.UserGetOpt) ([]*model.User, error) {
	cfg := &port.UserGetConfig{}
	for _, opt := range opts {
		opt.ApplyToUser(cfg)
	}

	records, err := r.uds.Fetch(ctx, UserDomainOptsToDatasourceOpts(cfg)...)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	users := make([]*model.User, len(records))
	for i, record := range records {
		users[i] = mapper.ToUserDomain(record)
	}
	return users, nil
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
