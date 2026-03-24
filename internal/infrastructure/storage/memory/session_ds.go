package memory

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"tictactoe/internal/application/port"
	datasourceModel "tictactoe/internal/infrastructure/storage/model"
)

type MemorySessionDataSource struct {
	sessions map[string]datasourceModel.SessionRecord
	mu       sync.RWMutex
}

func NewMemorySessionDataSource() *MemorySessionDataSource {
	return &MemorySessionDataSource{
		sessions: make(map[string]datasourceModel.SessionRecord),
	}
}
func (ds *MemorySessionDataSource) Fetch(_ context.Context, id string) (datasourceModel.SessionRecord, error) {
	ds.mu.RLock()
	record, exists := ds.sessions[id]
	ds.mu.RUnlock()

	if !exists {
		return datasourceModel.SessionRecord{}, fmt.Errorf("%w: %s", port.ErrSessionNotFound, id)
	}

	return record, nil
}
func (ds *MemorySessionDataSource) Store(_ context.Context, record datasourceModel.SessionRecord) error {
	if record.UUID == "" {
		record.UUID = uuid.New().String()
	}

	ds.mu.Lock()
	ds.sessions[record.UUID] = record
	ds.mu.Unlock()

	return nil
}
func (ds *MemorySessionDataSource) Delete(_ context.Context, id string) error {
	ds.mu.Lock()
	delete(ds.sessions, id)
	ds.mu.Unlock()

	return nil
}
