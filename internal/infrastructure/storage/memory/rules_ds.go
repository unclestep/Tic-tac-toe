package memory

import (
	"context"
	"fmt"
	"sync"
	"tictactoe/internal/application/port"
	datasourceModel "tictactoe/internal/infrastructure/storage/model"

	"github.com/google/uuid"
)

type RulesDataSource struct {
	rules map[string]datasourceModel.RulesRecord
	mu    sync.RWMutex
}

func NewMemoryRulesDataSource() *RulesDataSource {
	return &RulesDataSource{
		rules: make(map[string]datasourceModel.RulesRecord),
	}
}

func (ds *RulesDataSource) Fetch(_ context.Context, id string) (datasourceModel.RulesRecord, error) {
	ds.mu.RLock()
	record, exists := ds.rules[id]
	ds.mu.RUnlock()

	if !exists {
		return datasourceModel.RulesRecord{}, fmt.Errorf("%w: %s", port.ErrRulesNotFound, id)
	}

	return record, nil
}

func (ds *RulesDataSource) Store(_ context.Context, record datasourceModel.RulesRecord) error {
	// New rules - dont throw error
	if record.UUID == "" {
		record.UUID = uuid.New().String()
	}

	ds.mu.Lock()
	ds.rules[record.UUID] = record
	ds.mu.Unlock()

	return nil
}

func (ds *RulesDataSource) Delete(_ context.Context, id string) error {
	ds.mu.Lock()
	delete(ds.rules, id)
	ds.mu.Unlock()
	return nil
}
