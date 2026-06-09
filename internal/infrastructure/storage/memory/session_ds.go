package memory

import (
	"context"
	"fmt"
	"sync"
	"tictactoe/internal/application/port"
	dsmodel "tictactoe/internal/infrastructure/storage/model"
)

type SessionDataSource struct {
	sessions map[string]*dsmodel.SessionRecord
	mu       sync.RWMutex
}

func NewSessionDataSource() *SessionDataSource {
	return &SessionDataSource{
		sessions: make(map[string]*dsmodel.SessionRecord),
	}
}
func (ds *SessionDataSource) Fetch(_ context.Context, UUID string) (*dsmodel.SessionRecord, error) {
	ds.mu.RLock()
	record, exists := ds.sessions[UUID]
	ds.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("memory storage: %w: %s", port.ErrSessionNotFound, UUID)
	}

	return record.Clone(), nil
}
func (ds *SessionDataSource) Store(_ context.Context, record *dsmodel.SessionRecord) error {
	ds.mu.Lock()
	ds.sessions[record.UUID] = record.Clone()
	ds.mu.Unlock()
	return nil
}
func (ds *SessionDataSource) Delete(_ context.Context, UUID string) error {
	ds.mu.Lock()
	delete(ds.sessions, UUID)
	ds.mu.Unlock()
	return nil
}
