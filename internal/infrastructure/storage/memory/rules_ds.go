package memory

import (
	"context"
	"fmt"
	"sync"
	"tictactoe/internal/application/port"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
)

type RulesDataSource struct {
	rules map[string]*dsmodel.RulesRecord
	mu    sync.RWMutex
}

func NewRulesDataSource() *RulesDataSource {
	return &RulesDataSource{
		rules: make(map[string]*dsmodel.RulesRecord),
	}
}

func (ds *RulesDataSource) Fetch(_ context.Context, UUID string) (*dsmodel.RulesRecord, error) {
	ds.mu.RLock()
	record, exists := ds.rules[UUID]
	ds.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("memory storage: %w: %s", port.ErrRulesNotFound, UUID)
	}

	return record.Clone(), nil
}

func (ds *RulesDataSource) Store(_ context.Context, record *dsmodel.RulesRecord) error {
	ds.mu.Lock()
	ds.rules[record.UUID] = record.Clone()
	ds.mu.Unlock()
	return nil
}

func (ds *RulesDataSource) Delete(_ context.Context, UUID string) error {
	ds.mu.Lock()
	delete(ds.rules, UUID)
	ds.mu.Unlock()
	return nil
}
