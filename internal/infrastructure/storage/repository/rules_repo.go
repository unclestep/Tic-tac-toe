package repository

import (
	"context"
	"fmt"
	"tictactoe/internal/application/port"
	"tictactoe/internal/domain/model"
	"tictactoe/internal/infrastructure/storage/mapper"
	storagePort "tictactoe/internal/infrastructure/storage/port"

	"github.com/google/uuid"
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
		return nil, fmt.Errorf("%w: %s", port.ErrRulesNotFound, id)
	}
	return mapper.ToRulesDomain(record), nil
}

func (r *RulesRepo) Save(ctx context.Context, rules *model.Rules) error {
	if rules == nil {
		return fmt.Errorf("%w: given rules is nil", port.ErrRulesNotSaved)
	}

	// New rules - dont throw error
	if rules.UUID == "" {
		rules.UUID = uuid.New().String()
	}

	return r.ds.Store(ctx, mapper.ToRulesRecords(rules))
}

func (r *RulesRepo) Delete(ctx context.Context, id string) error {
	return r.ds.Delete(ctx, id)
}
