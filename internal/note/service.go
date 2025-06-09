package note

import (
	"context"
	"sync"
	"time"
)

type Payload struct {
	EventAt time.Time
	Text    string
}

type Note struct {
	CreatedAt time.Time
	Tags      []string
	Payload   Payload
}

type Service interface {
	Save(ctx context.Context, userID int64, note Note) error
	SetPending(userID int64, note Note)
	AddTagToPending(userID int64, tag string)
	GetPending(userID int64) (*Note, bool)
	SavePending(ctx context.Context, userID int64) error
}

type InMemoryService struct {
	mu      sync.Mutex
	notes   map[int64][]Note
	pending map[int64]*Note
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{notes: make(map[int64][]Note), pending: make(map[int64]*Note)}
}

func (s *InMemoryService) Save(ctx context.Context, userID int64, note Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notes[userID] = append(s.notes[userID], note)
	return nil
}

func (s *InMemoryService) SetPending(userID int64, note Note) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending[userID] = &note
}

func (s *InMemoryService) AddTagToPending(userID int64, tag string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, ok := s.pending[userID]
	if !ok {
		return
	}
	for _, t := range n.Tags {
		if t == tag {
			return
		}
	}
	n.Tags = append(n.Tags, tag)
}

func (s *InMemoryService) GetPending(userID int64) (*Note, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, ok := s.pending[userID]
	return n, ok
}

func (s *InMemoryService) SavePending(ctx context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, ok := s.pending[userID]
	if !ok {
		return nil
	}
	s.notes[userID] = append(s.notes[userID], *n)
	delete(s.pending, userID)
	return nil
}
