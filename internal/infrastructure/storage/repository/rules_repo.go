package repository

import (
	"context"
	"errors"
	"fmt"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/infrastructure/storage/mapper"
	storagePort "tictactoe/internal/infrastructure/storage/port"
)

type RulesRepo struct {
	ds storagePort.RulesDataSource
}

func NewRulesRepo(ds storagePort.RulesDataSource) *RulesRepo {
	return &RulesRepo{
		ds: ds,
	}
}

func (r *RulesRepo) Get(ctx context.Context, id string) (*model.Rules, error) {
	record, err := r.ds.Fetch(ctx, id)
	if err != nil {
		if errors.Is(err, port.ErrRulesNotFound) {
			return nil, fmt.Errorf("get rules: %w", port.ErrRulesNotFound)
		}
		return nil, fmt.Errorf("get rules: %w", err)
	}
	return mapper.ToRulesDomain(record), nil
}

func (r *RulesRepo) Save(ctx context.Context, rules *model.Rules) error {
	if rules == nil {
		return fmt.Errorf("save rules: nil rules")
	}
	if rules.UUID == "" {
		return fmt.Errorf("save rules: empty uuid")
	}
	if err := r.ds.Store(ctx, mapper.ToRulesStorage(rules)); err != nil {
		return fmt.Errorf("save rules: %w", err)
	}
	return nil
}

func (r *RulesRepo) Delete(ctx context.Context, id string) error {
	if err := r.ds.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete rules: %w", err)
	}
	return nil
}
