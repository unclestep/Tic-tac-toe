package memory

import (
	"context"
	"fmt"
	"sync"
	"tictactoe/internal/application/port"
	datasourceModel "tictactoe/internal/infrastructure/storage/model"

	"github.com/google/uuid"
)

type SessionDataSource struct {
	sessions map[string]datasourceModel.SessionRecord
	mu       sync.RWMutex
}

func NewMemorySessionDataSource() *SessionDataSource {
	return &SessionDataSource{
		sessions: make(map[string]datasourceModel.SessionRecord),
	}
}
func (ds *SessionDataSource) Fetch(_ context.Context, id string) (datasourceModel.SessionRecord, error) {
	ds.mu.RLock()
	record, exists := ds.sessions[id]
	ds.mu.RUnlock()

	if !exists {
		return datasourceModel.SessionRecord{}, fmt.Errorf("%w: %s", port.ErrSessionNotFound, id)
	}

	return record, nil
}
func (ds *SessionDataSource) Store(_ context.Context, record datasourceModel.SessionRecord) error {
	if record.UUID == "" {
		record.UUID = uuid.New().String()
	}

	ds.mu.Lock()
	ds.sessions[record.UUID] = record
	ds.mu.Unlock()

	return nil
}
func (ds *SessionDataSource) Delete(_ context.Context, id string) error {
	ds.mu.Lock()
	delete(ds.sessions, id)
	ds.mu.Unlock()

	return nil
}
